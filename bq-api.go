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
  getTables(client, "synset_prod")
  fmt.Println("Done")
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

func getTables(client *bigquery.Client, datasetName string) { //(*bigquery.Client.Dataset.Table, *bigquery.Client.Dataset.Table) {
  _, err1 := client.Dataset(datasetName).Metadata(CONTEXT)
  if err1 != nil { // create dataset
    if err := client.Dataset(datasetName).Create(CONTEXT); err != nil {
      log.Fatalf("Failed to create dataset: %v", err)
    }
  }

  _, err2 := client.Dataset(datasetName).Table("projects").Metadata(CONTEXT)
  if err2 != nil { // create projects table
    if err := client.Dataset(datasetName).Table("projects").Create(CONTEXT); err != nil {
      log.Fatalf("Failed to create 'projects' table: %v", err)
    }
  }

  _, err3 := client.Dataset(datasetName).Table("users").Metadata(CONTEXT)
  if err3 != nil { // create users table
    if err := client.Dataset(datasetName).Table("users").Create(CONTEXT); err != nil {
      log.Fatalf("Failed to create 'users' table: %v", err)
    }
  }

  projects := client.Dataset(datasetName).Table("projects")
  users := client.Dataset(datasetName).Table("users")
  fmt.Println(projects, users)
  //return projects, users
}
