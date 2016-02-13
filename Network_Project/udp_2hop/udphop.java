//Juan Taruc CSCI 156 Project - UPD hop w/ 2 hop relay
//This program will continually receive udp packets
//from the client and forward it to the server.


import java.io.*;
import java.net.*;

class UDP_Hop
{
   private static int num_bytes = 256; //the packet size.

   //main function continually receives a packet from the client and then
   //forwards it to the server.
   public static void main(String args[]) throws Exception
   {
      DatagramSocket serverSocket = new DatagramSocket(1234); //socket for receiving
      DatagramSocket clientSocket = new DatagramSocket(); //socket for sending
      
      byte[] receiveData = new byte[num_bytes]; //prepare the data buffer where the
                                                //data received from the socket will be placed.
      //packet for receiving the packet from the server.
      DatagramPacket receivePacket = 
            new DatagramPacket(receiveData, receiveData.length);
      //infinite loop for receiving and sending a packet.
      while(true)
      {
         serverSocket.receive(receivePacket); //receive a packet from the client
        
         InetAddress IPAddress = InetAddress.getByName("localhost"); //set the address of the server.
         //create a packet for sending with the data receive from the client together with the
         //appropriate address and port number.
         DatagramPacket sendPacket = 
            new DatagramPacket(receivePacket.getData(), receivePacket.getLength(), IPAddress, 4321);
         clientSocket.send(sendPacket); //send the packet to the server.
      }
   }
}
