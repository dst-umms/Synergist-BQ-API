package main

import (
  "fmt"
  "log"
  "encoding/json"

  "net/http"

  "cloud.google.com/go/bigquery"
  "golang.org/x/net/context"

  //my packages
  "./packages/schema"
)

var CONTEXT = context.Background()
var projects, users *bigquery.Table
const BQ_DATASET string = "devel"

func main() {
  client := getBqClient()
  projects, users = getTables(client, BQ_DATASET)
  http.HandleFunc("/LoadProjectData", loadProjectData)
  http.HandleFunc("/LoadUserData", loadUserData)
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

func getTables(client *bigquery.Client, datasetName string) (*bigquery.Table, *bigquery.Table) {
  if  _, err := client.Dataset(datasetName).Metadata(CONTEXT); err != nil {
    if err := client.Dataset(datasetName).Create(CONTEXT); err != nil {
      log.Fatalf("Failed to create dataset: %v", err)
    }
  }

  projects := client.Dataset(datasetName).Table("project")
  projectSchema, _ :=  bigquery.InferSchema(schema.Project{})
  users := client.Dataset(datasetName).Table("user")
  userSchema, _ :=  bigquery.InferSchema(schema.User{})

  if _, err := projects.Metadata(CONTEXT); err != nil {
    if err := projects.Create(CONTEXT, projectSchema); err != nil {
      log.Fatalf("Failed to create 'projects' table: %v", err)
    }
  }

  if _, err := users.Metadata(CONTEXT); err != nil {
    if err := users.Create(CONTEXT, userSchema); err != nil {
      log.Fatalf("Failed to create 'users' table: %v", err)
    }
  }

  return projects, users
}


func loadProjectData(res http.ResponseWriter, req *http.Request) {
  u := projects.Uploader()
  decoder := json.NewDecoder(req.Body)
  var projectData schema.Project
  _ = decoder.Decode(&projectData)
  if err := u.Put(CONTEXT, projectData); err != nil {
    log.Fatalf("Failed to load project data: %v", err)
    res.WriteHeader(http.StatusInternalServerError)
    res.Write([]byte("500 - Failed to save Project Data!"))
  }
  fmt.Println(projectData)
  res.WriteHeader(http.StatusOK)
  res.Write([]byte("Success!"))
  defer req.Body.Close()
}

func loadUserData(res http.ResponseWriter, req *http.Request) {
  u := users.Uploader()
  decoder := json.NewDecoder(req.Body)
  var userData schema.User
  _ = decoder.Decode(&userData)
  if err := u.Put(CONTEXT, userData); err != nil {
    log.Fatalf("Failed to load user data: %v", err)
    res.WriteHeader(http.StatusInternalServerError)
    res.Write([]byte("500 - Failed to save User Data!"))
  }
  fmt.Println(userData)
  res.WriteHeader(http.StatusOK)
  res.Write([]byte("Success!"))
  defer req.Body.Close()
}
