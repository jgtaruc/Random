//Juan Taruc CSCI 156 Project - TCP client w/ 2hop relay
//This program will continually send tcp packets
//to the server via 2hop relay scenario for the specified send time.
//It will then divide the total bits sent to the
//server by the send time to get the attempted
//data rate of the communication.

import java.io.*;
import java.net.*;
import java.util.Timer;
import java.util.concurrent.locks.*;
import java.util.TimerTask;


class TCP_Client
{
	private static int num_bytes = 256; //the packet size.
	private static int sending_period = 30000; //the specified time the client will send
										      //packets to the server in milliseconds.
	private static int sending_interval = 10; //delay for each packet sent in milliseconds.
	public static int total_bits_sent = 0; //the total bits sent during the send time.
	public static byte[] sendData = new byte[num_bytes]; //data that is going to be sent to the server.
	public static int flag = 1;

	public static void main(String argv[]) throws Exception
	{
		//fill the byte buffer sendData with dummy data.
		for(int i=0; i<sendData.length; i++)
			sendData[i] = 10;

		final ReentrantLock lck = new ReentrantLock(); //synchronization variable.

		//Create a timer that will execute the sending of packets in the specified
		//send time and also send those packets in an interval.
		Timer timer = new Timer(true);
		timer.scheduleAtFixedRate(new TimerTask()
		{
			@Override
			public void run()
			{
					try
					{
						Socket clientSocket = new Socket("juan-Inspiron-7548", 1234); //create a socket for sending.
						DataOutputStream outToServer = 
							new DataOutputStream(clientSocket.getOutputStream()); //create a stream for writing the
																			      //data to the server socket.
						//guard for sending the packet if the elapsed time is already reached.
						lck.lock();
						outToServer.write(sendData); //write the data to the server socket.
						total_bits_sent = total_bits_sent + (outToServer.size()*8); //calculate the total bits sent.
						outToServer.close(); //close the data stream.
						lck.unlock();					
					}
					catch(Exception e){}
			}
		},0,sending_interval); //schedule the sending for repeated fixed-rate
							   //execution(sending_interval) after
							   //the data, socket and packet is successfully created.
		Thread.sleep(sending_period); //stop sending packets after the send time specified is reached.
		lck.lock();
		timer.cancel(); //after the send time has elapsed cancel the remaining task 
						//scheduled by the timer.
		timer.purge();
		int attempted_data_rate = total_bits_sent / (sending_period/1000); //calculate the attempted data rate by 
																	  	   //dividing the total bits sent by the
																	       //send time in seconds
		System.out.println("Attempted data rate: " + attempted_data_rate + " bits/s"); //show the result.
	}
}