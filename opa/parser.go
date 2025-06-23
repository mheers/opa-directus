package opa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mheers/opa-directus/models"
	"github.com/open-policy-agent/opa/v1/ast"
)

// ParseStatementsWithOpts returns a slice of parsed statements. This is the
// default return value from the parser.
func ParseStatementsWithOpts(filename, input string, popts ast.ParserOptions) ([]ast.Statement, []*ast.Comment, error) {

	parser := ast.NewParser().
		WithFilename(filename).
		WithReader(bytes.NewBufferString(input)).
		WithProcessAnnotation(popts.ProcessAnnotation).
		WithFutureKeywords(popts.FutureKeywords...).
		WithAllFutureKeywords(popts.AllFutureKeywords).
		WithCapabilities(popts.Capabilities).
		WithSkipRules(popts.SkipRules).
		WithRegoVersion(popts.RegoVersion)

	stmts, comments, errs := parser.Parse()

	if len(errs) > 0 {
		return nil, nil, errs
	}

	return stmts, comments, nil
}

func GetAnnotations(filepath string) ([]*ast.Annotations, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	statements, comments, err := ParseStatementsWithOpts(filepath, string(content), ast.ParserOptions{
		ProcessAnnotation: true,
		SkipRules:         false,
		AllFutureKeywords: true,
		RegoVersion:       ast.DefaultRegoVersion,
	})
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		fmt.Println(string(comment.Text))
	}

	annotations := []*ast.Annotations{}

	for _, statement := range statements {
		// check if statement is a rule (by casting to ast.Annotations)
		if a, ok := statement.(*ast.Annotations); ok {
			annotations = append(annotations, a)
		}
	}
	fmt.Println(annotations)

	return annotations, nil
}

func GetCustomSchemata(filepath string) ([]models.ObjectSchema, error) {
	schemata := []models.ObjectSchema{}
	annotations, err := GetAnnotations(filepath)
	if err != nil {
		return nil, err
	}

	for _, annotation := range annotations {
		// Check if the annotation is of type "schema"
		if annotation.Custom["schema"] != nil {
			schema, ok := annotation.Custom["schema"].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid schema type")
			}
			// Convert the schema to models.ObjectSchema through json unmarshalling
			objectSchema := models.ObjectSchema{}

			schemaJSON, err := json.Marshal(schema)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(schemaJSON, &objectSchema)
			if err != nil {
				return nil, err
			}

			schemata = append(schemata, objectSchema)
		}
	}

	// Print the schemata
	for _, schema := range schemata {
		fmt.Printf("Schema: %v\n", schema)
		for propName, prop := range schema.Properties {
			fmt.Printf("  Property: %s, Type: %s\n", propName, prop.Type)
		}
	}
	return schemata, nil
}
