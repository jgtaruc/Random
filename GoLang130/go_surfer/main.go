/*
Juan Taruc CSCI 130 Challenge 1
This program runs a web server and serves the main.html template in "/".
Also it can serve the files in side the static folder by going to "/resources".
*/

package main

import (
  "net/http"
  "html/template"
  "log"
)


//Serves main.html, all the required static files are inside the static directory.
func serve_the_surfer(res http.ResponseWriter, req *http.Request) {
  tpl, err := template.ParseFiles("main.html")
  if err != nil {
    log.Fatalln(err)
  }

  tpl.Execute(res, nil)
}


func main() {
  http.HandleFunc("/", serve_the_surfer) // "/" serves main.html in the default address.
  //Loads the files inside the static folder
  http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("static"))))

  log.Println("Listening...")
  http.ListenAndServe(":8000", nil) // "Ip address":8000
}
