/*
Juan Taruc
This program serves a webpage over https using tls.
*/

package main

import (
	"io"
	"net/http"
)

func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Hello, https")
}

func main() {
	http.HandleFunc("/", serve_the_webpage)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)
}
