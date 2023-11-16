package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Schema struct {
	Definitions map[string]interface{} `json:"definitions"`
}

func main() {
	agonesSchema, err := FetchAgonesSchema()
	if err != nil {
		fmt.Println("Failed to fetch Agones schema: ", err)
		return
	}

	fmt.Println("Fetching native schema")
	nativeSchema, err := FetchNativeSchema()
	if err != nil {
		fmt.Println("Failed to fetch native schema: ", err)
		return
	}

	schema := MergeSchemas([]Schema{*agonesSchema, *nativeSchema})
	err = SaveSchema("gen/merged.json", schema)
	if err != nil {
		fmt.Println("Failed to save merged schema: ", err)
		return
	}

	fmt.Printf("Schema successfully saved in gen/merged.json")
}

func FetchAgonesSchema() (*Schema, error) {
	file, err := os.Open("agones.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var schema Schema
	err = json.NewDecoder(file).Decode(&schema)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}

func FetchNativeSchema() (*Schema, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/kubernetes/kubernetes/3cd242c51317aed8858119529ccab22079f523b1/api/openapi-spec/swagger.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var schema Schema
	err = json.NewDecoder(resp.Body).Decode(&schema)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}

func MergeSchemas(schemas []Schema) Schema {
	merged := Schema{
		Definitions: make(map[string]interface{}),
	}
	for _, schema := range schemas {
		for key, value := range schema.Definitions {
			merged.Definitions[key] = value
		}
	}
	return merged
}

func SaveSchema(filename string, schema Schema) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(schema)
}
