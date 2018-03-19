package mem

import (
	"github.com/nu7hatch/gouuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
)

func newVisitor(req *http.Request) (*http.Cookie, error) {
	id, err := uuid.NewV4()
	m := initialModel(id.String())
	if err != nil {
		ctx := appengine.NewContext(req)
		log.Errorf(ctx, "ERROR newVisitor uuid.NewV4: %s", err)
		return nil, err
	}
	return makeCookie(m, id.String(), req)
}

func currentVisitor(m model, id string, req *http.Request) (*http.Cookie, error) {
	return makeCookie(m, id, req)
}

func initialModel(id string) model {
	m := model{
		Name:  "",
		State: false,
		Pictures: []string{
			"one.jpg",
		},
		Id: id,
	}
	return m
}
