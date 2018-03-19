package file 

import (
	"net/http"
	"encoding/base64"
	"fmt"
	"encoding/json"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)


func genCookie(res http.ResponseWriter, req *http.Request, filename string) (map[string]bool, error) {
	mss := make(map[string]bool)
	cookie, _ := req.Cookie("file-names")

	//cookie exists
	if cookie != nil {
		bs, err := base64.URLEncoding.DecodeString(cookie.Value)
		if err != nil {
			return nil, fmt.Errorf("*** genCookie ERROR: base64.URLEncoding.DecodeString(): %s ***", err)
		}
		err = json.Unmarshal(bs, &mss)
		if err != nil {
			return nil, fmt.Errorf("*** genCookie ERROR: json.Unmarshal: %s ***", err)
		}
	}

	//cookie not exists
	mss[filename] = true
	bs, err := json.Marshal(mss)
	if err != nil {
		return mss, fmt.Errorf("genCookie ERROR: json.Marshal: %s", err)
	}
	b64 := base64.URLEncoding.EncodeToString(bs)

	ctx := appengine.NewContext(req)
	log.Infof(ctx, "Cookie Json: %s", string(bs))

	http.SetCookie(res, &http.Cookie{
		Name: "file-names",
		Value: b64,
		})

	return mss, nil
}