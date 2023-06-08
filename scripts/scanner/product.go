package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Product represents a cloud service catalog block.
type Product struct {
	Short     string `json:"short,omitempty"`
	Name      string `json:"name,omitempty"`
	Catalog   string `json:"catalog,omitempty"`
	NameCN    string `json:"name_cn,omitempty"`
	CatalogCN string `json:"catalog_cn,omitempty"`
}

type ProductInfo struct {
	Name      string `json:"name,omitempty"`
	Catalog   string `json:"catalog,omitempty"`
	NameCN    string `json:"name_cn,omitempty"`
	CatalogCN string `json:"catalog_cn,omitempty"`
}

func loadProductDetails() (map[string]*ProductInfo, error) {
	productDetails := make(map[string]*ProductInfo)

	content, err := os.ReadFile(CatalogPath)
	if err != nil {
		fmt.Printf("[ERROR] failed to read config: %s\n", err)
		return nil, err
	}

	if err := json.Unmarshal(content, &productDetails); err != nil {
		fmt.Printf("[ERROR] failed to unmarshal config: %s\n", err)
		return nil, err
	}

	return productDetails, nil
}

// BuildProduct is used to get the product infomation of a resource or data source
func BuildProduct(key string, productDetails map[string]*ProductInfo) *Product {
	if productDetails == nil {
		return nil
	}

	name := getProductName(key)
	if name == "" {
		fmt.Printf("[WARN] failed to get the product name of %s\n", key)
		return nil
	}

	shortName := getShortName(name, productDetails)
	if shortName == "" {
		fmt.Printf("[WARN] failed to get the product short name of %s\n", key)
		return nil
	}

	detail, isExist := productDetails[shortName]
	if !isExist {
		fmt.Printf("[WARN] failed to find details of %s in config\n", shortName)
		return nil
	}

	return &Product{
		Short:     shortName,
		Name:      detail.Name,
		Catalog:   detail.Catalog,
		NameCN:    detail.NameCN,
		CatalogCN: detail.CatalogCN,
	}
}

func getProductName(key string) string {
	if name, isExist := specialResourceMap[key]; isExist {
		return name
	}

	for k, v := range specialResourceKeyMap {
		if strings.Contains(key, k) {
			return v
		}
	}

	nameSlice := strings.Split(key, "_")
	if len(nameSlice) < 2 {
		fmt.Printf("[WARN] the resource %s is invalid\n", key)
		return ""
	}
	return nameSlice[1]
}

func getShortName(key string, productDetails map[string]*ProductInfo) string {
	// get name from specialShortNameMap or upper key
	name, isExist := specialShortNameMap[key]
	if !isExist {
		name = strings.ToUpper(key)
	}

	// try to use `name` as shortname
	if _, ok := productDetails[name]; ok {
		return name
	}
	// try to use `key` as shortname
	if _, ok := productDetails[key]; ok {
		return key
	}

	return ""
}
