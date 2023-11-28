package main

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	TypeProvider int = iota
	TypeResource
	TypeData
)

// Providers is the top-level object returned when exporting provider schemas
type Providers struct {
	Schemas map[string]*ProviderSchema `json:"provider_schemas,omitempty"`
}

type ProviderSchema struct {
	Provider          *Schema            `json:"provider,omitempty"`
	ResourceSchemas   map[string]*Schema `json:"resource_schemas,omitempty"`
	DataSourceSchemas map[string]*Schema `json:"data_source_schemas,omitempty"`
}

type Schema struct {
	Block   *Block   `json:"block,omitempty"`
	Product *Product `json:"product,omitempty"`
}

func BuildSchema(provider *schema.Provider, sourceType int, resName string, ignored bool) *Providers {
	providerSchema := buildProviderSchema(provider, sourceType, resName, ignored)
	return &Providers{
		Schemas: map[string]*ProviderSchema{
			ProdierName: providerSchema,
		},
	}
}

func buildProviderSchema(provider *schema.Provider, sourceType int, name string, ignored bool) *ProviderSchema {
	var ret ProviderSchema

	if provider == nil {
		return &ret
	}

	switch sourceType {
	case TypeResource:
		ret.ResourceSchemas = buildResourceBlock(provider.ResourcesMap, name, ignored)
	case TypeData:
		ret.DataSourceSchemas = buildResourceBlock(provider.DataSourcesMap, name, ignored)
	default:
		ret.Provider = buildProviderBlock(provider, ignored)
	}

	return &ret
}

func buildResourceBlock(resources map[string]*schema.Resource, key string, ignored bool) map[string]*Schema {
	if resources == nil {
		return nil
	}

	products, err := loadProductDetails()
	if err != nil {
		fmt.Printf("[WARN] failed to load product details: %s\n", err)
	}

	var rs map[string]*Schema
	if key == "all" {
		rs = buildSchemas(resources, products, ignored)
	} else {
		rs = buildSchemaByName(resources, products, key, ignored)
	}

	return rs
}

func buildSchemas(resources map[string]*schema.Resource, products map[string]*ProductInfo, ignored bool) map[string]*Schema {
	if resources == nil {
		return map[string]*Schema{}
	}

	ret := make(map[string]*Schema, len(resources))
	r := regexp.MustCompile("_v[1-9]$")

	for k, v := range resources {
		if r.MatchString(k) {
			continue
		}

		if v.DeprecationMessage != "" && ignored {
			continue
		}

		if isInternalResource(v, k) {
			fmt.Printf("[WARN] %s is only used for internal!\n", k)
			continue
		}

		resSchema := Schema{
			Block:   BuildBlockSchema(v, ignored),
			Product: BuildProduct(k, products),
		}
		ret[k] = &resSchema
	}

	return ret
}

func buildSchemaByName(resources map[string]*schema.Resource, products map[string]*ProductInfo, name string, ignored bool) map[string]*Schema {
	if resources == nil {
		return map[string]*Schema{}
	}

	resource, exist := resources[name]
	if !exist {
		return map[string]*Schema{}
	}

	resSchema := Schema{
		Block:   BuildBlockSchema(resource, ignored),
		Product: BuildProduct(name, products),
	}
	return map[string]*Schema{
		name: &resSchema,
	}
}

func buildProviderBlock(provider *schema.Provider, ignored bool) *Schema {
	if provider == nil {
		return &Schema{}
	}

	providerSchema := Schema{
		Block: configSchema(provider.Schema, ignored),
	}
	return &providerSchema
}
