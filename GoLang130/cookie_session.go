/*
Juan Taruc
This program creates a cookie into the client's machine with uuid
as it's value.
*/

package main

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"io"
	"net/http"
)

func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
	//create a cookie "session"
	cookie, err := req.Cookie("session")
	if err != nil {
		id, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:  "session",
			Value: id.String(),
			// Secure: true,
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)
	}

	io.WriteString(res, "Cookie: "+cookie.Value)
}

func main() {
	http.HandleFunc("/", serve_the_webpage)

	http.ListenAndServe(":8080", nil)
}
