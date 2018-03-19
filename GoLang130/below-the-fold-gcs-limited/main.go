//query is limited to 4 jpgs from photos folder

package below_the_fold

import(
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/cloud/storage"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"io/ioutil"
)

const bucket = "throw-away-bucket"
var tpl *template.Template
type app struct{
	ctx 		context.Context
	response 	http.ResponseWriter
	bucket 		*storage.BucketHandle
	client		*storage.Client
}

func init(){
	tpl, _ = template.ParseGlob("*.html")
	http.HandleFunc("/", main)

	//serve the css files in a file server instead of uploading it to gcs and querying it.
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
	//http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("img"))))
	//http.Handle("/photos/", http.StripPrefix("/photos", http.FileServer(http.Dir("photos"))))
}


func main(response http.ResponseWriter, request *http.Request){
	if request.URL.Path != "/"{
		http.NotFound(response, request)
		return
	}

	ctx := appengine.NewContext(request)
	client, err := storage.NewClient(ctx)
	if err != nil{
		log.Errorf(ctx, "Error in main client err")
	}
	defer client.Close()

	a := &app{
		ctx: 		ctx,
		response: 	response,
		bucket: 	client.Bucket(bucket),
		client: 	client,
	}

	//upload all the jpgs in the img and photos folder
	a.uploadFiles()


	var s []string
	//query is now limited to 4 images
	query := &storage.Query{
		Prefix: "photos/",
		Delimiter: "/",
		MaxResults: 4,
	}

	objs, err := a.bucket.List(a.ctx, query)
	if err != nil{
		log.Errorf(a.ctx, "ERROR in query_delimiter")
		return
	}

	//store the public links of the queried result to a slice of string
	for _, obj := range objs.Results {
		s = append(s, obj.MediaLink)
	}

	//execute index.html with the slice s to show the photos
	tpl.ExecuteTemplate(response, "index.html", s)
}


//function upload files explicitly passes the jpgs to 
//copyFile
func (a *app) uploadFiles(){
	for _, fileName := range []string{
						"img/adventure-lg.jpg",
						"img/adventure-sm.jpg",
						"img/adventure.jpg",
						"photos/0.jpg",
						"photos/2.jpg",
						"photos/3.jpg",
						"photos/3.jpg",
						"photos/4.jpg",
						"photos/5.jpg",
						"photos/6.jpg",
						"photos/7.jpg",
						"photos/8.jpg",
						"photos/10.jpg",
						"photos/11.jpg",
						"photos/12.jpg",
						"photos/13.jpg",
						"photos/14.jpg",
						"photos/15.jpg",
						"photos/16.jpg",
						"photos/17.jpg",
						"photos/18.jpg",
						"photos/19.jpg",
						"photos/20.jpg",
						"photos/21.jpg",
						"photos/22.jpg",
						"photos/23.jpg",
						"photos/24.jpg",
					}{
						a.copyFile(fileName)
					}
}


//function copyFile takes in the string fileName and reads that
//fileName in the system then writes it in the gcs in short makes
//a copy of the file to gcs
func (a *app) copyFile(fileName string){
	writer := a.bucket.Object(fileName).NewWriter(a.ctx)
	writer.ACL = []storage.ACLRule{
		{storage.AllUsers, storage.RoleReader},
	}
	writer.ContentType = "image/jpg"

	//read file from folders img and photos
	b, err := ioutil.ReadFile(fileName)
	if err != nil{
		log.Errorf(a.ctx, "Error in copyFile: ", err)
		return
	}

	//writes the read bytes to gcs bucket.
	_, err = writer.Write(b)
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