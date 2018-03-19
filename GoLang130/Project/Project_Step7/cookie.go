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
cookie.go - handles the cookie related codes for the website

cookie structure = uuid|hmac|base64
*/

package main

import (
	"fmt"
	"strings"
	"net/http"
	"github.com/nu7hatch/gouuid"
)


//gets the cookie saved in the browser. If no cookie is found, then create
//a new one.
func getCookie(res http.ResponseWriter, req *http.Request) *http.Cookie {
  cookie, err := req.Cookie("session-fino")

  //cookie does not exists, create new cookie with initial userinfo:
  //	username = user
  //	password = 1234
  //	state = false
  if err != nil {
    model := initialModel()
    id, _ := uuid.NewV4()
    return createCookie(id.String(), model) 
  }

  //checks the cookie if it is tampered, if it is
  //give the user a new cookie with initial user info
  if tampered(cookie.Value) {
   	fmt.Println("cookie has been tampered!")
  	model := initialModel()
    id, _ := uuid.NewV4()
    return createCookie(id.String(), model) 
  }

  return cookie
}


//Edits the cookie save in the browser whenever the user changes 
//username and or password.
func editCookie(cookie *http.Cookie, model User) *http.Cookie {
	xs := strings.Split(cookie.Value, "|")
	uuid := xs[0]

	bs := marshalModel(model)
	b64 := getB64(bs)
	code := getHMAC(b64)

	cookie.Value = uuid + "|" + code + "|" + b64
	return cookie
}


//creates a cookie if one does not exists in the browser
func createCookie(id string, model []byte) *http.Cookie {
	b64 := getB64(model) //convert to base64
	code := getHMAC(b64) //get hmac code of base64

	cookie := &http.Cookie{
		Name: "session-fino",
		Value: id + "|" + code + "|" + b64,
		//Secure: true,
		HttpOnly: true,
	}
	return cookie
}