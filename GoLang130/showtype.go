/*
Juan Taruc Assignment 2 Show types of variables.
This program prints the types variables using the
TypeOf() function in the reflect package.
*/
package main

import (
	"fmt"
	"reflect"
)

func main() {
	fmt.Println(reflect.TypeOf(1))
	fmt.Println(reflect.TypeOf(1.0))
	fmt.Println(reflect.TypeOf('a'))
	fmt.Println(reflect.TypeOf("abc"))
	fmt.Println(reflect.TypeOf(true))
}
