package code

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func Test_Entity(t *testing.T) {
	// Example YAML data (replace this with your actual YAML input)
	yamlData := `
function_description: |
  This file defines the main service API, including request handlers and middleware configuration.
file_info:
  file_name: "service.go"
  package_name: "main"
  imports:
    - "fmt"
    - "net/http"
constants:
  - name: "MaxRetries"
    value: "5"
    description: "The maximum number of retries for network requests."
structs:
  - name: "Handler"
    fields:
      - field_name: "Name"
        field_type: "string"
      - field_name: "Timeout"
        field_type: "int"
    methods:
      - name: "ServeHTTP"
        params: "w http.ResponseWriter, r *http.Request"
        return_values: ""
        description: "Handles incoming HTTP requests."
methods:
  - name: "StartServer"
    params: ""
    return_values: ""
    description: "Starts the HTTP server."
`

	// Parse the YAML into our ParsedYAML struct
	var parsedData ParsedYAML
	err := yaml.Unmarshal([]byte(yamlData), &parsedData)
	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	// Output the parsed data for demonstration
	fmt.Printf("Function Description: %s\n", parsedData.FunctionDescription)
	fmt.Printf("File Info:\n")
	fmt.Printf("  File Name: %s\n", parsedData.FileInfo.FileName)
	fmt.Printf("  Package Name: %s\n", parsedData.FileInfo.PackageName)
	fmt.Printf("  Imports: %v\n", parsedData.FileInfo.Imports)

	fmt.Println("Constants:")
	for _, c := range parsedData.Constants {
		fmt.Printf("  - Name: %s\n    Value: %s\n    Description: %s\n", c.Name, c.Value, c.Description)
	}

	fmt.Println("Structs:")
	for _, s := range parsedData.Structs {
		fmt.Printf("  - Name: %s\n", s.Name)
		fmt.Println("    Fields:")
		for _, f := range s.Fields {
			fmt.Printf("      - Field Name: %s\n        Field Type: %s\n", f.FieldName, f.FieldType)
		}
		fmt.Println("    Methods:")
		for _, m := range s.Methods {
			fmt.Printf("      - Method Name: %s\n        Params: %s\n        Return Values: %s\n        Description: %s\n",
				m.Name, m.Params, m.ReturnValues, m.Description)
		}
	}

	fmt.Println("Methods:")
	for _, m := range parsedData.Methods {
		fmt.Printf("  - Name: %s\n    Params: %s\n    Return Values: %s\n    Description: %s\n",
			m.Name, m.Params, m.ReturnValues, m.Description)
	}
}
