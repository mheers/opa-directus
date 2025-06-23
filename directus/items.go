package directus

import (
	"context"
	"errors"

	"github.com/altipla-consulting/directus-go/v2"
)

func NewItemsClient() *directus.ItemsClient[any] {
	client := NewClient()
	return directus.NewItemsClient[any](client, collectionName)
}

func FetchConfig() (any, error) {
	client := NewItemsClient()
	data, err := client.Get(context.Background(), "1")
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New("data not found")
	}

	return data, nil
}
