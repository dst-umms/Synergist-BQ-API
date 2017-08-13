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
var bqClient *bigquery.Client 
const PROJECT_ID string = "synergist-170903"
const BQ_DATASET string = "devel"
const PROJECT_TABLE = "project"
const USER_TABLE = "user"
var projects, users *bigquery.Table

// Program entrypoint
func main() {
  bqClient = getBqClient()
  projects, users = getTables()
  http.HandleFunc("/CreateUser", CreateUser)
  http.HandleFunc("/CreateProject", CreateProject)
  http.HandleFunc("/GetProjects", GetProjects)
  http.ListenAndServe(":8080", nil)
}

func getBqClient() (*bigquery.Client) {
  client, err := bigquery.NewClient(CONTEXT, PROJECT_ID)
  if err != nil {
    Error.Printf("Failed to create client: %v", err)
    os.Exit(1)
  }
  return client
}

func getTables() (*bigquery.Table, *bigquery.Table) {
  if  _, err := bqClient.Dataset(BQ_DATASET).Metadata(CONTEXT); err != nil {
    if err := bqClient.Dataset(BQ_DATASET).Create(CONTEXT); err != nil {
      Error.Printf("Failed to create dataset: %v", err)
      os.Exit(1)
    }
  }

  projects := bqClient.Dataset(BQ_DATASET).Table(PROJECT_TABLE)
  projectSchema, _ :=  bigquery.InferSchema(schema.Project{})
  users := bqClient.Dataset(BQ_DATASET).Table(USER_TABLE)
  userSchema, _ :=  bigquery.InferSchema(schema.User{})

  if _, err := projects.Metadata(CONTEXT); err != nil {
    if err := projects.Create(CONTEXT, projectSchema); err != nil {
      Error.Printf("Failed to create '%s' table: %v", PROJECT_TABLE, err)
      os.Exit(1)
    }
  }

  if _, err := users.Metadata(CONTEXT); err != nil {
    if err := users.Create(CONTEXT, userSchema); err != nil {
      Error.Printf("Failed to create '%s' table: %v", USER_TABLE, err)
      os.Exit(1)
    }
  }

  return projects, users
}


func CreateProject(res http.ResponseWriter, req *http.Request) {
  u := projects.Uploader()
  decoder := json.NewDecoder(req.Body)
  var projectData schema.Project
  _ = decoder.Decode(&projectData)
  if err := u.Put(CONTEXT, projectData); err != nil {
    res.WriteHeader(http.StatusInternalServerError)
    res.Write([]byte("Error"))
    Error.Printf("Failed to save project data: %v", projectData)
  } else {
    Info.Printf("Saved project data: %v", projectData)
    res.WriteHeader(http.StatusOK)
    res.Write([]byte("Success"))
  }
  defer req.Body.Close()
}

func GetProjects(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var userData schema.User
  _ = decoder.Decode(&userData)
  // Fetch all the projects associated with the userData.Email
  queryString := fmt.Sprintf(`
    SELECT Name, Desc FROM [%s:%s.%s] WHERE Owner = '%s'
  `, PROJECT_ID, BQ_DATASET, PROJECT_TABLE, userData.Email)
  query := bqClient.Query(queryString)
  query.QueryConfig.UseStandardSQL = false // the above query only works with legacy SQL
  iter, _ := query.Read(CONTEXT)
  type UserProject struct {
    Name  string  `json:"name"`
    Desc  string  `json:"desc"`
  }
  var result struct {
    Data []UserProject  `json:"data"`
  }
  for {
    var current UserProject
    err := iter.Next(&current)
    if err != nil {
      break
    }
    result.Data = append(result.Data, current)
  }
  jsonData, _ := json.Marshal(result)
  res.WriteHeader(http.StatusOK)
  res.Write(jsonData)
  defer req.Body.Close()
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
  u := users.Uploader()
  decoder := json.NewDecoder(req.Body)
  var userData schema.User
  _ = decoder.Decode(&userData)
  if err := u.Put(CONTEXT, userData); err != nil {
    res.WriteHeader(http.StatusInternalServerError)
    res.Write([]byte("Error"))
    Error.Printf("Failed to save user data: %v", userData)
  } else {
    Info.Printf("Saved user data: %v", userData)
    res.WriteHeader(http.StatusOK)
    res.Write([]byte("Success"))
  }
  defer req.Body.Close()
}
