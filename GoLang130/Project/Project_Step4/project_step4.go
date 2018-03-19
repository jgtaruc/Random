/*
This program continues from step3.
This program saves the name and age of the 
user to a struct type user. It also Marshals
the type user to JSON and encodes to base 64 
and saves it to the cookie.
*/


package main

import (
  "net/http"
  "github.com/nu7hatch/gouuid"
  "html/template"
  "log"
  "encoding/json"
  "encoding/base64"
  "strings"
)


type User struct {
  Name string
  Age string
}


//saves the submitted form values to a type use and encodes it to base64.
func set_user(req *http.Request, user *User) string {
  user.Name = req.FormValue("user_name")
  user.Age = req.FormValue("user_age")

  bs, _ := json.Marshal(user) //marshal to type user to json
  //fmt.Println(string(bs))
  b64 := base64.URLEncoding.EncodeToString(bs) //encode type user to base64
  //fmt.Println(b64)

  return b64
}


//Serves index.html
func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
  tpl, err := template.ParseFiles("index.html")
  if err != nil {
    log.Fatalln(err)
  }

  //get cookie.
  cookie, err := req.Cookie("session-fino")
  //create cookie if not exists.
  if err != nil {
    id, _ := uuid.NewV4()
    cookie = &http.Cookie{
      Name: "session-fino",
      Value: id.String(),
      // Secure: true,
      HttpOnly: true,
    }
  }

  //updates the cookie value in every form submission.
  user := new(User) //create new User instance.
  b64 := set_user(req, user) //get the encoded user type.

  //if form is not blank(valid) then update cookie value.
  if req.FormValue("user_name") != "" {
    //Not sure about this part and the instruction did not say anything,
    //I then decided that in every form submission made by the user the 
    //cookie value changes except for the uuid.
    cookie_id := strings.Split(cookie.Value, "|")
    cookie.Value = cookie_id[0] + "|" + b64
  }

  http.SetCookie(res, cookie)

  tpl.Execute(res, nil)
}


func main() {
  http.HandleFunc("/", serve_the_webpage) 
  http.Handle("/favicon.ico", http.NotFoundHandler()) 

  log.Println("Listening...")
  http.ListenAndServe(":8080", nil)
}
