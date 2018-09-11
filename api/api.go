package api

import (
  "encoding/json"
  //"fmt"
  "net/http"
  "bytes"
  "../app"
)

/*
  Constant Declarations
*/
const API_ENDPOINT = "http://localhost:8080/apps"

func SendApp(appl *app.App) error {

  // extract only the necessary app fields
  applSan := map[string]interface{}{
    "Company": appl.Company,
    "Position": appl.Position,
    "Contact": appl.Contact,
    "Source": appl.Source,
    "Heading": appl.Heading,
    "Note1": appl.Note1,
    "Note2": appl.Note2,
    "Skill1": appl.Skill1,
    "Skill2": appl.Skill2,
    "Skill3": appl.Skill3,
    "Local": appl.Local,
    "Url": appl.Url,
    "MailTo": appl.MailTo,
  }

  // convert appl to json string
  jsonApplSan, err := json.Marshal(applSan)
  if err != nil {
    panic(err)
  }
  // fmt.Println("JsonData")
  // fmt.Println(jsonApplSan)

  // send appl to
  resp, err := http.Post(API_ENDPOINT, "application/json", bytes.NewBuffer(jsonApplSan))
  if err != nil {
    panic(err)
  }

  var result map[string]interface{}

  json.NewDecoder(resp.Body).Decode(&result)

  // fmt.Println(result)
  // fmt.Println(result["data"])

  return nil
}
