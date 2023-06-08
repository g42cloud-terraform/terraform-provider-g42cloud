package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
)

// Block represents a configuration block.
//
// When converted to a value, a Block always becomes an instance of an object
// type derived from its defined attributes and nested blocks
type Block struct {
	// Attributes describes any attributes that may appear directly inside
	// the block.
	Attributes map[string]*Attribute `json:"attributes,omitempty"`

	// BlockTypes describes any nested block types that may appear directly
	// inside the block.
	BlockTypes map[string]*NestedBlock `json:"block_types,omitempty"`

	Description string `json:"description,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty"`
}

// NestedBlock represents the embedding of one block within another.
type NestedBlock struct {
	Block    *Block      `json:"block,omitempty"`
	Mode     NestingMode `json:"nesting_mode,omitempty"`
	ForceNew bool        `json:"forcenew,omitempty"`
	MinItems int         `json:"min_items,omitempty"`
	MaxItems int         `json:"max_items,omitempty"`
}

// Attribute represents a configuration attribute, within a block.
type Attribute struct {
	// Type is a type specification that the attribute's value must conform to.
	Type cty.Type `json:"type,omitempty"`

	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Optional    bool   `json:"optional,omitempty"`
	Computed    bool   `json:"computed,omitempty"`
	ForceNew    bool   `json:"forcenew,omitempty"`
	Sensitive   bool   `json:"sensitive,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty"`

	Default interface{} `json:"default,omitempty"`
}

// NestingMode is an enumeration of modes for nesting blocks inside other blocks.
type NestingMode string

const (
	// NestingInvalid indicates that the mode of block is invalid.
	NestingInvalid NestingMode = "invalid"

	// NestingObject indicates that only a single instance of a given
	// block type is permitted, with no labels, and its content should be
	// provided directly as an object value.
	NestingObject NestingMode = "object"

	// NestingList indicates that multiple blocks of the given type are
	// permitted, with no labels, and that their corresponding objects should
	// be provided in a list.
	NestingList NestingMode = "list"

	// NestingSet indicates that multiple blocks of the given type are
	// permitted, with no labels, and that their corresponding objects should
	// be provided in a set.
	NestingSet NestingMode = "set"

	// NestingMap indicates that multiple blocks of the given type are
	// permitted, each with a single label, and that their corresponding
	// objects should be provided in a map whose keys are the labels.
	NestingMap NestingMode = "map"
)

var (
	// SchemaDescriptionBuilder converts helper/schema.Schema Descriptions to Attribute
	// and Block Descriptions.
	SchemaDescriptionBuilder = func(s *schema.Schema) string {
		return s.Description
	}

	// ResourceDescriptionBuilder converts helper/schema.Resource Descriptions to Block
	// Descriptions at the resource top level.
	ResourceDescriptionBuilder = func(r *schema.Resource) string {
		return r.Description
	}
)

// ImpliedType returns the cty.Type that would result from decoding a
// configuration block using the receiving block schema.
//
// ImpliedType always returns a result, even if the given schema is
// inconsistent.
func (b *Block) ImpliedType() cty.Type {
	if b == nil {
		return cty.EmptyObject
	}

	atys := make(map[string]cty.Type)

	for name, attrS := range b.Attributes {
		atys[name] = attrS.Type
	}

	for name, blockS := range b.BlockTypes {
		if _, exists := atys[name]; exists {
			panic("invalid schema, blocks and attributes cannot have the same name")
		}

		childType := blockS.Block.ImpliedType()
		atys[name] = childType

		switch blockS.Mode {
		case NestingObject:
			atys[name] = childType
		case NestingList:
			// We prefer to use a list where possible, since it makes our
			// implied type more complete, but if there are any
			// dynamically-typed attributes inside we must use a tuple
			// instead, which means our type _constraint_ must be
			// cty.DynamicPseudoType to allow the tuple type to be decided
			// separately for each value.
			if childType.HasDynamicTypes() {
				atys[name] = cty.DynamicPseudoType
			} else {
				atys[name] = cty.List(childType)
			}
		case NestingSet:
			if childType.HasDynamicTypes() {
				panic("can't use cty.DynamicPseudoType inside a block type with NestingSet")
			}
			atys[name] = cty.Set(childType)
		case NestingMap:
			// We prefer to use a map where possible, since it makes our
			// implied type more complete, but if there are any
			// dynamically-typed attributes inside we must use an object
			// instead, which means our type _constraint_ must be
			// cty.DynamicPseudoType to allow the tuple type to be decided
			// separately for each value.
			if childType.HasDynamicTypes() {
				atys[name] = cty.DynamicPseudoType
			} else {
				atys[name] = cty.Map(childType)
			}
		default:
			panic("invalid nesting type")
		}
	}

	return cty.Object(atys)
}

func BuildBlockSchema(resource *schema.Resource, ignored bool) *Block {
	if resource == nil {
		// We return an actual (empty) object here, rather than a nil,
		// because a nil result would mean that we don't have a schema at
		// all, rather than that we have an empty one.
		return &Block{}
	}

	deprecated := resource.DeprecationMessage != ""
	if deprecated && ignored {
		return &Block{}
	}

	resourceBlock := configSchema(resource.Schema, ignored)

	desc := ResourceDescriptionBuilder(resource)
	// Only apply Resource Description and Deprecation at top level
	resourceBlock.Description = desc
	resourceBlock.Deprecated = deprecated

	if resourceBlock.Attributes == nil {
		resourceBlock.Attributes = map[string]*Attribute{}
	}

	// Add the implicitly computed "id" field if it doesn't exist
	if resourceBlock.Attributes["id"] == nil {
		resourceBlock.Attributes["id"] = &Attribute{
			Type:     cty.String,
			Computed: true,
		}
	}

	// degrade "region" to computed only
	if regionAttr := resourceBlock.Attributes["region"]; regionAttr != nil {
		regionAttr.Optional = false
		regionAttr.ForceNew = false
	}

	return resourceBlock
}

func configSchema(schemas map[string]*schema.Schema, ignored bool) *Block {
	if len(schemas) == 0 {
		return &Block{}
	}

	ret := &Block{
		Attributes: map[string]*Attribute{},
		BlockTypes: map[string]*NestedBlock{},
	}

	for name, schemaObject := range schemas {
		if isDeprecatedField(name) {
			continue
		}

		if schemaObject.Elem == nil {
			if attribute := configSchemaAttribute(schemaObject, ignored); attribute != nil {
				ret.Attributes[name] = attribute
			}
			continue
		}
		if schemaObject.Type == schema.TypeMap {
			// For TypeMap, it isn't valid for Elem to be a *Resource.
			if _, isResource := schemaObject.Elem.(*schema.Resource); isResource {
				panic(fmt.Errorf("invalid Schema.Elem should be *Schema or *Resource for TypeMap"))
			}
		}

		if schemaObject.Computed && !schemaObject.Optional {
			// Computed-only schemas are always handled as attributes,
			// because they never appear in configuration.
			if attribute := configSchemaAttribute(schemaObject, ignored); attribute != nil {
				ret.Attributes[name] = attribute
			}
			continue
		}

		switch schemaObject.Elem.(type) {
		case *schema.Schema:
			if attribute := configSchemaAttribute(schemaObject, ignored); attribute != nil {
				ret.Attributes[name] = attribute
			}
		case *schema.Resource:
			if nestedBlock := configSchemaBlock(schemaObject, ignored); nestedBlock != nil {
				ret.BlockTypes[name] = nestedBlock
			}
		default:
			// Should never happen for a valid schema
			panic(fmt.Errorf("invalid Schema.Elem %#v; need *Schema or *Resource", schemaObject.Elem))
		}
	}

	return ret
}

// configSchemaAttribute prepares a Attribute representation
// of a schema. This is appropriate only for primitives or collections whose
// Elem is an instance of Schema. Use configSchemaBlock for collections
// whose elem is a whole resource.
func configSchemaAttribute(s *schema.Schema, ignored bool) *Attribute {
	var deprecated bool

	reqd := s.Required
	opt := s.Optional
	computed := s.Computed
	forceNew := s.ForceNew

	// get extent attributes from description
	desc := SchemaDescriptionBuilder(s)
	extent := parseExtentAttribute(desc)

	// update Deprecated field
	if s.Deprecated != "" || hasExtentAttribute(extent, "Deprecated") || hasExtentAttribute(extent, "Internal") {
		deprecated = true
	}

	if deprecated && ignored {
		return nil
	}

	if s.Required || hasExtentAttribute(extent, "Required") {
		reqd = true
		opt = false
		computed = false
	}

	// set the filed as Computed only
	if hasExtentAttribute(extent, "Computed") {
		computed = true
		reqd = false
		opt = false
		forceNew = false
	}

	defaultVal, _ := s.DefaultValue()
	// ignore the default value near to the current time
	// currently, it's used to filter `start_time` in **as_policy** resource
	if v, ok := defaultVal.(string); ok {
		RFC3339NoSecond := "2006-01-02T15:04Z"
		if t, err := time.Parse(RFC3339NoSecond, v); err == nil {
			current := time.Now().UTC()
			diff := math.Abs(float64(current.Unix()) - float64(t.Unix()))
			if diff < 3600 {
				defaultVal = nil
			}
		}
	}

	return &Attribute{
		Type:        configSchemaType(s),
		Optional:    opt,
		Required:    reqd,
		Computed:    computed,
		ForceNew:    forceNew,
		Sensitive:   s.Sensitive,
		Default:     defaultVal,
		Description: desc,
		Deprecated:  deprecated,
	}
}

// configSchemaBlock prepares a NestedBlock representation of
// a schema. This is appropriate only for collections whose Elem is an instance
// of Resource, and will panic otherwise.
func configSchemaBlock(s *schema.Schema, ignored bool) *NestedBlock {
	if s.Deprecated != "" && ignored {
		return nil
	}

	ret := &NestedBlock{}

	nestedResource := s.Elem.(*schema.Resource)
	if nested := configSchema(nestedResource.Schema, ignored); nested != nil {
		ret.Block = nested

		desc := SchemaDescriptionBuilder(s)
		// set these on the block from the attribute Schema
		ret.Block.Description = desc
		ret.Block.Deprecated = s.Deprecated != ""
	}

	switch s.Type {
	case schema.TypeList:
		ret.Mode = NestingList
	case schema.TypeSet:
		ret.Mode = NestingSet
	case schema.TypeMap:
		ret.Mode = NestingMap
	default:
		// Should never happen for a valid schema
		panic(fmt.Errorf("invalid s.Type %s for s.Elem being resource", s.Type))
	}

	ret.ForceNew = s.ForceNew
	ret.MinItems = s.MinItems
	ret.MaxItems = s.MaxItems

	if s.Required && s.MinItems == 0 {
		// configschema doesn't have a "required" representation for nested
		// blocks, but we can fake it by requiring at least one item.
		ret.MinItems = 1
	}
	if s.Optional && s.MinItems > 0 {
		// Historically helper/schema would ignore MinItems if Optional were
		// set, so we must mimic this behavior here to ensure that providers
		// relying on that undocumented behavior can continue to operate as
		// they did before.
		ret.MinItems = 0
	}
	if s.Computed && !s.Optional {
		// MinItems/MaxItems are meaningless for computed nested blocks, since
		// they are never set by the user anyway. This ensures that we'll never
		// generate weird errors about them.
		ret.MinItems = 0
		ret.MaxItems = 0
	}

	return ret
}

// configSchemaType determines the core config schema type that corresponds
// to a particular schema's type.
func configSchemaType(s *schema.Schema) cty.Type {
	switch s.Type {
	case schema.TypeString:
		return cty.String
	case schema.TypeBool:
		return cty.Bool
	case schema.TypeInt, schema.TypeFloat:
		// configschema doesn't distinguish int and float, so helper/schema
		// will deal with this as an additional validation step after
		// configuration has been parsed and decoded.
		return cty.Number
	case schema.TypeList, schema.TypeSet, schema.TypeMap:
		var elemType cty.Type
		switch set := s.Elem.(type) {
		case *schema.Schema:
			elemType = configSchemaType(set)
		case *schema.Resource:
			elemType = configSchema(set.Schema, true).ImpliedType()
		default:
			if set != nil {
				// Should never happen for a valid schema
				panic(fmt.Errorf("invalid Schema.Elem %#v; need *Schema or *Resource", s.Elem))
			}
			// Some pre-existing schemas assume string as default, so we need
			// to be compatible with them.
			elemType = cty.String
		}
		switch s.Type {
		case schema.TypeList:
			return cty.List(elemType)
		case schema.TypeSet:
			return cty.Set(elemType)
		case schema.TypeMap:
			return cty.Map(elemType)
		default:
			// can never get here in practice, due to the case we're inside
			panic("invalid collection type")
		}
	default:
		// should never happen for a valid schema
		panic(fmt.Errorf("invalid Schema.Type %s", s.Type))
	}
}

func parseExtentAttribute(desc string) map[string]bool {
	extra := make(map[string]bool)
	prefix := "schema:"

	if !strings.HasPrefix(desc, prefix) {
		return extra
	}

	validDesc := strings.SplitN(desc, ";", 2)[0]
	allAttr := strings.Split(validDesc[len(prefix):], ",")
	for _, ext := range allAttr {
		if attr := strings.TrimLeft(ext, " "); attr != "" {
			extra[attr] = true
		}
	}

	return extra
}

func hasExtentAttribute(extra map[string]bool, key string) bool {
	return extra[key]
}
