package main

import (
  "fmt"
  "log"
  "os"
  "encoding/json"

  "net/http"

  "cloud.google.com/go/bigquery"
  "golang.org/x/net/context"

  //my packages
  "./packages/schema"
)

// Log handling
var Error *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
var Warn *log.Logger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
var Info *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

// BigQuery globals 
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
    Error.Printf("Failed to create client: %v", err)
    os.Exit(1)
  }
  return client
}

func getTables(client *bigquery.Client, datasetName string) (*bigquery.Table, *bigquery.Table) {
  if  _, err := client.Dataset(datasetName).Metadata(CONTEXT); err != nil {
    if err := client.Dataset(datasetName).Create(CONTEXT); err != nil {
      Error.Printf("Failed to create dataset: %v", err)
      os.Exit(1)
    }
  }

  projects := client.Dataset(datasetName).Table("project")
  projectSchema, _ :=  bigquery.InferSchema(schema.Project{})
  users := client.Dataset(datasetName).Table("user")
  userSchema, _ :=  bigquery.InferSchema(schema.User{})

  if _, err := projects.Metadata(CONTEXT); err != nil {
    if err := projects.Create(CONTEXT, projectSchema); err != nil {
      Error.Printf("Failed to create 'projects' table: %v", err)
      os.Exit(1)
    }
  }

  if _, err := users.Metadata(CONTEXT); err != nil {
    if err := users.Create(CONTEXT, userSchema); err != nil {
      Error.Printf("Failed to create 'users' table: %v", err)
      os.Exit(1)
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
    res.WriteHeader(http.StatusInternalServerError)
    res.Write([]byte("500 - Failed to save Project Data!"))
    Error.Printf("Failed to save project data: %v", projectData)
  } else {
    Info.Printf("Saved project data: %v", projectData)
    res.WriteHeader(http.StatusOK)
    res.Write([]byte("Success!"))
  }
  defer req.Body.Close()
}

func loadUserData(res http.ResponseWriter, req *http.Request) {
  u := users.Uploader()
  decoder := json.NewDecoder(req.Body)
  var userData schema.User
  _ = decoder.Decode(&userData)
  if err := u.Put(CONTEXT, userData); err != nil {
    res.WriteHeader(http.StatusInternalServerError)
    res.Write([]byte("500 - Failed to save User Data!"))
    Error.Printf("Failed to save user data: %v", userData)
  } else {
    Info.Printf("Saved user data: %v", userData)
    res.WriteHeader(http.StatusOK)
    res.Write([]byte("Success!"))
  }
  defer req.Body.Close()
}
