/*
DESC: API Interface
Author: Joshua Haupt
Last Modified: 10-08-2018
*/

package api

import (
  "encoding/json"
  "fmt"
  "net/http"
  "bytes"
  "../app"
)

/*
  Constant Declarations
*/
const API_ENDPOINT = "http://localhost:8080/newapp"
const AUTH_API_ENDPOINT = "http://localhost:8080/auth"

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

  auth := map[string]interface{}{
    "Username": "admin",
    "Password": "password",
  }

  jsonAuth, err := json.Marshal(auth)
  if  err != nil {
    panic(err)
  }

  authResp, err := http.Post(AUTH_API_ENDPOINT, "application/json", bytes.NewBuffer(jsonAuth))
  if err != nil {
    panic(err)
  }

  var authResult map[string]interface{}

  json.NewDecoder(authResp.Body).Decode(&authResult)
  authToken := authResult["token"]

  fmt.Println("authResult")
  fmt.Println(authResult)
  fmt.Println("authToken")
  fmt.Println(authToken)

  var bearer = "Bearer " + authToken.(string)

  // convert appl to json string
  jsonApplSan, err := json.Marshal(applSan)
  if err != nil {
    panic(err)
  }

  req, err := http.NewRequest(http.MethodPost, API_ENDPOINT, bytes.NewBuffer(jsonApplSan))
  if err != nil {
    panic(err)
  }

  req.Header.Add("Authorization", bearer)
  fmt.Println("Request")
  fmt.Println(req)

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    panic(err)
  }

  var result map[string]interface{}

  json.NewDecoder(resp.Body).Decode(&result)

  fmt.Println("result")
  fmt.Println(result)

  return nil
}
