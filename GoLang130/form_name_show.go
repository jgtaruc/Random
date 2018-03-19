/*
Juan Taruc
Serves a form and allows user to enter their name, then
show user's name in the webpage.
*/
package main

import (
	"io"
	"net/http"
)

func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
	val := req.FormValue("name")
	io.WriteString(res, `<form method="POST">
		<input type="text" name="name" placeholder="your name" >
		<input type="submit">
	</form>`)

	io.WriteString(res, val)
}

func main() {
	http.HandleFunc("/", serve_the_webpage)
	http.ListenAndServe(":8080", nil)
}
