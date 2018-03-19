/*
Juan Taruc CSCI 176 Assignment 1.
This program creates 3 child processes:
  * 1st child will wait for user input n and sends it to
    2nd and child processes via pipes.
  * 2nd and 3rd child processes will read from pipe and computes the
    nth fibonacci number using recursive and iterative ways respectively.
  * 2nd and 3rd child processes will time the execution of computing
    the nth fibonacci number and then writes the execution time to pipes.
  * 1st child process will read the execution times of 2nd and 3rd processes
    via pipe and displays it on the screen.
*/

#include <iostream>
#include <cstdlib> //exit()
#include <unistd.h> //fork() and pipe()
#include <sys/wait.h> //wait()
#include <sys/time.h> //timing
using namespace std;

//Given code in the examples in class for getting the current time.
//macro (in-line expansion) for GET_TIME(double); needs <sys/time.h>
#define GET_TIME(now)\
{ struct timeval t; gettimeofday(&t, NULL);\
  now = t.tv_sec + t.tv_usec/1000000.0; }

unsigned Fibo_rec(int); //Recursive implementation of getting the nth Fibonacci number.
unsigned Fibo_iter(int); //Iterative implementation of getting the nth Fibonacci number.

int main()
{
  int pid, //process id
      status, //status code
      i, //loop counter
      n; //user input

  //Pipe descriptors
  int p1[2], //reading and writing between child1 and child2 for n
      p2[2], //reading and writing between child1 and child3 for n
      p3[2], //reading and writing between child1 and child2 for exec. time
      p4[2]; //reading and writing between child1 and child3 for exec. time
  pipe(p1);
  pipe(p2);
  pipe(p3);
  pipe(p4);

  double start, stop;
  double totalTime_rec, totalTime_iter;

  for(i = 1; i <= 3; i++) //creates 3 child processes
  {
    pid = fork();
    if(pid == 0 && i == 1) //child1
    {
      cout<<"in child1: Enter a number: ";
      cin>>n;

      write(p1[1], &n, sizeof(int)); //write n to pipe1 for Fibo_rec
      write(p2[1], &n, sizeof(int)); //write n to pipe2 for Fibo_iter

      read(p3[0], &totalTime_rec, sizeof(double)); //read exec. time for recursive way
      read(p4[0], &totalTime_iter, sizeof(double)); //read exec. time for iterative way

      cout<<"in child1: Time with recursion = "<<totalTime_rec<<" sec."<<endl;
      cout<<"in child1: Time with iteration = "<<totalTime_iter<<" sec."<<endl;
      exit(0);
    }

    else if(pid == 0 && i == 2) //child2
    {
      read(p1[0], &n, sizeof(int)); //read n from pipe1
      GET_TIME(start); //time checking starts
      cout<<"in child2: result = "<<Fibo_rec(n)<<endl; //compute nth fibonacci number using recursion
      GET_TIME(stop); //time checking ends
      totalTime_rec = stop - start; //get exec. time
      write(p3[1], &totalTime_rec, sizeof(double)); //write the exec. time to pipe3 for reading in child1
      exit(0);
    }
    else if(pid == 0 && i == 3) //child3
    {
      read(p2[0], &n, sizeof(int)); //read n from pipe2
      GET_TIME(start); //time checking starts
      cout<<"in child3: result = "<<Fibo_iter(n)<<endl; //compute nth fibonacci number using iteration
      GET_TIME(stop); //time checking ends
      totalTime_iter = stop - start; //get exec. time
      write(p4[1], &totalTime_iter, sizeof(double)); //write the exec. time to pipe4 for reading in child1
      exit(0);
    }
  }

  for(i = 1; i <= 3; i++) //wait for all child processes to finish
  {
    pid = wait(&status);
    cout<<"child (pid="<<pid<<") exited, status="<<status<<endl;
  }

  return 0;
}



//Finding the nth Fibonacci number using recursion
unsigned Fibo_rec(int n)
{
  if(n == 1 || n == 2) //base case
    return 1;
  else
    return Fibo_rec(n-1) + Fibo_rec(n-2);
}


//Finding the nth Fibonacci number using iteration
unsigned Fibo_iter(int n)
{
  int fib1 = 1; //1st Fibonacci number
  int fib2 = 1; //2nd Fibonacci number
  unsigned i, temp;
  if(n == 1 || n == 2) //base case
    return 1;
  else
  {
    for(i = 1; i < n-1; i++)
    {
      temp = fib1 + fib2; //compute the next Fibo number
      fib1 = fib2; //move fib1 to next Fibo number which is fib2
      fib2 = temp; //move fib2 to next Fibo number which is temp
    }
    return temp;
  }
}
