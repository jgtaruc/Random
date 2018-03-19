/*
Juan Taruc
This program creates a cookie that tracks the number
of times a user visited the webpage.
*/

package main

import (
	"io"
	"net/http"
	"strconv"
)

func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("user-cookie")
	//no user-cookie exist, create one
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "user-cookie",
			Value: "0",
		}
	}
	//update existing cookie
	i, _ := strconv.Atoi(cookie.Value)
	i++
	cookie.Value = strconv.Itoa(i)

	http.SetCookie(res, cookie)

	io.WriteString(res, cookie.Value)
}

func throw_the_icon(res http.ResponseWriter, req *http.Request) {
	//do nothing when favicon.ico is requested.
	//does not increment the user-cookie
}

func main() {
	http.HandleFunc("/", serve_the_webpage)
	http.HandleFunc("/favicon.ico", throw_the_icon)
	http.ListenAndServe(":8080", nil)
}
