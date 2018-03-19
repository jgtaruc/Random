/*
Juan Taruc
Displays the name of the url when the url
is in 127.0.0.1:8080/name using req.URL.Path
*/
package main

import (
	"fmt"
	"net/http"
	"strings"
)

func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
	fs := strings.Split(req.URL.Path, "/")
	fmt.Fprint(res, fs[1]) //only displays the first name after "/"
}

func main() {
	http.HandleFunc("/", serve_the_webpage)

	http.ListenAndServe(":8080", nil)
}
