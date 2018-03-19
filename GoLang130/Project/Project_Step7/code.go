/*
This program continues from step6.
Allows users to login and stores information
about whether or not a user is logged in in both
"User data type" in model.go and in the cookie.
Also shows a logout button when the user is logged in.
Also refactored the code to be more readable.
When the user is logged in, the user can edit his/her
credentials(name and password)

Initial username: user
Initial password: 1234
*/

/*
code.go - contains all functions related to encoding and decoding.
*/

package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)


//encode base64 code to hmac code
func getHMAC(data string) string {
	h := hmac.New(sha256.New, []byte("ourkey"))
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}


//encode marshalled model to base64
func getB64(data []byte) string {
	return base64.URLEncoding.EncodeToString(data) 
}


//decodes base64 string
func decodeB64(b64Data string) []byte {
	bs, err := base64.URLEncoding.DecodeString(b64Data)
	if err != nil {
		log.Println("Error decoding base64 code", err)
	}
	return bs
}


//check if cookie.Value has been tampered
func tampered(s string) bool {
	xs := strings.Split(s, "|")
	code := xs[1]
	data := xs[2]
	return code != getHMAC(data)
}


//marshals the User data type
func marshalModel(model User) []byte {
	bs, err := json.Marshal(model)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return bs
}


//unmarshals the JSON-encoded data
func unmarshalModel(bs []byte) User{
	var model User
	err := json.Unmarshal(bs, &model)
	if err != nil {
		fmt.Println("error unmarshalling: ", err)
	}
	return model
}

