package main 

import(
	"net/http"
	"html/template"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"io"
)

var tpl *template.Template

func init(){
	http.HandleFunc("/", index)
	http.HandleFunc("/api/check", checkName)
	http.Handle("/favicon.ico", http.NotFoundHandler()) 
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	tpl = template.Must(template.ParseGlob("*.html"))
}
type model struct{
	name string
}

func index(response http.ResponseWriter, request *http.Request){
	if request.Method == "POST"{
		var name model
		name.name = request.FormValue("name")
		ctx := appengine.NewContext(request)

		key := datastore.NewKey(ctx, "User_Name", name.name, 0, nil)
		_, err := datastore.Put(ctx, key, &name)
		if err != nil{
			return
		}
	}
	tpl.ExecuteTemplate(response, "index.html", nil)
}


func checkName(response http.ResponseWriter, request *http.Request){
	ctx := appengine.NewContext(request)
	bs, err := ioutil.ReadAll(request.Body)
	var name model
	name.name = string(bs)

	key := datastore.NewKey(ctx, "User_Name", name.name, 0, nil)
	err = datastore.Get(ctx, key, &name)
	if err != nil{
		tpl.ExecuteTemplate(response, "index.html", nil)
		return
	}
	io.WriteString(response, "true")
}