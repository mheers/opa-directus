package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mheers/opa-directus/directus"
	"github.com/spf13/cobra"
)

var (
	watch    bool
	fetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "fetches directus values for opa policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fetch()
		},
	}
)

// adds a watch flag to the command
func init() {
	fetchCmd.Flags().BoolVarP(&watch, "watch", "w", false, "watch for changes")
}

func fetch() error {
	dst := "bundle/parameters/data.json"

	// mkdir bundle/parameters
	if err := os.MkdirAll("bundle/parameters", 0755); err != nil {
		return err
	}

	if watch {

		// create a channel to receive changes
		changes := make(chan any)

		go directus.Watch(changes)

		// listen for changes on the channel
		for {
			// read from the channel
			select {
			case config := <-changes:
				if config == nil {
					continue
				}
				fmt.Println("Change:", config)

				// convert the config from any into []any
				configList, ok := config.([]any)
				if !ok {
					fmt.Println("Error: config is not a list")
					continue
				}

				configBytes, err := json.Marshal(configList[0]) // Marshal the first element of the list
				if err != nil {
					return err
				}
				if err := os.WriteFile(dst, configBytes, 0644); err != nil {
					return err
				}
			}
		}

	} else {
		config, err := directus.FetchConfig()
		if err != nil {
			return err
		}
		configBytes, err := json.Marshal(config)
		if err != nil {
			return err
		}

		if err := os.WriteFile(dst, configBytes, 0644); err != nil {
			return err
		}
	}

	return nil
}
