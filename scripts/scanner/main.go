package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud"
)

const (
	ProdierName string = "g42cloud/g42cloud"
	CatalogPath        = "./product_catalog.json"
	DocsPath           = "../../docs/"
)

var (
	sourceName       string
	outputDir        string
	isResource       bool
	isData           bool
	ignoreDeprecated bool
	ischecker        bool

	commandLine flag.FlagSet
)

func init() {
	commandLine.Init(os.Args[0], flag.ExitOnError)

	commandLine.BoolVar(&ischecker, "c", false, "Whether to check markdown by schema")
	commandLine.BoolVar(&isResource, "r", false, "Indicates the input name is a resource")
	commandLine.BoolVar(&isData, "d", false, "Indicates the input name is a data source")
	commandLine.BoolVar(&ignoreDeprecated, "ignore-deprecated", true, "Whether to ignore deprecated attributes and blocks")
	commandLine.StringVar(&sourceName, "name", "", "The resource or data source name.")
	commandLine.StringVar(&outputDir, "output-dir", "", "schema json file output directory")

	commandLine.Usage = func() {
		fmt.Fprintf(commandLine.Output(), "Usage of %s:\n\n", os.Args[0])
		commandLine.PrintDefaults()
	}
}

func main() {
	commandLine.Parse(os.Args[1:])

	if sourceName == "" && (isResource || isData) {
		fmt.Printf("-name must be specified\n")
		os.Exit(1)
	}

	resourceType := TypeProvider
	if isResource {
		resourceType = TypeResource
	} else if isData {
		resourceType = TypeData
	} else if sourceName != "" {
		// the name defaults to a resource
		resourceType = TypeResource
	}

	provider := g42cloud.Provider()

	if ischecker {
		var errCode int
		errCode += CheckResourceMarkdown(provider)
		errCode += CheckDataSourcesMarkdown(provider)
		os.Exit(errCode)
	}

	schemas := BuildSchema(provider, resourceType, sourceName, ignoreDeprecated)
	if schemas == nil {
		fmt.Printf("Failed to marshal provider schemas to json: not supported\n")
		os.Exit(2)
	}

	if outputDir == "" {
		jsonSchemas, err := JsonMarshalIndent(schemas, "", "  ")
		if err != nil {
			fmt.Printf("Failed to marshal provider schemas to json: %s\n", err)
			os.Exit(3)
		}

		fmt.Println(string(jsonSchemas))
		os.Exit(0)
	}

	// write the result into output directory
	var subDir string
	if resourceType == TypeResource {
		subDir = "resources"
	} else {
		subDir = "data-sources"
	}
	targetDir := path.Join(outputDir, subDir)
	if err := makeDirEmpty(targetDir); err != nil {
		fmt.Printf("Failed to clean up %s\n", targetDir)
		os.Exit(4)
	}

	if err := writeToFile(schemas, outputDir); err != nil {
		os.Exit(4)
	}
}

// JsonMarshalIndent is similar to json.MarshalIndent, but without escaping.
func JsonMarshalIndent(t interface{}, prefix, indent string) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(t)
	if err != nil {
		return nil, err
	}

	indexBuf := bytes.NewBuffer([]byte{})
	err = json.Indent(indexBuf, buffer.Bytes(), prefix, indent)
	if err != nil {
		return nil, err
	}
	return indexBuf.Bytes(), nil
}
