/*
Juan Taruc
Serves a form and allows users to upload txt file(does not need checking).
Once uploaded copy the text from txt file and displays it in the webpage using
req.Formfile and io.Copy.
*/
package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func serve_the_webpage(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		file, _, err := req.FormFile("the-file")
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer file.Close()

		src := io.LimitReader(file, 400)

		dst, err := os.Create(filepath.Join(".", "the-file.txt"))
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer dst.Close()

		io.Copy(dst, src)
		dat, err := ioutil.ReadFile("the-file.txt")
		res.Header().Set("Content-Type", "text/html")
		io.WriteString(res, string(dat))
	}
	res.Header().Set("Content-Type", "text/html")
	io.WriteString(res, `
      <form method="POST" enctype="multipart/form-data">
        <input type="file" name="the-file">
        <input type="submit">
      </form>`)
}

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(serve_the_webpage))
}
