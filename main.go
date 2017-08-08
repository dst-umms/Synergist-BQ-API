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
const BQ_DATASET string = "devel"

func main() {
  client := getBqClient()
  projects, _ := getTables(client, BQ_DATASET)
  loadProjectData(projects)
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


func loadProjectData(projects *bigquery.Table) {

  data := []byte(`
    {
  "name": "project1",
  "desc": "project1 description",
  "owner": "vangalamaheshh@gmail.com",
  "users": ["uma.vangala@umassmed.edu"],
  "type": {
    "ngs": {
      "rawdata": {
        "platform": {
          "desc": "This project used Illumina's sequencing platform",
          "keywords": ["Illumina", "HiSeq2500", "Deep sequencing core"]
        },
        "libprep": {
          "desc": "We have used Paired end Illumina reagent kit.",
          "keywords": ["PE", "Paired-End", "Illumina Truseq protocol"]
        },
        "sample": [{
          "name": "sample1",
          "desc": "sample1 - belongs to control group.",
          "keywords": ["control"],
          "files": ["/path/to/sample1_leftmate.fastq.gz", "/path/to/sample1_rightmate.fastq.gz"]
        }, {
          "name": "sample2",
          "desc": "sample2 belongs to treatment group.",
          "keywords": ["treatment"],
          "files": ["/path/to/sample2_leftmate.fastq.gz", "/path/to/sample2_rightmate.fastq.gz"]
        }]
      }
    }
  }
}
  `)
  u := projects.Uploader()

  var projectData schema.Project
  _ = json.Unmarshal(data, &projectData)

  if err := u.Put(CONTEXT, projectData); err != nil {
    log.Fatalf("Failed to create client: %v", err)
  }
  fmt.Println(projectData)
}
