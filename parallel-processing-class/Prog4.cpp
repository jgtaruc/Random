/*
Juan Taruc CSCI 176 OpenMP Count Sort:
	This program takes p(number of processors) and 
	n(size of array) as command line arguments then
	creates 3 arrays with random values from 0 ~ n-1.
	It then sorts the 3 arrays respectively:
		array1 via serial_count_sort
		array2 via parallel_count_sort1
		array3 via parallel_count_sort2
	The program will show the unsorted and sorted
	arrays, execution time, and if the array is sorted
	properly for each sorting function.
*/



#include <iostream>
#include <stdlib.h>
#include <omp.h>
#include <sys/time.h>

using namespace std;

#define GET_TIME(now)\
{ struct timeval t; gettimeofday(&t, NULL);\
  now = t.tv_sec + t.tv_usec/1000000.0; }


//Globals
int p, //number of threads
	n; //size of array


/*
function initArray(array size)
	creates an array of integers with size n
	and initializes that array with random 
	values < n.
*/
int* initArray(int n, int* values){
	int *arr = new int[n]; //creates size n array of ints 
	
	//initialize array values
	for(int i = 0; i < n; i++){
		arr[i] = values[i];
	}

	return arr;
}


/*
function randomValues(array size)
	generates random numbers < n, stores those
	values to an array then returns that array.
	The returned array will be used in the 
	initialization process of the three arrays
	for count sort since those arrays need to have
	identical random values.
*/
int* randomValues(int n){
	int *arr = new int[n];

	for(int i = 0; i < n; i++){
		arr[i] = rand() % n;
	}

	return arr;
}


/*
function serial_count_sort(array, size)
	count sort algorithm derived
	from the assignment handout.
*/
void serial_count_sort(int a[], int n){
	int i, j, count; 
	int *temp = new int[n];

	for(i = 0; i < n; i++){
		count = 0;
		for(j = 0; j < n; j++){
			if(a[j] < a[i]) count++;
			else if(a[j] == a[i] && j < i) count++;
		}
		temp[count] = a[i];
	}

	for(i = 0; i < n; i++){
		a[i] = temp[i];
	}
	free(temp);
}


/*
function parallel_count_sort1(array, size)
	count sort algorithm derived from
	the assignment using openmp to parallelize
	the for i loop.
*/
void parallel_count_sort1(int a[], int n){
	int i, j, k, count; 
	int *temp = new int[n];

	# pragma omp parallel for num_threads(p) \
		default(none) shared(a, n, temp) private(i, j, count)
	for(i = 0; i < n; i++){
		count = 0;
		for(j = 0; j < n; j++){
			if(a[j] < a[i]) 
				count++;
			else if(a[j] == a[i] && j < i) 
				count++;
		}		
		temp[count] = a[i];
	}

	for(k = 0; k < n; k++){
		a[k] = temp[k];
	}
	free(temp);
}


/*
function parallel_count_sort2(array, size)
	variation of parallel_count_sort1 function but
	overwriting the array with the temporary array
	is parallel.
*/
void parallel_count_sort2(int a[], int n){
	int i, j, k, count; 
	int *temp = new int[n];

	# pragma omp parallel for num_threads(p) \
		default(none) shared(a, n, temp) private(i, j, count)
	for(i = 0; i < n; i++){
		count = 0;
		for(j = 0; j < n; j++){
			if(a[j] < a[i]) count++;
			else if(a[j] == a[i] && j < i) count++;
		}
		temp[count] = a[i];
	}

	#pragma omp parallel for num_threads(p) \
		default(none) shared(a, n, temp) private(k)
	for(k = 0; k < n; k++){
		a[k] = temp[k];
	}
	free(temp);
}


/*
function checkArray(array, size)
	checks the array if it is sorted properly.
*/
bool checkArray(int a[], int n){
	int i;
	for(i = 0; i < n; i++){
		if(a[i] > a[i+1] && a[i] != a[i+1] && (i+1) < n) 
			return false;
	}
	return true;
}


/*
function display(array, size, range)
	display the first(range) and 
	last(range) of array.
	in this assignment displaying the 
	first20 and last20 elements of array.
*/
void display(int a[], int n, int range){
	int i;

	cout<<"First "<<range<<": ";
	for(i = 0; i < range; i++)
		cout<<a[i]<<" ";
	cout<<endl<<endl;

	cout<<"Last  "<<range<<": ";
	for(i = n-range; i < n; i++)
		cout<<a[i]<<" ";
	cout<<endl;

}


int main(int argc, char* argv[]){
	if(!argv[1] or !argv[2]){
		cout<<"Needs command line arguments p n."<<endl;
		exit(0);
	}

	double start, stop;
	p = atoi(argv[1]); //get number of threads
	n = atoi(argv[2]); //get size of array

	//create 3 identical arrays with size n
	int *values = randomValues(n);
	int *arr1 = initArray(n, values);
	int *arr2 = initArray(n, values);
	int *arr3 = initArray(n, values);
	free(values);


	//serial count sort
	cout<<endl<<endl;
	cout<<"// Serial Count Sort"<<endl;
	//display unsorted array1(first20 ~ last20)
	cout<<"****************** Unsorted array1 **********************"<<endl;
	display(arr1, n, 20);
	cout<<"*********************************************************"<<endl;
	GET_TIME(start);
	serial_count_sort(arr1, n);
	GET_TIME(stop);
	//display sorted array1(first20 ~ last20)
	cout<<"******************* Sorted array1 ***********************"<<endl;
	display(arr1, n, 20);
	cout<<"*********************************************************"<<endl;
	cout<<"Exec time: "<<stop-start<<"(sec)"<<endl;
	if(checkArray(arr1, n)) cout<<"Sorting successful"<<endl;
	else cout<<"Sorting unsuccessful"<<endl;
	cout<<"// Done Serial Count Sort"<<endl;
	cout<<endl<<endl;


	//parallel count sort 1
	cout<<endl<<endl;
	cout<<"// Parallel Count Sort 1"<<endl;
	//display unsorted array2(first20 ~ last20)
	cout<<"****************** Unsorted array2 **********************"<<endl;
	display(arr2, n, 20);
	cout<<"*********************************************************"<<endl;
	GET_TIME(start);
	parallel_count_sort1(arr2, n);
	GET_TIME(stop);
	//display sorted array2(first20 ~ last20)
	cout<<"******************* Sorted array2***********************"<<endl;
	display(arr2, n, 20);
	cout<<"*********************************************************"<<endl;
	cout<<"Exec time: "<<stop-start<<"(sec)"<<endl;
	if(checkArray(arr2, n)) cout<<"Sorting successful"<<endl;
	else cout<<"Sorting unsuccessful"<<endl;
	cout<<"// Done parallel Count Sort 1"<<endl;
	cout<<endl<<endl;


	cout<<endl<<endl;
	cout<<"// Parallel Count Sort 2"<<endl;
	//parallel count sort 2
	//display unsorted array3(first20 ~ last20)
	cout<<"****************** Unsorted array3 **********************"<<endl;
	display(arr3, n, 20);
	cout<<"*********************************************************"<<endl;
	GET_TIME(start);
	parallel_count_sort1(arr3, n);
	GET_TIME(stop);
	//display sorted array3(first20 ~ last20)
	cout<<"******************* Sorted array3***********************"<<endl;
	display(arr3, n, 20);
	cout<<"*********************************************************"<<endl;
	cout<<"Exec time: "<<stop-start<<"(sec)"<<endl;
	if(checkArray(arr3, n)) cout<<"Sorting successful"<<endl;
	else cout<<"Sorting unsuccessful"<<endl;
	cout<<"// Done parallel Count Sort 2"<<endl;
	cout<<endl<<endl;

 	return 0;
}


