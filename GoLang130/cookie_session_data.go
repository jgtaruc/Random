/*
Juan Taruc
This program saves user information to a cookie "session_data"
user HMAC to ensure user did not tamper with data.
*/

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"io"
	"net/http"
	"strings"
)

func HMAC_encode(data string) string {
	h := hmac.New(sha256.New, []byte("my_key"))
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func valid(res http.ResponseWriter, req *http.Request) {
	//get cookie.
	cookie, err := req.Cookie("session")

	//cookie not exist.
	if err != nil {
		http.Redirect(res, req, "/", 303)
		return
	}

	//get cookie values
	xs := strings.Split(cookie.Value, "|")

	name := strings.Split(xs[1], "=")     //name[0] = raw name, name[1] = HMAC encoded name
	password := strings.Split(xs[2], "=") //password[0] = raw password, password[1] = HMAC encoded password

	//check if match
	if HMAC_encode(name[0]) != name[1] || HMAC_encode(password[0]) != password[1] {
		io.WriteString(res, "Not Matched!")
	} else {
		io.WriteString(res, "Matched!")
	}

}

func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
	//create/get a cookie "session"
	cookie, err := req.Cookie("session")
	if err != nil {
		id, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:  "session",
			Value: id.String(),
			// Secure: true,
			HttpOnly: true,
		}
	}

	// save user name and password to cookie value using HMAC.
	if req.FormValue("name") != "" && !strings.Contains(cookie.Value, "name") {
		name := req.FormValue("name")
		password := req.FormValue("password")
		cookie.Value = cookie.Value + "|" +
			name + "=" + HMAC_encode(name) + "|" +
			password + "=" + HMAC_encode(password)
	}

	http.SetCookie(res, cookie)

	res.Header().Set("Content-Type", "text/html")
	io.WriteString(res,
		`<form method="POST">
      <label>Name: </label><input type="text" name="name" placeholder="your name" >
      <br><br>
      <label>Password: </label><input type="password" name="password" placeholder="your password" >
      <br><br>
      <button type="submit">Submit</button>
    </form>
    <a href="/validation">Check HMAC</a>`)
}

func main() {
	http.HandleFunc("/", serve_the_webpage)
	http.HandleFunc("/validation", valid)

	http.ListenAndServe(":8080", nil)
}
