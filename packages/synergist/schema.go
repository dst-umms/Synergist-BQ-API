package schema

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

type User struct {
  Name  string   `json:"name"`
  Email string   `json:"email"`
  Lab   []string `json:"lab"`
}


