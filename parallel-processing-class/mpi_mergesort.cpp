/*
Juan Taruc CSCI 176 MPI Parallel Merge Sort:
	This program takes all inputs from the command line.
	
	To compile/run:
		> mpic++ mpi_mergesort.cpp
		> mpiexec -n P ./a.out N
		P = communication size, N = problem size

	This program creates dynamic arrays each of which has size
	N/P with random numbers from 0 ~ 99 as its values. It then
	prints the unsorted arrays in order by rank then uses, sorts
	it using qsort then uses tree reduction for mergin operation. 
	Finally, it prints the global array with size N in sorted order.
*/

#include <iostream>
#include <stdlib.h> //for qsort, srand, rand
#include <mpi.h>  

using namespace std;

//compare function is use by qsort to compare two elements
int compare(const void *a_p, const void *b_p){
	int a = *(int *)a_p;
	int b = *(int *)b_p;
	return a - b;
}

//merge function merges two arrays in order producing an array
//double the size of the two arrays.
//size argument is already double the size of the two arrays a and b.
int *merge(int *a, int *b, int size){
	//temp will contain the sorted elements of the two arrays
	int *temp = new int[size];
	int i, j, k;
	i = j = k = 0;
	for(k = 0; k < size; k++){
		//left sublist exhausted? then copy remaining values of b array to temp
		if(i == size/2){
			temp[k] = b[j++];
		}
		//right sublist exhausted? then copy the remaining values of a array to temp
		else if(j == size/2){
			temp[k] = a[i++];
		}
		else{
			//compare each value of array a and b and put to temp array in order.
			if(a[i] <= b[j])
				temp[k] = a[i++];
			else
				temp[k] = b[j++];
		}
	}
	return temp;
}

int main(int argc, char* argv[]) {
	int n,         //problem size
		comm_size, //communication size
		my_rank,   //id of each process
		root = 0;  //root is process_0
	int local_n; 
	int *local_array; //array for each process
	int i; //loop counter

	MPI_Init(NULL, NULL); //start up MPI
	MPI_Comm_size(MPI_COMM_WORLD, &comm_size);//get number of process from -n x
	MPI_Comm_rank(MPI_COMM_WORLD, &my_rank); //get my_rank among all

	//1. process_0 accesses command line argument for n(problem size)
	//and broadcasts it to all other processes
	if(my_rank == root){
		n = atoi(argv[1]); //get size of array
			
	}
	//process 0(root) broadcasts n to all other processes
	MPI_Bcast(&n, 1, MPI_INT, root, MPI_COMM_WORLD);
	
	//2. Each process creates a dynamic array of integers with n/comm_size elements and
	//fills it up with random numbers from 0 ~ 99.
	local_n = n / comm_size;
	local_array = new int[local_n];
	srand(my_rank + 1);
	for(i = 0; i < local_n; i++){
		local_array[i] = rand() % 100;
	}

	//3. process 0 collects local array of all processes and prints them in order by rank_id.
	if(my_rank != root)
		MPI_Send(&local_array[0], local_n, MPI_INT, root, 0, MPI_COMM_WORLD);

	if(my_rank == root){
		cout<<"process_0 unsorted array: ";
		for(i = 0; i < local_n; i++)
			cout<<local_array[i]<<" ";
		cout<<endl;

		for(i = 1; i < comm_size; i++) {
			int *rcv_array = new int[local_n];
			MPI_Recv(&rcv_array[0], local_n, MPI_INT, i, 0, MPI_COMM_WORLD, MPI_STATUS_IGNORE);
			cout<<"process_"<<i<<" unsorted array: ";
			for(int j = 0; j < local_n; j++)
				cout<<rcv_array[j]<<" ";
			cout<<endl;
		}	
	}
	
	//4. Each process performs qsort to sort their local arrays in ascending order.
	qsort(local_array, local_n, sizeof(int), compare);


	//5. Using tree reduction for merging operation using method in HW1(#3).
	int divisor = 2;
	int core_diff = 1;
	int q = local_n;
	int r = n % comm_size;
	int *global_array = local_array;
	int partner;
	bool done = false;
	int size = local_n;

	while(divisor <= comm_size && !done){
		if(my_rank % divisor == 0){
			partner = my_rank + core_diff;
			//for odd number of processes
			if(partner <= (comm_size - 1)){
				int *rcv_array = new int[size];
				MPI_Recv(&rcv_array[0], size, MPI_INT, partner, 0, MPI_COMM_WORLD, MPI_STATUS_IGNORE);
				size *= 2; //double the size since merging two arrays
				global_array = merge(global_array, rcv_array, size);
			}
		}
		else{
			partner = my_rank - core_diff;
			MPI_Send(&global_array[0], size, MPI_INT, partner, 0, MPI_COMM_WORLD);
			done = true;
		}
		divisor *= 2;
		core_diff *= 2;
	}

	//6. process 0 prints the global sorted array
	if(my_rank == root){
		cout<<endl;
		cout<<"process_0 global sorted array: "<<endl;
		for(i = 0; i < size; i++){
			cout<<global_array[i]<<" ";
		}
		cout<<endl;
	}

	MPI_Finalize(); //shut down MPI
	return 0;
}