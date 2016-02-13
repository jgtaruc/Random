//Juan Taruc CSCI 156 Project - TCP server w/ 2hop relay
//This program will continually receive tcp packets
//from the client via 2hop relay scenario for the specified received time.
//It will then divide the total bits received from the
//client by the received time to get the actual
//data rate of the communication.

import java.io.*;
import java.net.*;


class TCP_Server
{
   private static int num_bytes = 256; //the packet size.
   private static int received_time = 30000; //the specified time the server will repeatedly
                                               //received packets from the client in milliseconds.
   public static int total_bits_received = 0; //the total bits received during the received time.
   public static int timer_flag = 0; //indicated if the timer thread if it has already started.
   
   //timer thread simulates the time duration in which the server will allow
   //receiving of packets from the client. Once the received_time specified
   //elapsed this thread will also show the actual data rate of communication.
   public static Thread timer = new Thread()
      {
         public void run()
         {
            //finishtime is the time when this thread will stop the server
            //from receiving packets from the client.
            long finishtime = System.currentTimeMillis() + received_time; 

            //wait until the current time is == to the finishtime
            while(System.currentTimeMillis() <= finishtime)
               ;
            int a = total_bits_received; //save the total bits received to a different
                                         //variable since the main function will still
                                         //accepts packet after the received time has
                                         //elapsed.
            int actual_data_rate = a / (received_time/1000); //calculate the actual data rate
                                                             //by dividing the total bits received
                                                             //by the received time in seconds.
            //show the actual data rate.
            System.out.println("Actual data rate: " + actual_data_rate + " bits/s");
            System.exit(0);
         }
      };


   //main function creates a socket and data stream for receiving tcp packets.
   public static void main(String args[]) throws Exception
   {
      ServerSocket serverSocket = new ServerSocket(4321); //create a serverSocket with appropriate
                                                          //port number for receiving.
      byte[] receivedData = new byte[num_bytes]; //prepare the data buffer where the
                                                //data received from the socket will be placed.

      while(true) //continue receiving packets from the client.
      {
         Socket connectionSocket = serverSocket.accept(); //accept the connection from client.
         if(timer_flag == 0) //when the first packet is received start the timer thread.
         {
            timer.start();
            timer_flag = 1;
         }
         //create input stream for reading the packets sent by the client.
         InputStream stream = connectionSocket.getInputStream();
         int length = stream.read(receivedData); //read the data and store how many
                                                 //bytes have been read from the stream.
         total_bits_received = total_bits_received+(length*8); //calculate the bits 
                                                               //received from the client.
      }
   }
} 