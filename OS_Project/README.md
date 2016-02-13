README

included:
	Source file – traffic_simul.cpp
	this README – README.md
	Project description – project_outline.pdf

Description:
	This project simulates an Intelligent Traffic Light System and Stop Sign Rule using multi threading. The driver of the program gives users two choices: run the simulation without limit and get the average response time.

Option 1: Run the simulation without limit
		The program will run without stopping and the only way to stop the 	execution is discussed below.
Option 2: Get the average response time
		The driver of the program will ask the user how many cars will be 	considered in the computation of average response time effectively limiting 	the execution of the program to that number of cars.

Requirements:
	Since threads were used in the implementation of the system, some compilation arguments must be added to run the program.
	To compile the program add -std=c++11 and -pthread to the compilation arguments.
			g++ filename.cpp -std=c++11 -pthread
	Since the user is given the option to run the program without stopping, the only way to terminate the program is by directly interrupting the execution. To do this press:
			ctrl + z or ctrl + c

