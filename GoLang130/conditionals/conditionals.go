/*
Juan Taruc CSCI 130 Assignment 1(Template with conditional).
This program checks if the number entered is even or odd. It
then opens a file(index.html) and shows the result of the 
check in that template.
*/
package main

import (
	"fmt" 
	"log"
	"os"
	"html/template"
)

//Need a struct because I'm passing two variables
//to the template.
type cond struct {
	Number int
	Is_even bool
}

func main() {
	var check bool //even/odd
	var input int //the user input

    _, err := fmt.Scanf("%d", &input)

    //check if input is a number else show error
    if err != nil {
		log.Fatalln(err)
	}

	//check if the input is even or odd then store the result
	//to variable check.
	if input%2 == 0 {
       check = true
    } else {
        check = false
    }

    //open the template(index.html)
	index, err := template.ParseFiles("index.html")	
	//check for errors in opening the file
	if err != nil {
		log.Fatalln(err)
	}

	//Since there are two variables(input and check) that needs
	//to be passed to the template it needs to be in a struct(cond).
	err = index.Execute(os.Stdout,cond{input,check})
	//check for errors
	if err != nil {
		log.Fatalln(err)
	}
}