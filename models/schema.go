package models

// Represents a schema of type "object"
type ObjectSchema struct {
	Properties map[string]PropertySchema `json:"properties"`
	Type       SchemaType                `json:"type"`
}

// Strongly-typed string for schema type (usually "object")
type SchemaType string

const (
	SchemaTypeObject SchemaType = "object"
)

// Represents individual property definitions (boolean, integer, string, etc.)
type PropertySchema struct {
	Note string       `json:"note"`
	Type PropertyType `json:"type"`
	Enum []string     `json:"enum,omitempty"`
}

// Strongly-typed string for property types
type PropertyType string

const (
	PropertyTypeBoolean PropertyType = "boolean"
	PropertyTypeInteger PropertyType = "integer"
	PropertyTypeString  PropertyType = "string"
)
