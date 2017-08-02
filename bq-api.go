package main

import (
	"fmt"
	"log"
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	projectID := "synergist-170903"
	client, err := bigquery.NewClient(ctx, projectID)
        if err != nil {
                log.Fatalf("Failed to create client: %v", err)
        }
	datasetName := "ngs_dataset"
	dataset := client.Dataset(datasetName)
	if err := dataset.Create(ctx); err != nil {
                log.Fatalf("Failed to create dataset: %v", err)
        }

        fmt.Printf("Dataset created\n")
}
