package main

import (
  "fmt"
  "log"

  "net/http"

  "cloud.google.com/go/bigquery"
  "golang.org/x/net/context"
)

var CONTEXT = context.Background()

func main() {
  client := getBqClient()
  getTables(client, "synergist_dataset_prod")
  //following 2 lines are hacks - need to refactor as table references being returned from getTables function
  projects := client.Dataset("synergist_dataset_prod").Table("project")
  users := client.Dataset("synergist_dataset_prod").Table("user")
  fmt.Println(projects, users)
  http.ListenAndServe(":8080", nil)
}

func getBqClient() (*bigquery.Client) {
  projectID := "synergist-170903"
  client, err := bigquery.NewClient(CONTEXT, projectID)
  if err != nil {
    log.Fatalf("Failed to create client: %v", err)
  }
  return client
}

func getTables(client *bigquery.Client, datasetName string) { //(projects *bigquery.Client.Dataset.Table, users *bigquery.Client.Dataset.Table) {
  if  _, err := client.Dataset(datasetName).Metadata(CONTEXT); err != nil {
    if err := client.Dataset(datasetName).Create(CONTEXT); err != nil {
      log.Fatalf("Failed to create dataset: %v", err)
    }
  }

  projects := client.Dataset(datasetName).Table("project")

  emptySchema := bigquery.Schema{
    &bigquery.FieldSchema{Name: "empty", Required: false, Type: bigquery.StringFieldType},
  }

  genericSchema := bigquery.Schema{
    &bigquery.FieldSchema{Name: "desc", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "keywords", Required: true, Type: bigquery.StringFieldType},
  }

  sampleSchema := bigquery.Schema{
    &bigquery.FieldSchema{Name: "name", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "files", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "desc", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "keywords", Required: true, Type: bigquery.StringFieldType},
  }

  ngsRawSchema := bigquery.Schema{
    &bigquery.FieldSchema{Name: "platform", Required: true, Type: bigquery.RecordFieldType, Schema: genericSchema},
    &bigquery.FieldSchema{Name: "libprep", Required: true, Type: bigquery.RecordFieldType, Schema: genericSchema},
    &bigquery.FieldSchema{Name: "sample", Required: true, Type: bigquery.RecordFieldType, Schema: sampleSchema},
  }

  ngsSchema := bigquery.Schema{
    &bigquery.FieldSchema{Name: "rawdata", Required: false, Type: bigquery.RecordFieldType, Schema: ngsRawSchema},
    &bigquery.FieldSchema{Name: "analysis", Required: false, Type: bigquery.RecordFieldType, Schema: emptySchema},
  }

  typeSchema := bigquery.Schema{
    &bigquery.FieldSchema{Name: "ngs", Required: false, Type: bigquery.RecordFieldType, Schema: ngsSchema},
    &bigquery.FieldSchema{Name: "imaging", Required: false, Type: bigquery.RecordFieldType, Schema: emptySchema},
  }

  projectSchema := bigquery.Schema{
    &bigquery.FieldSchema{Name: "name", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "desc", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "type", Required: true, Type: bigquery.RecordFieldType, Schema: typeSchema},
  }

  users := client.Dataset(datasetName).Table("user")
//  userSchema := 
  if _, err := projects.Metadata(CONTEXT); err != nil {
    if err := projects.Create(CONTEXT, projectSchema); err != nil {
      log.Fatalf("Failed to create 'projects' table: %v", err)
    }
  }

  if _, err := users.Metadata(CONTEXT); err != nil {
    if err := users.Create(CONTEXT); err != nil {
      log.Fatalf("Failed to create 'users' table: %v", err)
    }
  }

}

