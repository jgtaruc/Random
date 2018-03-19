/*
Juan Taruc
Displays the name in the url when the url
is in 127.0.0.1:8080/?n="name" using req.FormValue
*/
package main

import (
	"io"
	"net/http"
)

func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
	key := "n"
	name := req.FormValue(key)
	io.WriteString(res, name)
}

func main() {
	http.HandleFunc("/", serve_the_webpage)
	http.ListenAndServe(":8080", nil)
}
