/**********************************************************
Juan Taruc: CSCI 176 Parallel matrix multiplication in java
This program performs matrix multiplication between matrix
A and B storing the results to matrix C using java threads
and synchronization variables. The program will also time 
the execution of the process base on the number of threads 
and the matrix dimensions.

	*** To run: ***
	> javac P3_matrix.java
	> java P3_matrix L m n p
		L, m, n are the matrix dimensions of A and B while
		p is the number of threads to run.

	*** Note ***
		The ouput of matrix C is run on 1920x1080 screen 
		resolution, so if this program is run in a screen
		that has lower resolution the output of the program
		wraps around the screen and therefore will be 
		unreadable.
**********************************************************/


//main class
public class P3_matrix {
	public static Globals matrix; //shared object Globals creation.
	public static Communication comm; //shared object Communication creation.
	public static int L, m, n, p;
	public static long start, stop;
	public static void main(String args[]) {
		matrix = new Globals(); //shared object memory allocation.
		comm = new Communication(); //shared object memory allocation.
		
		//get matrix dimensions from command line arguments + number of threads.
		L = Integer.parseInt(args[0]);
		m = Integer.parseInt(args[1]);
		n = Integer.parseInt(args[2]);
		p = Integer.parseInt(args[3]); //number of threads to run.
		
		Mythread[] thread = new Mythread[p];
		//initialize matrices with given dimensions
		init("A", L, m); //initlize matrix A with values i+j.
		init("B", m, n); //initlize matrix B with values i+j+1.
		init("C", L, n); //initlize matrix C with value 0.

		System.out.println("\n");
		System.out.println("L=" + L + ", m=" + m + ", n=" + n);
		start = System.currentTimeMillis();
		for(int i = 0; i < p; i++) {
			thread[i] = new Mythread(i, matrix, comm);
			thread[i].start();
		}
		for(int i = 0; i < p; i++) {
			try{thread[i].join();}
			catch(InterruptedException e){}
		}
		System.out.println("Total computation time(sec) = " + (double)(stop - start)/1000);
	}

	//matrix initializer
	public static void init(String flag, int row, int col) {
		int i, j; //loop counters.
		if (flag == "A") {
			matrix.A = new double[row][col];
			for(i = 0; i < row; i++) {
				for(j = 0; j < col; j++) {
					matrix.A[i][j] = i + j;
				}
			}
		}
		else if (flag == "B") {
			matrix.B = new double[row][col];
			for(i = 0; i < row; i++) {
				for(j = 0; j < col; j++) {
					matrix.B[i][j] = i + j + 1;
				}
			}
		}
		else if (flag == "C") {
			matrix.C = new double[row][col];
			//values for matrix C is 0s.
		}
	}
}


//
class Mythread extends Thread {
	static Globals matrix; //local for shared_obj
	static Communication comm; //local for shared_obj
	private int id;
	
	private P3_matrix p3_matrix = new P3_matrix();
	private int my_first_i, my_last_i, q, r;

	public Mythread(int k, Globals matrix, Communication comm) {
		////assign parameters and local here
		this.matrix = matrix;
		this.comm = comm;
		this.id = k;
		q = p3_matrix.L / p3_matrix.p;
		r = p3_matrix.L % p3_matrix.p;

		//my_first_i and my_last_i computation.
		if(id < r) {
			my_first_i = id * (q + 1);
			my_last_i = my_first_i + q;
		}
		else {
			my_first_i = id * q + r;
			my_last_i = my_first_i + (q - 1);
		}

		System.out.println("Thread_"+ id +": start_row="+ my_first_i+ ", end_row="+ my_last_i);
	}

	public void run() {
		int i, j, k; //loop counters.

		//matrix computation
		for(i = my_first_i; i <= my_last_i; i++) {
			for(j = 0; j < p3_matrix.n; j++) {
				for(k = 0; k < p3_matrix.m; k++)
					matrix.C[i][j] += matrix.A[i][k] * matrix.B[k][j];
			}
		}

		comm.increment(id);//wait or notify threads.

		//thread_0 prints matrix C after all threads are done.
		if(id == 0) {
			p3_matrix.stop = System.currentTimeMillis();

			System.out.println();
			System.out.println("Matrix C: first_20 x first_10");
			print_Cmatrix(0, 0);

			System.out.println();
			System.out.println("Matrix C: last_20 x last_10");
			print_Cmatrix(981, 1990);
		}
	}

	public void print_Cmatrix(int start_row, int start_col) {
		int i, j, k; //lopp counters.
		for(k = 0; k < 181; k++) {
			System.out.print("=");
		}
		System.out.println();

		for(i = start_row; i < start_row+20; i++) {
			for(j = start_col; j < start_col+10; j++)
				System.out.print(String.format("| %-15s ", Double.toString(p3_matrix.matrix.C[i][j])));
			System.out.println("|");
			for(k = 0; k < 181; k++) {
				System.out.print("=");
			}
			System.out.println();
		}
		System.out.println();
	}
}


//shared variables class
class Globals {
	public double[][] A; //matrix A
	public double[][] B; //matrix B
	public double[][] C; //matrix C
}


class Communication {
	private static int counter = 0;
	private static P3_matrix p3_matrix = new P3_matrix();
	public synchronized void increment(int id) {//mutex
		counter++;
		if(id == 0 && counter < p3_matrix.p) {
			try{wait();}
      		catch(InterruptedException e){}
		}
		else if (counter == p3_matrix.p) {
			notify();
		}
	}
}