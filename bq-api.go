package main

import (
  "fmt"
//  "log"
  "encoding/json"

  //"net/http"

  //"cloud.google.com/go/bigquery"
  //"golang.org/x/net/context"
)

//var CONTEXT = context.Background()
//const BQ_DATASET string = "development" 

func main() {
  /*client := getBqClient()
  getTables(client, BQ_DATASET)
  //following 2 lines are hacks - need to refactor as table references being returned from getTables function
  projects := client.Dataset("DATASET").Table("project")
  users := client.Dataset("DATASET").Table("user")
  fmt.Println(projects, users)
  http.ListenAndServe(":8080", nil)*/
  loadProjectData()
}

/*func getBqClient() (*bigquery.Client) {
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
    &bigquery.FieldSchema{Name: "keywords", Required: true, Repeated: true, Type: bigquery.StringFieldType},
  }

  sampleSchema := bigquery.Schema{
    &bigquery.FieldSchema{Name: "name", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "files", Required: true, Repeated: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "desc", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "keywords", Required: true, Repeated: true, Type: bigquery.StringFieldType},
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
    &bigquery.FieldSchema{Name: "owner", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "users", Required: false, Repeated: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "type", Required: true, Type: bigquery.RecordFieldType, Schema: typeSchema},
  }

  users := client.Dataset(datasetName).Table("user")

  userSchema := bigquery.Schema{
    &bigquery.FieldSchema{Name: "name", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "email", Required: true, Type: bigquery.StringFieldType},
    &bigquery.FieldSchema{Name: "lab", Required: false, Repeated: true, Type: bigquery.StringFieldType},
  }

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

}
*/

func loadProjectData() {

  type Project struct {
      Name  string `json:"name"`
      Desc  string `json:"desc"`
      Owner string `json:"owner"`
      Users []string `json:"users"`
      Type  struct {
        Ngs     struct {
          Rawdata  struct {
            Platform struct {
              Desc     string   `json:"desc"`
              Keywords []string `json:"keywords"`
            } `json:"platform"`
            Libprep struct {
              Desc     string   `json:"desc"`
              Keywords []string `json:"keywords"`
            } `json:"libprep"`
            Sample []struct {
              Name     string   `json:"name"`
              Desc     string   `json:"desc"`
              Keywords []string `json:"keywords"`
              Files    []string `json:"files"`
            } `json:"sample"`
          } `json:"rawdata"`
          Analysis interface{} `json:"analysis"`
        } `json:"ngs"`
        Imaging interface{} `json:"imaging"`
      } `json:"type"`
    }

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
      },
      "analysis": null
    },
    "imaging": null
  }
}
  `)

  var projectData Project
  _ = json.Unmarshal(data, &projectData)
  fmt.Println(projectData)
}
