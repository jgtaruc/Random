/*
Juan Taruc CSCI 176 Parallel Matrix Multiplication.
This program performs matrix multiplication between matrix
A and B storing the results to matrix C using pthreads. The
program will also time the execution of the process base on 
the number of threads and the matrix dimensions.
	
	*** To run: ./executable L m n p 
	where L m n are the matrix dimensions and p is the number 
	of threads that will be created.
*/
#include <iostream>
#include <pthread.h>
#include <cstdlib>
#include <sys/time.h>
#include <iomanip>

using namespace std;

//Given code in the examples in class for getting the current time.
//macro (in-line expansion) for GET_TIME(double); needs <sys/time.h>
#define GET_TIME(now)\
{ struct timeval t; gettimeofday(&t, NULL);\
  now = t.tv_sec + t.tv_usec/1000000.0; }

//globals.
double **A;
double **B;
double **C;
int L, m, n, //matrix dimensions
		p; //number of threads
int q, r; //for dividing the task to threads.
pthread_mutex_t mutex; //to protect cout

/* Matrix initializer.
This function initializes all three matrices. Since this
function is being reused by the three matrices we need some 
way to identify which matrix we want to initialize. The
- flag - argument handles this by its value of either
'A', 'B' or 'C' corresponding to the names of the matrices.
*/
void init_matrix(char flag, int row, int col) {
	int i, j;

	//initialize matrix A with values i+j.
	if(flag == 'A') {
		A = new double*[row];
		for(i = 0; i < row; i++) {
			A[i] = new double[col];
		}
		for(i = 0; i < row; i++) {
			for(j = 0; j < col; j++) {
				A[i][j] = i + j;
			}
		}
	}

	//initialize matrix B with values i+j+1.
	else if(flag == 'B') {
		B = new double*[row];
		for(i = 0; i < row; i++) {
			B[i] = new double[col];
		}
		for(i = 0; i < row; i++) {
			for(j = 0; j < col; j++) {
				B[i][j] = i + j + 1;
			}
		}
	}

	//initialize matrix C.
	else if(flag == 'C') {
		C = new double*[row];
		for(i = 0; i < row; i++) {
			C[i] = new double[col];
		}
	}
}


/* Matrix Multiplication
This function is invoked by the threads to calculate
the dot product of matrix A and B and stores the results to
matrix C. 
*/
void *matrix_multiply(void *thread_id) {
	long my_rank = (long)thread_id; //get and convert the thread_id and store
								//store it to my_rank
	int my_first_i, my_last_i; //used to assign the values each thread will work on.
	int i, j, k; //loop counters

	//hw1
	if(my_rank < r) {
		my_first_i = my_rank * (q + 1);
		my_last_i = my_first_i + q;
	}
	else {
		my_first_i = my_rank * q + r;
		my_last_i = my_first_i + (q - 1);
	}

	pthread_mutex_lock(&mutex); //try to acquire mutex 
	cout<<"Thread_"<<my_rank<<": start_row="<<my_first_i<<", end_row="<<my_last_i<<endl;
	pthread_mutex_unlock(&mutex);

	//perform matrix multiplication with my_first_i and my_last_i.
	for(i = my_first_i; i <= my_last_i; i++) {
		for(j = 0; j < n; j++) {
			for(k = 0; k < m; k++)
				C[i][j] += A[i][k] * B[k][j];
		}
	}
	return NULL;
}


/* print matrix C
This function prints matrix C with its first and last 20x10 values.
*/
void print(int start_row, int start_col) {
	int i, j, b; //loop counters

	for(b = 0; b < 181; b++) { //print the borders
		cout<<"=";
	}
	cout<<endl;

	for(i = start_row; i < start_row+20; i++) {
		for(j = start_col; j < start_col+10; j++) {
			cout<<"| "<<setw(15)<<C[i][j]<<" "; //print the divisions on the left side of each value.
		}
		cout<<"|"<<endl; //print the divisions on the last value of each rows.
		for(b = 0; b < 181; b++) { //print the borders.
		cout<<"=";
		}
		cout<<endl;
	}
}


//main function.
int main(int argc, char* argv[])
{
	int i; //loop counter.
	long thread_id; //thread id.

	//for timing exec time.
	double start, stop, totalTime;

	cout<<endl;
	//if one of the cmd line arguments is missing, terminate the 
	//program.
	if(!argv[1] or !argv[2] or !argv[3] or !argv[4]) {
		cout<<"Missing command line arguments: be sure it is L m n p"<<endl;
		exit(0);
	}
	//save the command line arguments to its respective variables.
	L = atoi(argv[1]);
	m = atoi(argv[2]);
	n = atoi(argv[3]);
	p = atoi(argv[4]);

	//initialize the three matrices.
	init_matrix('A', L, m);
	init_matrix('B', m, n);
	init_matrix('C', L, n);

	//reference for assigning the respective data each thread
	//will work on.
	q = L / p;
	r = L % p;

	cout<<"L = "<<L<<", m = "<<m<<", n = "<<n<<endl;

	pthread_mutex_init(&mutex, NULL); //initialize mutex

	//create threads passing multiply function.
	pthread_t threads[p];
	GET_TIME(start); //time checking starts
	for(thread_id = 0; thread_id < p; thread_id++) {
		pthread_create(&threads[thread_id], NULL, matrix_multiply, (void*)thread_id);
	}

	//wait for the threads to finish
	for(thread_id = 0; thread_id < p; thread_id++) {
		pthread_join(threads[thread_id], NULL);
	}
	GET_TIME(stop); //time checking ends

	//print the results.
	totalTime = stop - start;
	cout<<endl;
	cout<<"Total computation time(sec) = "<<totalTime<<endl;
	cout<<endl;

	cout<<"Matrix C: first_20 x first_10"<<endl;
	print(0, 0);
	
	cout<<endl;

	cout<<"Matrix C: last_20 x last_10"<<endl;
	print(981, 1990);
	return 0;	
}