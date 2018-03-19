package mem

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
	"errors"
)

var tpl *template.Template

func init() {
	tpl, _ = template.ParseGlob("templates/*.html")

	http.HandleFunc("/", index)
	http.HandleFunc("/login/?id=", login)
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/imgs/", fs)
	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {

	ctx := appengine.NewContext(req)
	id, err := get_ID(res, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		src, _, err := req.FormFile("data")
		if err != nil {
			log.Errorf(ctx, "ERROR index req.FormFile: %s", err)
			// TODO: create error page to show user
			http.Redirect(res, req, `/?id=`+id, http.StatusSeeOther)
		}
		err = uploadPhoto(src, id, req)
		if err != nil {
			log.Errorf(ctx, "ERROR index uploadPhoto: %s", err)
			// expired cookie may exist on client
			http.Redirect(res, req, "/logout", http.StatusSeeOther)
			return
		}
	}

	m, err := retrieveMemc(id, req)
	if err != nil {
		log.Errorf(ctx, "ERROR index retrieveMemc: %s", err)
		// expired cookie may exist on client
		http.Redirect(res, req, "/logout", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(res, "index.html", m)
}

func logout(res http.ResponseWriter, req *http.Request) {
	cookie, err := newVisitor(req)
	if err != nil {
		ctx := appengine.NewContext(req)
		log.Errorf(ctx, "ERROR logout getCookie: %s", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(res, cookie)
	http.Redirect(res, req, `/?id=`+cookie.Value, http.StatusSeeOther)
}

func login(res http.ResponseWriter, req *http.Request) {

	ctx := appengine.NewContext(req)
	id, err := get_ID(res, req)
	
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" && req.FormValue("password") == "secret" {
		m, err := retrieveMemc(id, req)
		if err != nil {
			log.Errorf(ctx, "ERROR index retrieveMemc: %s", err)
			// expired cookie may exist on client
			http.Redirect(res, req, "/logout", http.StatusSeeOther)
			return
		}
		m.State = true
		m.Name = req.FormValue("name")
		m.Id = id

		cookie, err := currentVisitor(m, id, req)
		if err != nil {
			log.Errorf(ctx, "ERROR login currentVisitor: %s", err)
			http.Redirect(res, req, `login/?id=`+cookie.Value, http.StatusSeeOther)
			return
		}
		http.SetCookie(res, cookie)

		http.Redirect(res, req, `login/?id=`+cookie.Value, http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(res, "login.html", nil)
}


func get_ID(res http.ResponseWriter, req *http.Request) (string, error) {
	ctx := appengine.NewContext(req)
	var id string
	cookie, err := req.Cookie("session-id")

	//no cookie found
	if err == http.ErrNoCookie {
		//get id from url
		id := req.FormValue("id")

		//no id in url
		if id == "" {
			http.Redirect(res, req, "/logout", http.StatusSeeOther)
			return id, errors.New("Error in get_ID: redirect to logout")
		}

		cookie = &http.Cookie{
			Name: "session-id",
			Value: id,
			//Secure: true,
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)
	}
	//cookie exists
	id = cookie.Value
	return id, nil
}