/*
This program continues from step6.
Allows users to login and stores information
about whether or not a user is logged in in both
"User data type" in model.go and in the cookie.
Also shows a logout button when the user is logged in.
Also refactored the code to be more readable.
When the user is logged in, the user can edit his/her
credentials(name and password)

Initial username: user
Initial password: 1234
*/

/*
main.go - handles the request of the users.
*/

package main

import (
  "net/http"
  "html/template"
)


var tpl *template.Template


//parse all hmtl pages
func init() {
  tpl, _ = template.ParseGlob("*.html")
}



func main() {
  http.HandleFunc("/", index)
  http.HandleFunc("/login", login)
  http.Handle("/favicon.ico", http.NotFoundHandler()) 

  http.ListenAndServe(":8000", nil)
}


//the page the user sees when he/she is logged in.
//In this page the user can edit user name and password.
func index(res http.ResponseWriter, req *http.Request) {
  cookie := getCookie(res, req) 
  
  user := getModel(cookie) //get the user information encoded to the cookie.

  //if user changes his/her username and or password save it to the cookie.
  if req.Method == "POST" && req.FormValue("user_name") != "" && 
    req.FormValue("password") != "" {
      user.Name = req.FormValue("user_name")
      user.Password = req.FormValue("password")
      user.State = true
  }
  cookie = editCookie(cookie, user)
  http.SetCookie(res, cookie)
  
  tpl.ExecuteTemplate(res, "index.html", user)
}


//login page where the user can enter username and password
//to be able to access the index page.
func login(res http.ResponseWriter, req *http.Request) {
  cookie := getCookie(res, req)

  user := getModel(cookie)
  if req.Method == "POST" && req.FormValue("password") == user.Password && 
    req.FormValue("user_name") == user.Name {
      user.State = true
      cookie = editCookie(cookie, user)
      http.SetCookie(res, cookie)
      http.Redirect(res, req, "/", 302)
      return
  } 

  http.SetCookie(res, cookie)  

  tpl.ExecuteTemplate(res, "login.html", nil)
}