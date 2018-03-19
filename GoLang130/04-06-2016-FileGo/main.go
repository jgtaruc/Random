package file

import (
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io"
)


//bucket from google cloud
const bucket = "throw-away-bucket"


func init() {
	http.HandleFunc("/", home_store)
	http.HandleFunc("/retrieve", retrieve)
	http.Handle("/favicon.ico", http.NotFoundHandler())
}


func home_store(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	html := `
		<h1>Upload File</h1>
	    <form method="POST" enctype="multipart/form-data">
			<input type="file" name="data">
			<input type="submit">
	    </form>
	`

	if req.Method == "POST" {
		mpf, hdr, err := req.FormFile("data")
		
		if err != nil {
			log.Errorf(ctx, "*** CTX ERROR: In home_store req.FormFile() ***", err)
			http.Error(res, "*** HTTP ERROR: Unable to upload file ***", http.StatusInternalServerError)
			return
		}
		defer mpf.Close()

		//attempt to store file in google cloud
		filename, err := storeFile(req, mpf, hdr)
		if err != nil {
			log.Errorf(ctx, "*** CTX ERROR: In home_store storeFile() ***")
			http.Error(res, "*** HTTP ERROR: In home_store storeFile(); unable to accept file ***\n" + err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		filenames, err := genCookie(res, req, filename)
		if err != nil {
			log.Errorf(ctx, "*** CTX ERROR: In home_store genCookie(): ", err, "***")
			http.Error(res, "*** HTTP ERROR: In home_store genCookie(); unable to accept file\n" + err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		html += `<h1>Files</h1>`
		for f, _ := range filenames{
			html += `<a href="/retrieve?object=` + f + `">` + f + `</a><br>`
		}
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(res, html)
}


func retrieve(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	object := req.FormValue("object")
	rdr, err := get(ctx, object)
	if err != nil {
		log.Errorf(ctx, "*** CTX ERROR: In retrieve get() ***", err)
		http.Error(res, "*** HTTP ERROR: In retrieve get(); unable to retrieve file " + object + " ***\n" + err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	defer rdr.Close()
	io.Copy(res, rdr)
}