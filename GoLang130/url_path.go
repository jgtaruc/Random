/*
Juan Taruc
Displays the url path using req.URL.Path
*/
package main

import (
	"io"
	"net/http"
)

func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, req.URL.Path)
}

func main() {
	http.HandleFunc("/", serve_the_webpage)

	http.ListenAndServe(":8080", nil)
}
