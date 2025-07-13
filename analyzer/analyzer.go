package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/elastic/go-elasticsearch/v9"
)

func main() {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		APIKey:    "dlhfUkFaZ0JoYnZmbWl6SXZhVk86X1lhdVJIZUg4U1NsSi1rOTVJV1BkQQ==",
	})

	if err != nil {
		log.Fatalf("couldn't start elastic client %s\n", err)
	}

	client.Indices.Create("time")

	json.Marshal()

	client.Index("time", bytes.NewReader())

}
