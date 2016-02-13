//Juan Taruc CSCI 156 Project - UPD serverw/ 2 hop relay
//This program will continually receive udp packets
//from the client via 2hop relay scenario for the specified received time.
//It will then divide the total bits received from the
//client by the received time to get the actual
//data rate of the communication.


import java.io.*;
import java.net.*;


class UDP_Server
{
   private static int num_bytes = 256; //the packet size.
   private static int received_time = 30000; //the specified time the server will repeatedly
                                               //received packets from the client in milliseconds.
   public static int timer_flag = 0; //indicates if the timer thread has already started
   public static int total_bits_received = 0; //the total bits received during the received time.
   
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


   //main function creates a socket and packet for receiving udp packets.
   public static void main(String args[]) throws Exception
   {
      DatagramSocket serverSocket = new DatagramSocket(4321); //create a packet with the
                                                              //port number specified.
      byte[] receiveData = new byte[num_bytes]; //prepare the data buffer where the
                                                //data received from the socket will be placed.
      //create a packet for receiving.
      DatagramPacket receivePacket = new DatagramPacket(receiveData, receiveData.length);
      //continually received a packet from the client.
      while(true)
      {
         serverSocket.receive(receivePacket); //receive a packet
         //when the first packet is received start the time thread.
         if(timer_flag == 0) 
         {
            timer.start();
            timer_flag = 1; //timer_flag prevents the timer flag from starting again
                            //every time a packet is received.
         }
         //calculate the bits received from the client.
         total_bits_received += (receivePacket.getLength()*8);
      }
   }
} 