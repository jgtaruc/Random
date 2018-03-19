package file

import (
	"net/http"
	"mime/multipart"
	"strings"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"fmt"
)


func storeFile(req *http.Request, mpf multipart.File, hdr *multipart.FileHeader) (string, error) {
	ext, err := fileType(req, hdr) //get type of file based on file extension

	if err != nil {
		return "", err
	}	

	name := getSha(mpf) + `.` + ext //give unique name to the file
	mpf.Seek(0, 0)

	ctx := appengine.NewContext(req)
	return name, store(ctx, name, mpf)
}


//determines filetype base on file extension
func fileType(req *http.Request, hdr *multipart.FileHeader) (string, error) {
	ext := hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:] //string slicing
	ctx := appengine.NewContext(req)
	log.Infof(ctx, "+++ fileType INFO: File Extension: %s +++", ext)

	switch ext {
	case "jpg", "jpeg", "png", "txt", "md": //acceptable files
		return ext, nil
	}
	return ext, fmt.Errorf("*** fileType ERROR: File is of type %s. Only accepts jpg, jpeg, png, txt, md for now. ***", ext)

}


