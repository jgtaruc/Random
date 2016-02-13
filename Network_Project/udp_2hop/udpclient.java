//Juan Taruc CSCI 156 Project - UPD client w/ 2 hop relay
//This program will continually send udp packets
//to the server via 2hop relay scenario for the specified send time.
//It will then divide the total bits sent to the
//server by the send time to get the attempted
//data rate of the communication.

import java.io.*;
import java.net.*;
import java.util.concurrent.locks.*;
import java.util.Timer;
import java.util.TimerTask;


class UDP_Client
{
	private static int num_bytes = 256; //the packet size.
	private static int send_time = 30000; //the specified time the client will send
										  //packets to the server in milliseconds.
	private static int sending_interval = 10; //delay for each packet sent in milliseconds.
	public static int total_bits_sent = 0; //the total bits sent during the send time.
	public static byte[] sendData = new byte[num_bytes]; //data that is going to be sent to the server.
	

	//main function generates the data that will be sent to the server, creates a socket for sending,
	//schedules the sending of that data in the specified send time and calculates the attempted data 
	//rate of the communication.
	public static void main(String args[]) throws Exception
	{
		//fill the byte buffer sendData with dummy data.
		for(int i=0; i<sendData.length; i++)
			sendData[i] = 10;

		final ReentrantLock lck = new ReentrantLock(); //synchronization variable.

		//clientSocket, IPAddress, sendPacket needs to be declared final since it is being
		//accessed within innerclass Timer
		final DatagramSocket clientSocket = new DatagramSocket(); //create a socket for sending
		final InetAddress IPAddress = InetAddress.getByName("localhost"); //assign the address where the data
																			       //will be sent.
		//create a packet that will be sent from the clientSocket with the data generated
		//and with the specified address and port number.
		final DatagramPacket sendPacket = new DatagramPacket(sendData, sendData.length, IPAddress, 1234); 
		//Create a timer that will execute the sending of packets in the specified
		//send time and also send those packets in an interval.
		Timer timer = new Timer(true);
		final long finishtime = System.currentTimeMillis() + send_time; //calculate the end time of the program
		timer.scheduleAtFixedRate(new TimerTask()
		{
			@Override
			public void run()
			{
					try
					{
						//guard for sending the packet if the elapsed time is already reached.
						lck.lock();
						//repeatedly send the packet to the server.	
						clientSocket.send(sendPacket);
						//calculate the total bits sent by the socket.
						total_bits_sent += (sendPacket.getLength()*8);
						lck.unlock();
					}
					catch(Exception e){}
			}
		},0,sending_interval); //schedule the sending for repeated fixed-rate
							   //execution(sending_interval) after
							   //the data, socket and packet is successfully created.
		Thread.sleep(send_time); //stop sending packets after the send time specified is reached.
		lck.lock();
		timer.cancel(); //after the send time has elapsed cancel the remaining task 
						//scheduled by the timer.
		timer.purge();
		int attempted_data_rate = total_bits_sent / (send_time/1000); //calculate the attempted data rate by 
																	  //dividing the total bits sent by the
																	  //send time in seconds
		System.out.println("Attempted data rate: " + attempted_data_rate + " bits/s"); //show the result.
	}
} 
