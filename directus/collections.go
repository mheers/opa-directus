package directus

import (
	"context"
	"fmt"
	"os"

	"github.com/altipla-consulting/directus-go/v2"
	"github.com/mheers/opa-directus/models"
)

func NewClient() *directus.Client {

	token := os.Getenv("DIRECTUS_TOKEN")
	if token == "" {
		panic("missing DIRECTUS_TOKEN")
	}

	url := os.Getenv("DIRECTUS_URL")
	if url == "" {
		panic("missing DIRECTUS_URL")
	}

	return directus.NewClient(url, token)
}

const collectionName = "pdp-config"

func CreateCollectionForSchemata(schemata []models.ObjectSchema) error {
	client := NewClient()
	collection := &directus.Collection{}

	// remove old collection if it exists
	exists := false
	_, err := client.Collections.Get(context.Background(), collectionName)
	if err == nil {
		exists = true
	}

	if exists {
		fmt.Println("Collection already exists, deleting...")
		if err := client.Collections.Delete(context.Background(), collectionName); err != nil {
			return err
		}
	}

	// fields, err := client.Fields.ListCollection(context.Background(), collectionName)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(fields)

	collection.Collection = collectionName
	collection.Meta = directus.CollectionMeta{
		Color:            "blue",
		Singleton:        true,
		Icon:             "settings",
		Note:             "OPA policy configuration",
		Versioning:       true,
		Hidden:           false,
		ArchiveAppFilter: true,
		Collapse:         directus.CollectionCollapseOpen,
		Accountability: directus.Nullable[directus.Accountability]{
			Value: directus.AccountabilityAll,
			Valid: true,
		},
	}
	collection.Schema = &directus.CollectionSchema{
		Name:    collectionName,
		Comment: "OPA policy configuration",
	}

	_, err = client.Collections.Create(context.Background(), collection)
	if err != nil {
		return err
	}

	// create fields
	fields := []*directus.Field{}
	for _, schema := range schemata {
		if schema.Type != models.SchemaTypeObject {
			continue
		}

		for name, prop := range schema.Properties {
			field := &directus.Field{
				Collection: collectionName,
				Field:      name,
				// one of [bigInteger, boolean, date, dateTime, decimal, float, integer, json, string, text, time, timestamp, binary, uuid, alias, hash, csv, geometry, geometry.Point, geometry.LineString, geometry.Polygon, geometry.MultiPoint, geometry.MultiLineString, geometry.MultiPolygon, unknown, o2m, m2m, m2a, o2a, files, translations, null]
				Type: directus.FieldType(prop.Type),
				Meta: directus.FieldMeta{
					Interface: func() string {
						if len(prop.Enum) > 0 {
							return "select-dropdown"
						}
						return ""
					}(),
					Note: prop.Note,
					Options: &directus.FieldOptions{
						Choices: directus.FieldChoices{
							Values: toAnySlice(prop.Enum),
						},
					},
					Width:    directus.FieldWidthFull,
					System:   false,
					Required: true,
				},
			}
			fields = append(fields, field)
		}
	}

	for _, field := range fields { // TODO: keep order
		_, err = client.Fields.Create(context.Background(), field)
		if err != nil {
			return err
		}
	}

	fmt.Println("Collection created successfully")

	return nil
}

func toAnySlice[T any](slice []T) []any {
	anySlice := make([]any, len(slice))
	for i, v := range slice {
		anySlice[i] = v
	}
	return anySlice
}
