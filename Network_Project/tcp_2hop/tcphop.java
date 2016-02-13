//Juan Taruc CSCI 156 Project - TCP hop w/ 2 hop relay
//This program will continually receive tcp packets
//from the client and forward it to the server.

import java.io.*;
import java.net.*;


class TCP_Hop
{
   private static int num_bytes = 256; //the packet size.

   public static void main(String argv[]) throws Exception
      {
         ServerSocket serverSocket = new ServerSocket(1234); //create a serverSocket for 
                                                             //receiving.
         byte[] receivedData = new byte[num_bytes]; //prepare the data buffer where the
                                                //data received from the socket will be placed.
         while(true) //continue receiving packets from client and forward it to the client.
         {
            Socket connectionSocket = serverSocket.accept(); //accept connection with client.
            InputStream stream = connectionSocket.getInputStream(); //create an input stream
                                                                    //for reading.
            stream.read(receivedData); //read the data in the stream.

            Socket toServerSocket = new Socket("localhost", 4321); //create a socket
                                                                        //for sending to server.
            //create a stream for writing the data to the server socket.
            DataOutputStream outToServer = 
               new DataOutputStream(toServerSocket.getOutputStream());
            outToServer.write(receivedData); //write the data to server.
            toServerSocket.close();
         }
      }
}
