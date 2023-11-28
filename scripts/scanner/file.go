package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func makeDirEmpty(dir string) error {
	f, err := os.Stat(dir)
	if err != nil {
		return os.MkdirAll(dir, 0750)
	}

	if f.IsDir() {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, file := range entries {
			fileName := path.Join(dir, file.Name())
			os.RemoveAll(fileName)
		}
		return nil
	}

	return fmt.Errorf("%s is not a dir", dir)
}

func writeToFile(schemas *Providers, output string) error {
	for p, schema := range schemas.Schemas {
		for r, s := range schema.ResourceSchemas {
			fullName := strings.SplitN(r, "_", 2)
			if len(fullName) != 2 {
				fmt.Printf("[WARN] the format of resource %s is not valid\n", r)
				continue
			}
			fileName := fullName[1] + ".json"
			filePath := path.Join(output, "resources", fileName)
			fmt.Printf("[DEBUG] writing resource %s into %s\n", r, filePath)

			singleResource := &Providers{
				Schemas: map[string]*ProviderSchema{
					p: {
						ResourceSchemas: map[string]*Schema{r: s},
					},
				},
			}
			content, err := JsonMarshalIndent(singleResource, "", "  ")
			if err != nil {
				fmt.Printf("Failed to marshal resource %s schemas to json: %s\n", r, err)
				return err
			}

			// write to file
			if err := writeSingleFile(filePath, content); err != nil {
				fmt.Printf("Failed to writing resource %s into %s: %s\n", r, filePath, err)
				return err
			}
		}

		for d, s := range schema.DataSourceSchemas {
			fullName := strings.SplitN(d, "_", 2)
			if len(fullName) != 2 {
				fmt.Printf("[WARN] the format of data source %s is not valid\n", d)
				continue
			}
			fileName := fullName[1] + ".json"
			filePath := path.Join(output, "data-sources", fileName)
			fmt.Printf("[DEBUG] writing data source %s into %s\n", d, filePath)

			singleDataSource := &Providers{
				Schemas: map[string]*ProviderSchema{
					p: {
						DataSourceSchemas: map[string]*Schema{d: s},
					},
				},
			}
			content, err := JsonMarshalIndent(singleDataSource, "", "  ")
			if err != nil {
				fmt.Printf("Failed to marshal data source %s schemas to json: %s\n", d, err)
				return err
			}

			// write to file
			if err := writeSingleFile(filePath, content); err != nil {
				fmt.Printf("Failed to writing data source %s into %s: %s\n", d, filePath, err)
				return err
			}
		}
	}

	fmt.Println("Success!")
	return nil
}

func writeSingleFile(file string, body []byte) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(body); err != nil {
		return err
	}

	return nil
}
