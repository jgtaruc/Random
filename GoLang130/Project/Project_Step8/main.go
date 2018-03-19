/*
This program continues from step7.
Allows the user to logout when the user is logged in

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
  http.HandleFunc("/logout", logout)
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


//logout enables the user to logout once the user is logged in.
func logout(res http.ResponseWriter, req *http.Request) {
  cookie := getCookie(res, req)

  user := getModel(cookie)
  user.State = false
  cookie = editCookie(cookie, user)
  http.SetCookie(res, cookie)
  http.Redirect(res, req, "/login", 302)
  return
}