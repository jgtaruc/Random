/*
This program continues from step4.
This program integrates HMAC ensuring
that nobody tampers with the cookie.
*/


package main

import (
  "net/http"
  "github.com/nu7hatch/gouuid"
  "html/template"
  "log"
  "encoding/json"
  "crypto/sha256"
  "crypto/hmac"
  "encoding/base64"
  "strings"
  "io"
  "fmt"
)


type User struct {
  Name string
  Age string
}


func set_user(req *http.Request, user *User) string {
  user.Name = req.FormValue("user_name")
  user.Age = req.FormValue("user_age")

  bs, err := json.Marshal(user) //marshal to type user to json
  if err != nil {
    fmt.Println("error: ", err)
  }

  b64 := base64.URLEncoding.EncodeToString(bs) //encode type user to base64 to compensate for quotes in the cookie

  return b64
}


//encde user data to hmac
func getCode(data string) string {
  h := hmac.New(sha256.New, []byte("ourkey"))
  io.WriteString(h, data)
  return fmt.Sprintf("%x", h.Sum(nil))
}


//Serves index.html
func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
  tpl, err := template.ParseFiles("index.html")
  if err != nil {
    log.Fatalln(err)
  }

  user_data := set_user(req, new(User))
  hmac_code := getCode(user_data)

  //get cookie.
  cookie, err := req.Cookie("session-fino")
  //create cookie if not exists.
  if err != nil {
    id, _ := uuid.NewV4()
    cookie = &http.Cookie{
      Name: "session-fino",
      Value: id.String() + "|" + user_data + "|" + hmac_code,
      // Secure: true,
      HttpOnly: true,
    }
  }

  if req.FormValue("user_name") != "" {
    //Not sure about this part and the instruction did not say anything,
    //I then decided that in every form submission made by the user the 
    //cookie value changes except for the uuid.
    cookie_id := strings.Split(cookie.Value, "|")
    cookie.Value = cookie_id[0] + "|" + user_data + "|" + hmac_code
  }

  //get the base64 user data and hmac code in the cookie
  xs := strings.Split(cookie.Value, "|")
  data := xs[1]
  code := xs[2]

  /* decode base64 code
  data1, err := base64.StdEncoding.DecodeString(data)
  if err != nil {
    fmt.Println("error:", err)
    return
  }
  fmt.Printf("%q\n", data1)
  */

  encode := getCode(data) //get the hmac code of data
  msg := "" //message placeholder for code valid/invalid.
  if code == encode {
    msg = "code valid: " + code + " == " + encode
  } else {
   msg = "code invalid: " + code + " != " + encode
  }

  http.SetCookie(res, cookie)

  res.Header().Set("Content-Type", "text/html")
  io.WriteString(res, msg)
  
  tpl.Execute(res, nil)
}


func main() {
  http.HandleFunc("/", serve_the_webpage) 
  http.Handle("/favicon.ico", http.NotFoundHandler()) 

  log.Println("Listening...")
  http.ListenAndServe(":8080", nil)
}
