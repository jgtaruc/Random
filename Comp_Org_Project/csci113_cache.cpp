/*
Juan Taruc CSCI 113 Cache simulation
This program simulates caching using arrays.
It takes machine code as input located in input.txt.
It then parses each line of input then determines which
operation to run, either lw or sw. It update the memory or 
the cache depending on the circumstances.
It then prints out the contents of the cache and memory 
after all the operations are done. 
*/

#include <iostream>
#include <string>
#include <iomanip>
#include <fstream>
#include <stdlib.h>
#include <math.h>

using namespace std;

//name of the file where the machine codes are saved
string filename = "input.txt"; 

struct block	//[[valid][tag][data]]
{
	int valid;
	int tag;
	int data;
	int lrubit; //use for determining if the block is least recently used.
};

//Components
int register_file[8]; //$s0 ~ $s7
block cache[128][2]; //the cache
int memory[128]; //0th ~ 127th word

//function prototypes
void initialize(); //initialize the components(register file, cache, memory) 
void read_mcode(ifstream&, string&, string&, string&, string&); //read machine code from input file
int bin_to_dec(string);
string dec_to_bin(int,int);
void read(int,int,int,int);//loads a memory block to the register
void write(int,int,int,int);//writes the contents of the register to memory

int main()
{
	ifstream MIPS_file; //to operate on input file
	string op, rt, rs, offset; //strings for reading from MIPS_file

	int offsetvalue, reg_number, set, tag, opvalue;
	int i,j;

	initialize(); //initialize the register file, cache and memory
	MIPS_file.open(filename.c_str()); //open input file

	if(!MIPS_file.is_open()) //if file failed to open terminate the program
	{
		cout<<"Error: Failed to open "<<filename<<endl;
		exit(0); //exit
	}

	while(!MIPS_file.eof())
	{
		read_mcode(MIPS_file, op,rt,rs,offset);//read from file

		offsetvalue = (bin_to_dec(offset) + bin_to_dec(rs)) / 4; //convert byte address to word address
		reg_number = bin_to_dec(rt); //get the register number

		//determine cache set number and tag field
		set = offsetvalue % 8;
		tag = offsetvalue / 8;

		opvalue = bin_to_dec(op);
		if(opvalue == 35)
			read(set, tag, reg_number, offsetvalue);
		else if(opvalue  == 43)
			write(set, tag, reg_number, offsetvalue);
		else
		{
			cout<<"Error: Invalid operation."<<endl;
			exit(0);
		}
	}
	cout<<endl;
	
	//display the contents of the register file, cache, main memory
	cout<<"                 Register file"<<endl;
	cout<<"     ------------------------------------"<<endl;
	for(i = 0; i < 8; i++)
	{
		cout<<"$s"<<i<<"  | "<<setw(32)<<dec_to_bin(register_file[i],32)<<" |"<<endl;
	cout<<"     ------------------------------------"<<endl;
	}
	cout<<endl;

	cout<<"                                                                Cache"<<endl;
	cout<<"      ---------------------------------------------";
	cout<<"--------------------------------------------------------------------------"<<endl;
	cout<<"      |                         block0                 ";
	cout<<"          |                     block1                               |"<<endl;
	cout<<"      ---------------------------------------------";
	cout<<"--------------------------------------------------------------------------"<<endl;
	cout<<"      | valid | history|  tag |               data               |";
	cout<<" valid | history|  tag |               data               |"<<endl;
	cout<<"      ---------------------------------------------";
	cout<<"--------------------------------------------------------------------------"<<endl;
	for(i = 0; i < 8; i++)
	{
		cout<<"set"<<i<<"  | ";
		for(j = 0; j < 2; j++)
		{
			cout<<setw(5)<<cache[i][j].valid<<" | ";
			cout<<setw(6)<<cache[i][j].lrubit<<" | ";
			cout<<setw(4)<<dec_to_bin(cache[i][j].tag, 4)<<" | ";
			cout<<setw(32)<<dec_to_bin(cache[i][j].data,32)<<" | ";
		}
		cout<<endl;
		cout<<"      ---------------------------------------------";
		cout<<"--------------------------------------------------------------------------"<<endl;
	}


	cout<<"                  Memory"<<endl;
	cout<<"          ------------------------------------"<<endl;
	for(i = 0; i < 128; i++)
	{
		cout<<setw(5)<<"Mem["<<i<<"] ";
		cout<<"  | "<<setw(32)<<dec_to_bin(memory[i],32)<<" |"<<endl;
		if(i < 9)
			cout<<"          ------------------------------------"<<endl;
		else if(i < 99)
			cout<<"           ------------------------------------"<<endl;
		else
			cout<<"            ------------------------------------"<<endl;
	}
	cout<<endl;

	return 0;
}


//Initialize the components(register,cache,memory)
// - registers($s0 ~ $s7) in the register file are initially 0.
// - cache entries are empty and their valid bits are 0.
// - memory is initialize: mem[address] = address + 5.
void initialize()
{
	int i,j;
	for(i = 0; i < 8; i++) //Initialize register file
		register_file[i] = 0;

	for(i = 0; i < 8; i++) //Initialize cache
	{
		for(j = 0; j < 2; j++)
			{
				cache[i][j].valid = 0;
				cache[i][j].tag = 0;
				cache[i][j].data = 0;
				cache[i][j].lrubit = 1; //initially the least recently used bit for all blocks are 1
			}
	}

	for(i = 0; i < 128; i++) //Initialize memory
		memory[i] = i + 5;
}


//read each line from input file
void read_mcode(ifstream &MIPS_file, string &op, string &rt, string &rs, string &offset)
{
	string mcode; 
	MIPS_file >> mcode; //read machine code
	op = mcode.substr(0,6); //get opcode
	rs = mcode.substr(6,5); //get rs
	rt = mcode.substr(11,5); //get rt
	offset = mcode.substr(16); //get offset

	//Bit error checking.
	if(op.length() != 6)
	{
		cout<<"Error: Invalid opcode '"<<op<<"', must be 6 bits."<<endl;
		exit(0);
	}
	if(rt.length() != 5)
	{
		cout<<"Error: Invalid rt '"<<op<<"', must be 5 bits."<<endl;
		exit(0);
	}
	if(rs.length() != 5)
	{
		cout<<"Error: Invalid rs '"<<op<<"', must be 5 bits."<<endl;
		exit(0);
	}
	if(offset.length() != 16)
	{
		cout<<"Error: Invalid rt '"<<op<<"', must be 16 bits."<<endl;
		exit(0);
	}
}


//converts a binary number to decimal number
int bin_to_dec(string s)
{
	int i,value = 0;
	string s2 = "";
	for(i = 0; i < s.length(); i++)
		s2 = s[i] + s2;
	for(i = 0; i < s.length(); i++)
	{
		if(s2[i] == '1')
			value = value + pow(2,i);
	}
	
	return value;
}


//converts a decimal number to binary with given num_bits
string dec_to_bin(int d, int num_bits)
{
	string binary = "";
	int i;
	for(i = 0; i < num_bits; i++)
	{
		if(d % 2 == 0)
			binary = "0" + binary;
		else if(d % 2 == 1)
			binary = "1" + binary;
		d = d / 2;
	}
	return binary;
}

//if valid block but miss, save cache block data to mem then update cache step2
//???how to implement writeback policy
void read(int set, int tag, int reg_number, int offsetvalue)
{
	if(cache[set][0].tag == tag && cache[set][0].valid == 1) //hit on block 0
	{
		cout<<"read_hit:   Mem["<<offsetvalue<<"] on cache["<<set<<"][0]"<<endl;
		register_file[reg_number-16] = cache[set][0].data;
		cache[set][0].lrubit = 0; //recently used block
		cache[set][1].lrubit = 1; //least recently used block
	}
	else if(cache[set][1].tag == tag && cache[set][1].valid == 1) //hit on block1
	{
		cout<<"read_hit:   Mem["<<offsetvalue<<"] on cache["<<set<<"][1]"<<endl;
		register_file[reg_number-16] = cache[set][1].data;
		cache[set][0].lrubit = 1; //least recently used block
		cache[set][1].lrubit = 0; //recently used block
	}
	else //miss
	{
		if(cache[set][0].valid == 0) //initial block0
		{
			cout<<"read_miss:  Mem["<<offsetvalue<<"]"<<endl;
			cache[set][0].tag = tag;
			cache[set][0].data = memory[offsetvalue];
			cache[set][0].valid = 1;
			register_file[reg_number-16] = cache[set][0].data;
		}
		else if(cache[set][1].valid == 0) //initial block1
		{
			cout<<"read_miss:  Mem["<<offsetvalue<<"]"<<endl;
			cache[set][1].tag = tag;
			cache[set][1].data = memory[offsetvalue];
			cache[set][1].valid = 1;
			register_file[reg_number-16] = cache[set][1].data;
		}
		else //valid blocks but does not match tag field
		{
			//form the memory address by combining tag(left 4 bits) and set(right 3 bits)
			string bin_index = dec_to_bin(tag,4) + dec_to_bin(set,3); //get the address in binary
			int index = bin_to_dec(bin_index)/4; //convert the binary address to decimal and make it word address
			cout<<"read_miss:  Mem["<<offsetvalue<<"]"<<endl;
			if(cache[set][0].lrubit == 1) //if block0 is the least recently used block
			{
				memory[index] = cache[set][0].data; //cache write-back
				cache[set][0].data = memory[offsetvalue]; //victim block <- memory
				cache[set][0].valid = 1; //update valid
				cache[set][0].tag = tag; //update tag
				register_file[reg_number-16] = cache[set][0].data; //target register <- victim block
			}
			else //if block1 is the least recently used block
			{
				memory[index] = cache[set][1].data; //cache write-back
				cache[set][1].data = memory[offsetvalue]; //victim block <- memory
				cache[set][1].valid = 1; //update valid
				cache[set][1].tag = tag; //update tag
				register_file[reg_number-16] = cache[set][1].data; //target register <- victim block
			}
		}
	}
}


void write(int set, int tag, int reg_number, int offsetvalue)
{
	//write-hit
	if(cache[set][0].tag == tag && cache[set][0].valid == 1) //hit on block 0
	{
		cout<<"write_hit:  Mem["<<offsetvalue<<"] on cache["<<set<<"][0]"<<endl;
		cache[set][0].data = register_file[reg_number-16];
		cache[set][0].lrubit = 0; //recently used block
		cache[set][1].lrubit = 1; //least recently used block
	}
	else if(cache[set][1].tag == tag && cache[set][1].valid == 1) //hit on block1
	{
		cout<<"write_hit:  Mem["<<offsetvalue<<"] on cache["<<set<<"][1]"<<endl;
		cache[set][1].data = register_file[reg_number-16];
		cache[set][0].lrubit = 1; //least recently used block
		cache[set][1].lrubit = 0; //recently used block
	}
	//write-miss, used no write allocate
	else //miss
	{
		cout<<"write_miss: Mem["<<offsetvalue<<"]"<<endl;
		//form the memory address by combining tag(left 4 bits) and set(right 3 bits)
		string bin_index = dec_to_bin(tag,4) + dec_to_bin(set,3); //get the address in binary
		int index = bin_to_dec(bin_index)/4; //convert the binary address to decimal and make it word address
		memory[index] = register_file[reg_number-16];		
	}
}