package cmd

import (
	"github.com/mheers/opa-directus/directus"
	"github.com/mheers/opa-directus/opa"
	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "generates directus collections from OPA policies",
		Run: func(cmd *cobra.Command, args []string) {
			generate()
		},
	}
)

func generate() {
	// get metadata/schema from opa policies
	file := "bundle/demo/demo.rego"
	schemata, err := opa.GetCustomSchemata(file)
	if err != nil {
		panic(err)
	}

	// create metadata/schema in directus
	if err := directus.CreateCollectionForSchemata(schemata); err != nil {
		panic(err)
	}
}
