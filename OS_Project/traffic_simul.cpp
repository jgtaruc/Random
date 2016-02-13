/*
Juan Taruc CSCI 144 Fall 2015 Project
This project is a simulation of an intelligent traffic
light system with multithreading. This program generates
random car arriving from 4 directions(North,East,South,West)
with the probability of generating a car 50% in each time.
This system stores the RTP(Request To Pass) by cars when
it gets in a certain distance from the intersection, which
is 50 meters. It allows a car to pass without stopping if
its RTP is in front of the Queue. Also, this system allows
N cars to pass in same/opposite direction without waiting 
for each other. The system computes the response time for
each car by determining each car's wait time to cross the
intersection.

Given:
	*Each car takes 3 seconds to cross the intersection
	*The critical distance from the intersection is 50 meters
	*Each car generated has an random speed of 45~55mph
	*All the time related variables in the project is measured 
		in seconds(waittime,sleep)
	*Generating a car has the probability of 50%
	*Simulation does not include turns
	
*/

//To run, g++ <filename> -std=c++11 -pthread

#include <iostream>
#include <string> 
#include <stdlib.h> //for rand()
#include <sstream> //for int to string conversion w/o c++11
#include <time.h> //for time
#include <unistd.h> //for sleep()
#include <mutex> //for mutual exclusion(lock)
#include <thread>//for threads


using namespace std;

int car_id = 1; //car id starts with "Car 1"
mutex lck;	//the synchronization variable
double chance = 0.5; //the probability that a car is arriving in N,E,W,S direction
int critical_distance = 50; //the distance where the sensor will accept the rtp of each car.
int num_cars; //number of cars that will be used in comparing the average response time
			  //of the system and stop sign rule. 
int num_cars_i; //num_cars_i is use to track the number of cars serviced.
int total_rt = 0; //total response time will be use to calculate the average response time.
int traffic_system_flag; //1 for intelligent traffic light system
						 //0 for stop sign rule


//Car object - A car object is created in 4 directions with the percentage chance
struct Car
{
	string id; //Car i
	time_t arrival_time; //arrival time in ms
	string direction;	//North | East | South | West
	int speed;			//random number between 45 ~ 55
	Car(string i, string d) //Car constructor
	{
		id = i; //set id
		direction = d; //set direction
		arrival_time = time(NULL); //set arrival time to now
		speed = rand() % 10 + 45;//45~55 mph //random speed
	}
	string Get_Id()//return private Car Id
	{	return id;	}
	int Get_ArriveTime() //return private arrival time
	{	return arrival_time; }
	string Get_Direction() //return private direction
	{	return direction;	}
	int Get_Speed() //return private speed
	{	return speed;	}
};


//Queue object - the shared queue between all the threads in the program. It holds
//which Car object is allowed to pass in the intersection.
//The template class in the Queue object is not required to run the program, but
//I implemented it that way	as an exercise in polymorphism.
template <class T>
class Queue
{
public:
	Queue() //Queue constructor
	{
		//Initially pointers to front and rear of the Queue
		//is NULL.
		front = NULL;
		rear = NULL;
	}

	bool Empty() //Empty Queue?
	{
		return (front == NULL);
	}

	void Enqueue(T *c); //pushes a Car object to the Queue

	void Dequeue();//pops a Car object from the queue

	T *Get_Front() //returns whichever Car object is front of the Queue
	{	return front->obj;	}

	T *Get_Next() //return whichever Car Object is next to the front of Queue
	{	return front->next->obj;	}

private:
	struct Node //The node structure fills the Queue
	{
		T *obj; //obj holds the Car object
		Node *next; //pointer to next node
	};
	Node *front; //first Node of Queue
	Node *rear; //last Node of Queue
};


//Enqueue - pushes the new Car object to the Queue
//It also prints that a car is arriving in a direction and speed
template <class T>
void Queue<T>::Enqueue(T *c)
{
	Node *newnode = new Node;
	newnode->obj = c; 
	newnode->next = NULL;
	if(front == NULL) //pushing to an empty queue
	{	front = newnode;	}
	else //always push the new node to the rear of the non empty Queue
	{	rear->next = newnode;	}
	rear = newnode;
	//Car ID arriving from direction __ at __ mph
	cout<<"+ "<<newnode->obj->id<<" arriving from direction ";
	cout<<newnode->obj->direction<<" at "<<newnode->obj->Get_Speed()<<"mph"<<endl;
}


//Dequeue - pops a node which contains the Car object from the Queue
//It always pops the front of the Queue
template <class T>
void Queue<T>::Dequeue()
{
	if(Empty()) //if the Queue is empty don't Dequeue
		return;
	Node *temp = new Node; //temporary Node that will hold the Car object
	temp = front; //save the front to temp
	front = front->next; //The new front of the Queue is the next object
	temp->next = NULL; //prepare for deleting
	delete temp; //delete the front Car object of the Queue
}


//Generate_CarId - takes i as an argument, converts it to string
//and appends it to the string "Car " to make the Car id.
string Generate_CarId(int i)
{
	stringstream convert; //convert int to string w/o using c++11 to_string
						//since it is not compatible with some systems.
	convert << i;
	return "Car "+ convert.str();
}


//Get_Probability - returns a random number in the range 0 ~ 1(0% ~ 100%)
double Get_Probability()
{
	srand(time(NULL));
	return (double)rand() / RAND_MAX;//generate random number between 0~1
}


int Max(int x, int y)
{
	if(x > y)
		return x;
	return y;
}


//Converts mph to meter per second
double Convert_to_MeterPerSec(int x)
{
	//Conversion factor 0.44704
	return x * 0.44704;
}


//Calculates the average response time of the Intelligent Traffic Light
//System or Stop Sign Rule. After the number of cars specified had been
//all serviced, this function prints out the average response time and
//exits the program.
void Ave_RT(int rt)
{
	num_cars_i--;
	total_rt += rt;
	if(num_cars_i == 0)
	{
		cout<<endl;
		if(traffic_system_flag == 1)
			cout<<"Intelligent Traffic Light System"<<endl;
		else if(traffic_system_flag == 0)
			cout<<"Stop Sign Rule"<<endl;
		cout<<"Average Response Time for "<<num_cars<<" cars: "<<total_rt/num_cars;
		cout<<" seconds"<<endl;
		exit(0);
	}
}


//Drive_ITS - this implementation uses the Intelligent Traffic light system.
//It allows N cars to pass in one direction and also allows Cars in
//opposite direction to pass without waiting.
//***With some modification in the code of Drive function, this program
//can also simulate stop sign rule. Follow directions below.
void Drive_ITS(Queue<Car> &q)
{	
	
	if(!q.Empty())//if the Queue is not empty, safe to Dequeue
	{
		Car *prev_car = q.Get_Front();//save the Car in front of queue
		q.Dequeue(); //pop the car in front

		//calculate the response time.
		int t = critical_distance / Convert_to_MeterPerSec(prev_car->Get_Speed()); // D/S
		int rt = Max(t, difftime(time(NULL), prev_car->Get_ArriveTime()))+3;

		//Print the results
		cout<<"	- "<<prev_car->Get_Id()<<" driving from "<<prev_car->Get_Direction()<<" | Response Time: "<<rt<<" seconds"<<endl;


		Ave_RT(rt);


		//if Queue is not empty, check the next car if it can pass
		//* if system is going to use stop sign rule every car must wait before its turn regardless of their 
		//direction, so to achieve that need to comment below code to only allow 1 car to pass per turn.
		if(!q.Empty())
		{
			Car *next_car = q.Get_Front(); //get the new Car in front
			if(prev_car->direction == next_car->direction) //same direction, let it pass w/o waiting
				Drive_ITS(q);
			//opposite directions, let it pass w/o waiting
			if(((prev_car->direction == "North") && (next_car->direction == "South")) || 
				((prev_car->direction == "South") && (next_car->direction == "North")))
				Drive_ITS(q);
			if(((prev_car->direction == "East") && (next_car->direction == "West")) || 
				((prev_car->direction == "West") && (next_car->direction == "East")))
				Drive_ITS(q);
		}
	}
}


//Drive_SSR - this implementation uses the Stop Sign Rule.
//It allows N cars to pass in one direction and also allows Cars in
//opposite direction to pass without waiting.
//***With some modification in the code of Drive function, this program
//can also simulate stop sign rule. Follow directions below.
void Drive_SSR(Queue<Car> &q)
{	
	
	if(!q.Empty())//if the Queue is not empty, safe to Dequeue
	{
		Car *prev_car = q.Get_Front();//save the Car in front of queue
		q.Dequeue(); //pop the car in front

		//calculate the response time.
		int t = critical_distance / Convert_to_MeterPerSec(prev_car->Get_Speed()); // D/S
		int rt = Max((t*2), difftime(time(NULL), prev_car->Get_ArriveTime()))+3;
		//Print the results
		cout<<"	- "<<prev_car->Get_Id()<<" driving from "<<prev_car->Get_Direction()<<" | Response Time: "<<rt<<" seconds"<<endl;


		Ave_RT(rt); 


	}

}


//Intersection_Thread - handles the intersection with each call to Drive
void Intersection_Thread(Queue<Car> &q)
{
	while(true)
	{
		lck.lock(); //get the lock because Drive() pops a Car in the shared object Queue
		if(!q.Empty())//If non empty Queue, safe to pop from the queue
		{
			cout<<endl<<"============================= GO ==============================="<<endl;
			if(traffic_system_flag == 1)
				Drive_ITS(q); //pop a Car from the Queue
			else if(traffic_system_flag == 0)
				Drive_SSR(q);
			sleep(3); //wait for the car to pass the intersection
			cout<<"============================ STOP =============================="<<endl<<endl;
		}
		lck.unlock();//release the lock
	}
}


//North, East, West, South Threads pushes a new Car object to the Queue with the appropriate
//direction.
//It pushes a Car object if the probability <= chance(%)
//One thread handles each direction so the generation of N Cars is completely random best. 
void North_Thread(Queue<Car> &q)
{
	while(true)
	{
		if(Get_Probability() <= chance)//generate car with probability of 50%
		{
			lck.lock();
			q.Enqueue(new Car(Generate_CarId(car_id++), "North"));
			lck.unlock();
		}
		usleep(5);
	}
}

void East_Thread(Queue<Car> &q)
{
	while(true)
	{
		if(Get_Probability() <= chance)//generate car with probability of 50%
		{
			lck.lock();
			q.Enqueue(new Car(Generate_CarId(car_id++), "East"));
			lck.unlock();
		}
		usleep(5);
	}
}

void South_Thread(Queue<Car> &q)
{
	while(true)
	{
		if(Get_Probability() <= chance)//generate car with probability of 50%
		{
			lck.lock();
			q.Enqueue(new Car(Generate_CarId(car_id++), "South"));
			lck.unlock();
		}
		usleep(5);
	}
}

void West_Thread(Queue<Car> &q)
{
	while(true)
	{
		if(Get_Probability() <= chance)//generate car with probability of 50%
		{
			lck.lock();
			q.Enqueue(new Car(Generate_CarId(car_id++), "West"));
			lck.unlock();
		}
		usleep(5);
	}
}



int main()
{
	Queue<Car> queue; //declare a Queue instance
	srand(time(NULL));//seed the random number generator

	//Driver program
	string input;
	cout<<"Options:"<<endl;
	cout<<"1 - Run the program without limit."<<endl;
	cout<<"2 - Measure the average response time."<<endl; 
	cout<<"q - Exit."<<endl;
	cout<<">> "; cin>>input;
	cout<<endl;
	if(input == "1") //Run intelligent traffic light / stop sign simulation w/o limit
	{
		cout<<"Which traffic system would you like to use?"<<endl;
		cout<<"1 - Intelligent Traffic Light System"<<endl;
		cout<<"2 - Stop Sign Rule"<<endl;
		cout<<"q - Exit."<<endl;
		cout<<">> "; cin>>input;
		cout<<endl;

		if(input == "1")
			traffic_system_flag = 1;
		else if(input == "2")
			traffic_system_flag = 0;
		else if(input == "q")
			exit(0);
		else
		{
			cin.clear();
			cout<<"*** Input not recognize! ***"<<endl;
			main();
		}

		cout<<"Press ctrl-z or ctrl-c to exit."<<endl;
		cout<<"Starting in..."<<endl;
		for(int i=3; i>0; i--)
		{
			cout<<i<<endl;
			sleep(1);
		}

		//start the threads
		thread thread1(North_Thread, ref(queue));
		thread thread2(East_Thread, ref(queue));
		thread thread3(South_Thread, ref(queue));
		thread thread4(West_Thread, ref(queue));
		thread intersection(Intersection_Thread, ref(queue));

		//wait for each thread to finish
		thread1.join();
		thread2.join();
		thread3.join();
		thread4.join();
		intersection.join();
	}
	else if(input == "2") //measure the average response time of traffic system chosen
	{
		int n;
		cout<<"Enter the number of cars: "; //how many number of cars to consider
		cin>>n;
		cout<<endl;
		if(cin)
			num_cars = num_cars_i = n;
		else
		{
			cout<<"*** Invalid input! ***"<<endl;
			main();
		}

		cout<<"Which traffic system would you like to use?"<<endl;
		cout<<"1 - Intelligent Traffic Light System"<<endl;
		cout<<"2 - Stop Sign Rule"<<endl;
		cout<<"q - Exit."<<endl;
		cout<<">> "; cin>>input;
		cout<<endl;

		if(input == "1")
			traffic_system_flag = 1;
		else if(input == "2")
			traffic_system_flag = 0;
		else if(input == "q")
			exit(0);
		else
		{
			cin.clear();
			cout<<"*** Input not recognize! ***"<<endl;
			main();
		}

		cout<<"Press ctrl-z or ctrl-c to exit."<<endl;
		cout<<"Starting in..."<<endl;
		for(int i=3; i>0; i--)
		{
			cout<<i<<endl;
			sleep(1);
		}

		//start the threads
		thread thread1(North_Thread, ref(queue));
		thread thread2(East_Thread, ref(queue));
		thread thread3(South_Thread, ref(queue));
		thread thread4(West_Thread, ref(queue));
		thread intersection(Intersection_Thread, ref(queue));

		//wait for each thread to finish
		thread1.join();
		thread2.join();
		thread3.join();
		thread4.join();
		intersection.join();
	}
	else if(input == "q")
		exit(0);
	else
	{
		cin.clear();
		cout<<"*** Input not recognize! ***"<<endl;
		main();
	}

	return 0;
}
