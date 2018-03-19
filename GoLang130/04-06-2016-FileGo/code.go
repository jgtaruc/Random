package file

import (
	"crypto/sha1"
	"mime/multipart"
	"io"
	"fmt"
)


//gives unique name to file
func getSha(mpf multipart.File) string {
	codename := sha1.New()
	io.Copy(codename, mpf)
	return fmt.Sprintf("%x", codename.Sum(nil))
}