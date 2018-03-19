package query_delimiter

import (
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/cloud/storage"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"io"
)

const bucket = "throw-away-bucket"
var tpl *template.Template
type app struct{
	ctx 		context.Context
	response 	http.ResponseWriter
	request 	*http.Request
	bucket 		*storage.BucketHandle
	client		*storage.Client
}


func init(){
	tpl, _ = template.ParseGlob("*.html")
	http.HandleFunc("/", main)
}


func main(response http.ResponseWriter, request *http.Request){
	if request.URL.Path != "/"{
		http.NotFound(response, request)
		return
	}

	ctx := appengine.NewContext(request)
	client, err := storage.NewClient(ctx)
	if err != nil{
		log.Errorf(ctx, "*** ERROR: main NewClient: ", err, " ***")
	}
	defer client.Close()

	a := &app{
		ctx: 		ctx,
		response: 	response,
		request:	request,
		bucket: 	client.Bucket(bucket),
		client: 	client,
	}

	a.createFiles()
	a.showFiles()
	a.query_delimiter()
}


func (a *app) createFiles(){
	for _, fileName := range []string{"file1", "file2", 
							"folder1/file3", 
							"folder1/file4", 
							"folder1/folder2/file5", 
							"folder3/"}{
	a.createFile(fileName)
	}
}


func (a *app) createFile(fileName string){
	writer := a.bucket.Object(fileName).NewWriter(a.ctx)
	writer.ContentType = "text/plain"

	_, err := writer.Write([]byte(fileName))
	if err != nil{
		log.Errorf(a.ctx, "createFile: unable to write data to bucket")
		return
	}
	err = writer.Close()
	if err != nil{
		log.Errorf(a.ctx, "createFile Close")
		return
	}
}


func (a *app) query_delimiter(){
	query := &storage.Query{
		Delimiter: "/",
	}

	objs, err := a.bucket.List(a.ctx, query)
	if err != nil{
		log.Errorf(a.ctx, "ERROR in query_delimiter")
		return
	}

	html := `<h1>Results with given delimiter /</h1>`
	for _, obj := range objs.Results{
		html += obj.Name + `<br>`
	}

	html += `<br>
			 <h1>Prefixes found (storage.ObjectList Prefixes) </h1>`
	for _, p := range objs.Prefixes{
		html += p + `<br>`
	}

	io.WriteString(a.response, html)
}


func (a *app) showFiles(){
	client, err := storage.NewClient(a.ctx)
	if err != nil{
		log.Errorf(a.ctx, "%v", err)
		return
	}
	defer client.Close()

	objs, err := client.Bucket(bucket).List(a.ctx, nil)
	if err != nil{
		log.Errorf(a.ctx, "%v", err)
		return
	}

	tpl.ExecuteTemplate(a.response, "index.html", objs.Results)

}