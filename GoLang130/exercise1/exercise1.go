/*
Juan Taruc Exercise 1 CSCI 130
*/
package main 

import (
	"fmt"
	"strconv"
	)

func main() {
	fmt.Println("Part 1")
	//Prints Hello <NAME>
	var my_name string
	fmt.Print("Enter your name: ")
	fmt.Scan(&my_name)
	fmt.Println("Hello " + my_name + "\n")

	/*
	A program that prints to the terminal asking for a user to enter 
	a small number and a larger number. Print the remainder of the larger 
	number divided by the smaller number. 
	*/
	var num1,num2 int
	fmt.Print("Enter first number: ")
	fmt.Scan(&num1)
	fmt.Print("Enter second number: ")
	fmt.Scan(&num2)
	if num1 > num2 {
		fmt.Println(strconv.Itoa(num1) + "%" + strconv.Itoa(num2) + " = " + strconv.Itoa(num1%num2) + "\n")
	} else {
		fmt.Println(strconv.Itoa(num2) + "%" + strconv.Itoa(num1) + " = " + strconv.Itoa(num2%num1) + "\n")
	}

	//Print all of the even numbers between 0 and 100.
	fmt.Println("Even numbers between 0 ~ 100")
	for i := 0; i <= 100; i ++ {
		if i%2 == 0 {
			fmt.Print(strconv.Itoa(i) + " ")
		}
	}
	fmt.Println("\n")

	/*
	A program that prints the numbers from 1 to 100. But for 
	multiples of three print "Fizz" instead of the number and for the 
	multiples of five print "Buzz". For numbers which are multiples of both 
	three and five print "FizzBuzz".
	*/
	fmt.Println("Multiples of 3 and 5")
	for i := 1; i <= 100; i++ {
		if i%3 == 0 && i%5  == 0 {
			fmt.Print("FizzBuzz ")
		} else if i%3 == 0 {
			fmt.Print("Fizz ")
		} else if i%5 == 0 {
			fmt.Print("Buzz ")
		} else {
			fmt.Print(strconv.Itoa(i) + " ")
		}
	}
	fmt.Println("\n")

	/*
	Find the sum of all the multiples of 3 or 5 below 1000. 
	*/
	fmt.Print("Sum of the numbers below 100 that are multiples of 3 or 5: ")
	var sum int = 0
	for i := 1; i < 1000; i++ {
		if i%3 == 0 || i%5  == 0 {
			sum += i
		}
	}
	fmt.Print(sum)
	fmt.Println("\n")
	

	fmt.Println("Part 2")

	/*
	function(tuple()) will have two returns. The first return should be the argument divided by 2. The 2.
	second return should be a bool that let’s us know whether or not the 
	argument was even.
	*/
	var i int
	fmt.Println("Returns the number divided by 2 and if it is even")
	fmt.Print("Enter a number: ")
	fmt.Scan(&i)
	fmt.Println(tuple(i))
	fmt.Println("\n")

	/*
	A function(max()) with one variadic parameter that finds the greatest 
	number in a list of numbers. 
	*/
	fmt.Println("Function with one variadic parameter. max(1,4,2,45,6) will return biggest number.")
	fmt.Println(max(1,4,2,45,6))
	fmt.Println("\n")

	/*Question 4 on part 2
	What's the value of this expression: (true && false) || (false && true) || 
	!(false && false)? 
	Answer:
		false || false || true
		true
	*/

	//function foo() can be called in below ways
	foo(1,2)
	foo(1,2,3)
	aSlice := []int{1,2,3,4}
	foo(aSlice)
	foo()

	//Problem with projecteuler.net
	var a int = 2
	var b,c int
	var index int = 1
	var is_prime bool
	fmt.Println("Finding 10001st prime number from projecteuler...(almost took 1 minute on my machine)")
	for index <= 10001 {
		is_prime = true
		for b = 2; b < a; b++ {
			if a%b == 0 {
				is_prime = false
			}
		}
		if is_prime {
			c = a
			index++
		}
		a++
	}
	fmt.Println("10001st prime number is " + strconv.Itoa(c))
}


func tuple (i int) (int, bool) {
	return i/2, i%2 == 0
}

func max(nums ...int) int {
   var maxx int = 0
   for _, num := range nums {
   	if maxx < num {
   		maxx = num
   	}
   }
   return maxx
}

func foo(vars ...interface{}) {
	//do something...
}

