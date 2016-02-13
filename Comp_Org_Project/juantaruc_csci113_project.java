/*
Juan Taruc: CSCI 113 Project Fall 2015
This program is a simulation program for the 16-bit ALU 
that performs AND, OR, ADD and SUB operations. This program also
simulates unsigned number multiplication.

The program asks for 4 inputs: A, B, op and Binv
A = the 1st operand
B = the 2nd operand
op = the operation that will be executed {00-AND,01-OR,10-ADD}
Binv = inverts the bits in B for SUB operation

The program will display the results of operations in the form:
A      = 1st operand
B      = 2nd operand
op     = operation
Binv   = B invert signal
result = result of operation 
cout   = any carry out
of     = overflow in addition

Also, the program will ask the user if he/she wanted to proceed to
multiplication operation:
If yes: Multiplication is executed showing the state of registers in every
		iteration/cycle. Finally the program will display the results in the form:
		MD(A)   = the multiplicand
		MQ(B)   = the multiplier
		product = MD * MQ
If no: the program termninates.

To compile/run in terminal: javac filename -> java ALU
*/


import javax.swing.JOptionPane; //for JOptionPane dialog boxes


//main class
class ALU
{
	private static int max_bits = 16; //number of bits
	public static int sum; //reference value for result of ADD
	public static int cout; //reference value for carry out
	public static int of; //reference value for overflow
	public static String message = new String();
	public static void main(String args[])
	{
		String inputA; //first user input(A)
		String inputB; //second user input(B)
		String inputOp; //third user input(op)
		String inputBinv; //fourth user input(Binv)

		int a[] = new int[max_bits]; //the 1st operand
		int b[] = new int[max_bits]; //the 2nd operand
		int binv; //B invert signal

		int i;//loop counter
		//get a[]
		while(true) //while unsuccessful
		{
			inputA = JOptionPane.showInputDialog("Enter A:");//get 1st user input
			try 
			{
				for(i=0; i<max_bits; i++)//store the input to a[]
				{
					char at = inputA.charAt(i);
					a[i] = Character.getNumericValue(at);
					if(a[i] > 1 || a[i] < 0 || inputA.length() != max_bits)//if errors ask input again
					{
						Exception e = new Exception();
						throw e;
					}
				}
				break;
			}
			catch(Exception e)//error handler
			{
				if(inputA == null)//if user press cancel, terminate program
					System.exit(0);
				JOptionPane.showMessageDialog(null, "Invalid A!");
			}
		}

		//get b[]
		while(true)//while unsuccessful
		{
			inputB = JOptionPane.showInputDialog("Enter B:");//get 2nd user input
			try 
			{
				for(i=0; i<max_bits; i++)//store input to b[]
				{
					char at = inputB.charAt(i);
					b[i] = Character.getNumericValue(at);
					if(b[i] > 1 || b[i] < 0 || inputB.length() != max_bits)//if errors ask input again
					{
						Exception e = new Exception();
						throw e;
					}
				}
				break;
			}
			catch(Exception e)//error handler
			{
				if(inputB == null)//if user press cancel, terminate program
					System.exit(0);
				JOptionPane.showMessageDialog(null, "Invalid B!");
			}
		}
		
		//get op
		while(true)//while unsuccessful
		{
			inputOp = JOptionPane.showInputDialog("Enter operation:");//get 3rd user input
			if(inputOp.equals("00") || inputOp.equals("01") 
					|| inputOp.equals("10"))//op match acceptable operations
				break;
			else//error handler
			{
				if(inputOp == null)//if user press cancel, terminate program
					System.exit(0);
				JOptionPane.showMessageDialog(null, "Invalid op! Choices{00, 01, 10}");
			}
		}

		//get Binv
		while(true)
		{
			inputBinv = JOptionPane.showInputDialog("Enter Binv:");//get 4th user input
			if(inputBinv.equals("0") || inputBinv.equals("1"))//Binv match acceptable signal(s)
			{
				char at = inputBinv.charAt(0);
				binv = Character.getNumericValue(at);
				break;
			}
			else//error handler
			{
				if(inputBinv == null)//if user press cancel, terminate program
					System.exit(0);
				JOptionPane.showMessageDialog(null, "Invalid Binv! Choices{0,1}");
			}
		}
		
		//show the results of the operation
		JOptionPane.showMessageDialog(null, 
			"A:      " + inputA + "\n" + 
			"B:      " + inputB + "\n" +
			"op:     " + inputOp + "\n" +
			"Binv:   " + inputBinv + "\n" +
			"result: " + print(ALU_16bit(a,b,inputOp.toCharArray(),binv),0) + "\n" +
			"cout:   " + cout + "\n" +
			"of:     " + of);
		
		//ask the user if he/she wants to proceed to multiplication
		int mul = JOptionPane.showConfirmDialog(null,"Multiplication?");
		//if yes, then show the states of registers in Multiplier_16bit in every
		//cycle/iteration then show product
		if(mul == JOptionPane.YES_OPTION)
		{
			Multiplier_16bit(a,b);		}
	//done main()
	}


	//AND gate accepts two inputs and returns (a and b)
	public static int AND_gate(int a, int b)
	{
		if(a == 1 && b == 1)
			return 1;
		else
			return 0;
	}


	//OR gate accepts two inputs and returns (a or b)
	public static int OR_gate(int a, int b)
	{
		if(a == 1 || b == 1)
			return 1;
		else
			return 0;
	}


	//Full Adder adds a and b and also manages carry in and carry out
	public static void Full_Adder_1bit(int a, int b, int cin)
	{
		int xor_ab = 0;
		if((a == 1 && b == 0) || (a == 0 && b == 1))//a XOR b
			xor_ab = 1;

		//sum = (a XOR b) XOR cin
		if((xor_ab == 1 && cin == 0) || (xor_ab == 0 && cin == 1))
			sum = 1;
		else
			sum = 0;

		int cin_and_xor_ab = AND_gate(cin,xor_ab);//(cin and (a xor b))
		int a_and_b = AND_gate(a,b);//(a and b)
		cout = OR_gate(cin_and_xor_ab, a_and_b);//cout = (cin and (a xor b)) or (a and b)
	}


	//4x1 Multiplexor takes the results of AND, OR, ADD/SUB operations and selects
	//which one is the appropriate result for the op signal.
	//dummy as it's name implies just a dummy with value 0
	public static int Mux_4x1(int and, int or, int sum, int dummy, char sel[])
	{
		//4way decoder
		if(sel[0] == '0' && sel[1] == '0')//AND operation
			return and;
		else if(sel[0] == '0' && sel[1] == '1')//OR operation
			return or;
		else if(sel[0] == '1' && sel[1] == '0')//ADD/SUB operation
			return sum;
		else//return 0
			return dummy;
	}


	//1-bit ALU performs AND,OR,ADD/SUB operations and returns results
	public static int ALU_1bit(int a, int b, int cin, char op[])
	{
		int and = AND_gate(a,b);//AND
		int or = OR_gate(a,b);//OR
		Full_Adder_1bit(a,b,cin);//ADD/SUB
		return Mux_4x1(and,or,sum,0,op);//determine which result to return
	}


	//1-bit ALU with overflow checking. Same as the 1-bit ALU with the addition
	//of overflow checking.
	public static int ALU_1bit_OF(int a, int b, int cin, char op[])
	{
		int and = AND_gate(a,b);//AND
		int or = OR_gate(a,b);//OR
		Full_Adder_1bit(a,b,cin);//ADD/SUB
		//check overflow
		if(a == b)//if sign bit of a and b are same then possible overflow
		{
			if(a != sum)//if sign bit of sum is not the same as a and b, then overflow
				of = 1;
		}
		else//no overflow
			of = 0;
		return Mux_4x1(and,or,sum,0,op);//determine which result to return
	}


	//16-bit ALU calls 1-bit ALU n-1 times then calls 1-bit ALU with overflow checking
	//for the last/sign bit.
	public static int[] ALU_16bit(int a[], int b[], char op[], int binv)
	{
		int result[] = new int[max_bits];//initialize the result to return
		int newb;//for Binv
		int i;//loop counter
		int cin = 0;//initially, cin = 0
		//if Binv = 1 then need to convert b to 2's complement
		//to do that, the initial cin needs to be 1
		if(binv == 1)
			cin = 1;
		//reverse loop to start at the rightmost bit
		for(i=max_bits-1; i>0; i--)
		{
			//2x1 multiplexor
			if(binv == 1)//if Binv = 1 then need to invert b[i]
			{
				if(b[i] == 0)
					newb = 1;
				else
					newb = 0;
			}
			else//else if Binv = 0 keep original b[i]
				newb = b[i];
			result[i] = ALU_1bit(a[i],newb,cin,op);//call 1-bit ALU n-1 times
			cin = cout;//manage carry in and carry out
		}
		//2x1 multiplexor for the sign bit
		if(binv == 1)//if Binv = 1 then need to invert b[0]
		{
			if(b[0] == 0)
				newb = 1;
			else
				newb = 0;
		}
		else//else if Binv = 0 keep original b[i]
			newb = b[0];
		//call 1-bit ALU with overflow checking for last/sign bit
		result[0] = ALU_1bit_OF(a[0],newb,cin,op);

		return result;
	}


	//16-bit Multiplier multiplies two unsigned binary numbers.
	//This function also display the register states in each iteration/cycle.
	//Uses the 16-bit ALU to perform addition
	public static void Multiplier_16bit(int md[], int mq[])
	{
		int acc[] = new int[max_bits];//accumulator
		int product[] = new int[max_bits*2];//the product(needs to 
									//have size max_bits*2 for the result to fit)
		int counter = max_bits;//cycle counter
		int temp1 = 0;//for right shift
		int temp2 = 0;//for rigth shift
		int i,itr=1;//iteration
		
		//show the initial state of the registers md,mq,acc
		//the 6th parameter is the initial flag so the
		//print helper function knows to print the initial
		//state of the registers.
		//the 7th parameter is the step(s) in the multiplication algorithm
		show_registers(md,acc,mq,counter,itr,1,0);

		while(counter > 0)//repeat n times
		{
			//STEP 1:
			//if mq[0] == 1 then acc = acc + md
			//since the index of mq is in the order 0...max_bits-1,
			//then mq[0] is really mq[max_bits-1]
			if(mq[max_bits-1] == 1)
			{ 
				acc = ALU_16bit(acc,md,"10".toCharArray(),0);
			}
			
			//show the register states after the 1st step
			show_registers(md,acc,mq,counter,itr,0,1);

			//STEP 2:
			//right shift acc/mq >> 1
			//initially temp1 = 0, so in every initial rigth shift
			//0 is inserted in the leftmost bit.
			for(i=0; i<max_bits; i++)
			{
				temp2 = acc[i];//save the current bit
				acc[i] = temp1;//change acc[i] with the previous bit
				temp1 = temp2;//save current bit to temp1 which will 
							//become the previous bit for next iteration
			}
			temp2=0;//reset temp2
			for(i=0; i<max_bits; i++)
			{
				temp2 = mq[i];//save the current bit
				mq[i] = temp1;//change mq[i] with the previous bit
				temp1 = temp2;//save current bit to temp1 which will 
							//become the previous bit for next iteration
			}
			temp1=0;temp2=0;//reset temp1 and temp2 for next iteration/cycle

			//show the register states after the 2nd step
			show_registers(md,acc,mq,counter,itr,0,2);

			counter--;
			itr++;
		}
		//set the product
		for(i=0; i<max_bits; i++)//the first half of the product is the acc
			product[i] = acc[i];
		for(i=max_bits; i<max_bits*2; i++)//the second half of the product is the mq
			product[i] = mq[i-max_bits];
		message = message + "\n" + "Product:  " + print(product,1) + "\n";
		JOptionPane.showMessageDialog(null,message);
	}


	//This function is the print helper function
	//r[] is the array that will be printed
	//flag = 0 for AND,OR,ADD/SUB (result[max_bits])
	//flag = 1 for Multiplication (result[max_bits*2])
	public static String print(int r[], int flag)
	{
		String result = "";
		if(flag == 0)//for AND, OR, ADD/SUB
		{
			for(int i=0; i < max_bits; i++)
			{
				result = result + Integer.toString(r[i]);
			}
		}
		else//for Multiplication
		{
			for(int i=0; i<max_bits*2; i++)
			{
				result = result + Integer.toString(r[i]);
			}
		}
		return result;
	}


	//This function is the show_registers helper function.
	//It shows the states of the registers in each step in the Multiplication
	//operation.
	//parameters:
	//md = multiplicand
	//acc = accumulator
	//mq = multiplier
	//counter = cycle counter
	//itr = iteration counter
	//iflag = initial state flag
	//step = step in multiplication algorithm
	public static void show_registers(int md[], int acc[], int mq[], 
		int counter, int itr, int iflag, int step)
	{
		String s_itr;//for some readability. For 1st,2nd,3rd,4th.... strings
		int temp = itr;//for determining is st,nd,rd,or th
		String format_label = "%-40s%-40s%-40s%-45s%s\n";
		String format_content = "%-30s%-30s%-30s%-40s%s\n";
		if(iflag == 1)//show initial state
		{
			message = message + String.format(format_label,"Counter","MD","ACC","MQ","Step");
			message = message + String.format(format_content,print_i(counter),print(md,0),print(acc,0),print(mq,0),"--Initial");
			return;//done
		}
		if(itr > 20)//to determine if st,nd,rd,th
			temp = temp % 10;
		switch(temp)
		{
			case 1://_1st
				s_itr = "st";
			break;
			case 2://_2nd
				s_itr = "nd";
			break;
			case 3://_3rd
				s_itr = "rd";
			break;
			default://_xth
				s_itr = "th";
			break;
		}
		String state = "--After step" + step + 
			" of " + itr + s_itr + " iteration";
		state = 
		message = message + String.format(format_content,print_i(counter),print(md,0),print(acc,0),print(mq,0),state);
	}

	//Function print_i prints the counter in binary
	public static String print_i(int x)
	{
		char[] array = new char[5];//the array of binary bits
		String counter = "";
		for (int i=0; i < 5; i++)
		{
	    	if (x % 2 == 0)//0 bit
	          array[(5-1)-i] = '0';
	    	else if (x % 2 == 1)
	          array[(5-1)-i] = '1'; //1 bit
	    	x = x / 2;
		}
		for(int i=0; i < 5; i++)//convert each binary char to string and concat them
			counter = counter + array[i];
		return counter;//return the counter
	}
}
