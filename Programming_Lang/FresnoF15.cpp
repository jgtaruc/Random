/***************************************************************
Juan Taruc CSCI 117 Lab3 Assignment FresnoF15.cpp Interpreter.
This program reads a file (w/c in this case is "FresnoF15.txt") 
that is a simulationof a 'programming language'. This program is
the interpreter for FresnoF15 w/c reads each line of that file 
and executes each statement according to syntax in BNF given. 
If there are any syntax,semantic or lexical error(s) present in
the code, the interpreter will catch all error(s) it could 
before showing it in the screen.  
***************************************************************/

//Header files
#include <iostream>
#include <stdio.h>
#include <stdlib.h>
#include <fstream>
#include <string.h>
#include <sstream>
#include <vector>
#include <math.h>

using namespace std;


//Global variables
ifstream programFile;
string filename;

//The struct object that represents the variable 'Id' and it's value
struct sym_table {
    char id;
    int val;
};

//vector containing the sym_tabl(althoug it is not really a table)
vector<sym_table> var;
string aword;       //string used in reading each string in the file

bool compiling;         //if true then the interpreter is just looking for error(s) w/o doing
                        //the essential functions of input,output and assignment.
        
bool has_errors=false;    //if true then interpreter has find errors in the code. 
vector<string> error_msgs;//vector containing all the error messages of the code.
int total_errors = 0;     //the number of error(s) present in the code.


//Function Prototypes
void Prog(), //<Prog> ::= program <Declarations> begin <Statements> endl
    DeclS(), //<Declarations> ::= <Declaration> | <Declaration> <Declarations>
    Decl(),  //var <Id-list> ; | ε
    StmtS(), //<Statements> ::= <Statement> <Statements> | ε
    Stmt(string), //<Assign-St> | <Input-St> | <Output-St>
    Input_Stmt(), //<Input-St> ::= input <Id> ;
    Output_Stmt(), //<Output-St> ::= output <Id> ; | output <Exp> ;
    Assign_Stmt(string); //<Assign-St> ::= <Id> = <Exp> ;
bool Id(char); //<Id> ::= a|b|c| ... |z|A|B|C| ... |Z

bool test_sym_table(char);//function for checking if variable Id is known/declared
int get_sym_index(char); //function that returns the index which will be use to get 
                         //variable Id
void display_errors();  //function for displaying error(s)

//Lab2 provided by TA
int Exp(), Term(), Exp2(int), Term2(int), Fact(), Pwr();


//main driver function
int main()
{
    filename = "FresnoF15.txt";
    programFile.open(filename.c_str());
    compiling = true;
    Prog();
    if(has_errors)
    {
        display_errors();
        exit(0);
    }
    programFile.close();

    var.clear();
    compiling = false;

    programFile.open(filename.c_str());
    Prog();
    programFile.close();
    return 0;
}


//<Prog> ::= program <Declarations> begin <Statements> end
void Prog()
{
    programFile >> aword;
    //call Declarations to examine the variable(s) and make it known
    if(aword != "program")
    {
        has_errors = true;
        total_errors++;
        error_msgs.push_back(filename+": error: Expected 'program' at the start of program file.");
    } 
    DeclS();
    //call Statements to examine each statement and execute it
    StmtS();
}


//<Declarations> ::= <Declaration> | <Declaration> <Declarations>
void DeclS()
{
    
    programFile >> aword;
    if(aword == "begin")    //if 'begin' then return to Prog() and execute StmtS()
        return;
    else if(aword == "var") //if 'var' then call Decl() to assign the variables to global vector 
        Decl();
    else                    //else throw error
    {
        has_errors = true;
        total_errors++;
        error_msgs.push_back(filename+": error: Invalid variable declaration(s) or 'begin' not found.");
        programFile >> aword;
        return;
    }
    DeclS();      //DeclS() calls itself to accept multiple variable declarations
}


//<Declaration> ::= var <Id-list> ; | ε
void Decl()
{
    char aId;
    programFile >> aId;
    if(Id(aId))                       //is it a valid id?
    {
        if(test_sym_table(aId))       //is it already declared?
        {
            has_errors = true;
            total_errors++;
            error_msgs.push_back(filename+": error: Multiple/Invalid variable declarations.");
        }
        else                          //if it is not yet declared store it to variable vector
        {
            sym_table v;
            v.id = aId;
            v.val = 0;
            var.push_back(v);
            programFile >> aId;
            if(aId != ',')            //if aId is not ',' then it better be ';' else error 
            {
                if(aId != ';')
                {
                    has_errors = true;
                    total_errors++;
                    error_msgs.push_back(filename+": error: Expected ';' after variable declaration(s).");
                    return;
                }
            }
            else                     //else if aId is ';' then read next variable    
                Decl();
        }
    }
    else
    {
        has_errors = true;
        total_errors++;
        error_msgs.push_back(filename+": error: Invalid variable declaration(s).");
    }
}


//<Statements> ::= <Statement> <Statements> | ε
void StmtS()
{ 
    programFile >> aword;
    if(aword == "end")    //if next line is 'end' finish the program else execute each Statement
    {
        return;
    }
    else if(aword == "input" or aword == "output" or (aword.length() == 1 or strchr(aword.c_str(), '=')))
        Stmt(aword);
    else
    {
        has_errors = true;
        total_errors++;
        error_msgs.push_back(filename+": error: Invalid statement(s) or 'end' not found.");
        return;
    }            
    StmtS();
}


//<Statement> ::= <Assign-St> | <Input-St> | <Output-St>
void Stmt(string aword)
{
    if(aword == "input")
        Input_Stmt();
    else if(aword == "output")
        Output_Stmt();
    else 
        Assign_Stmt(aword);
}



//<Input-St> ::= input <Id> ;
void Input_Stmt()
{
    char aId,next;
    int temp,index;
    programFile >> aId;

    if(Id(aId) and test_sym_table(aId))     //is it a valid id?
    {
        if(not(compiling))                  //if aId is declared and not compiling get input from keyboard
        {
            cout<<aId<<" >> ";cin>>temp;
            index = get_sym_index(aId);
            var[index].val = temp;         
        }
        programFile >> next;                            
        if(next != ';')                      //after reading the input variable check if the next char after it is
        {                                    //a ';' if not then error
            has_errors = true;
            total_errors++;
            error_msgs.push_back(filename+": error: Expected ';' after statement 'input "+aId+"'.");
            programFile.putback(next);
        }                      
    }
    else                                     //else if aId is not valid or unknown 
    {
        has_errors = true;
        total_errors++;
        error_msgs.push_back(filename+": error: Unknown/Invalid variable in 'input' statement(s).");
    }
}


//<Output-St> ::= output <Id> ; | output <Exp> ;
void Output_Stmt()
{
    char aId,next;
    int index,temp;
    programFile >> aId;

    if(Id(aId) and test_sym_table(aId))//is it a valid id?
    {
        if(not(compiling))             //if aid is declared and not compiling then print value of aId
        {    
            index = get_sym_index(aId);
            cout<<var[index].val<<endl;
        }
        programFile >> next;
        if(next != ';')                //after printing out the value of aId check the next character if
        {                              //it is a ';' if not then error
            has_errors = true;
            total_errors++;
            error_msgs.push_back(filename+": error: Expected ';' after statement 'output "+aId+"'.");
            programFile.putback(next);
        }                                
    }

    else if(aId == '(' or isdigit(aId))//if aId is an expression then call expression and print it
    {
        programFile.putback(aId);
        temp = Exp();
        if(not(compiling))
        {
            cout<<temp<<endl;
        }
        programFile >> next;
        if(next != ';')                //after printing out the value of aId check the next character if
        {                              //it is a ';' if not then error
            has_errors = true;
            total_errors++;
            error_msgs.push_back(filename+": error: Expected ';' after statement 'output <Expression>.");
            programFile.putback(next);
            display_errors();
            exit(0);
        }                  
    }
    else                               //else if aId is not valid or unknown
    {
        has_errors = true;
        total_errors++;
        error_msgs.push_back(filename+": error: Unknown/Invalid variable or <Expression> in 'output' statement(s).");
    }
}


//<Assign-St> ::= <Id> = <Exp> ;
void Assign_Stmt(string word)
{
    char aId1,aId2, next;
    int temp,index;
    for(int i=0; i<word.size(); i++)   //restore the string word to programFile then just read char by char
        programFile.unget();
    programFile >> aId1;
    if(Id(aId1) and test_sym_table(aId1))//is it a valid Id?
    {
        programFile >> next;
        if(next == '=')                  //if next char is '=' then read next char to see if it an expression or Id
        {
            programFile >> aId2;
            if(aId2 == '(' or isdigit(aId2))//if expression then assign the result of expression to aId1
            {
                programFile.putback(aId2);
                temp = Exp();
                if(not(compiling))
                {
                    index = get_sym_index(aId1);
                    var[index].val = temp;
                }
                programFile >> next;        //if next char after assignment statement is not a ';' then error
                if(next != ';')
                {
                    has_errors = true;
                    total_errors++;
                    error_msgs.push_back(filename+": error: Expected ';' after '"+aId1+" <= <Expression> in 'assignment' statement.");
                    display_errors();
                    exit(0);
                }      
            }
            else if(Id(aId2) and test_sym_table(aId2))//else if aId2 is a valid Id then assign value of aId2 to aId1
            {
                if(not(compiling))
                {
                    temp = get_sym_index(aId2)  ;
                    index = get_sym_index(aId1);
                    var[index].val = var[temp].val;
                }
                programFile >> next;          //if next char after assignment statement is not a ';' then error
                if(next != ';')
                {
                    has_errors = true;
                    total_errors++;
                    error_msgs.push_back(filename+": error: Expected ';' after "+aId1+" <= "+aId2+" in 'assignment' statement.");
                    programFile.putback(next);
                }      
            }
            else      //else if aId is not an expression or known variable then error
            {
                has_errors = true;
                total_errors++;
                error_msgs.push_back(filename+": error: Unknown/Invalid variable or <Expression> in 'assignment' statement(s).");
            }                         
        }
        else          //else if next char is not '=' then it is unknown/invalid statement
        {
            has_errors = true;
            total_errors++;
            error_msgs.push_back(filename+": error: Unknown/Invalid statement(s).");
            programFile.putback(next);
        }

    }
    else              //else if aId is not a variable then unknown/invalid statement
    {
        has_errors = true;
        total_errors++;
        error_msgs.push_back(filename+": error: Unknown/Invalid statement(s).");
    }
}


//<Id> ::= a|b|c| ... |z|A|B|C| ... |Z
bool Id(char aId)
{
    //check aId if it is in the range a ~ Z else throw error
    if((aId >= 'A' and aId <= 'Z') or (aId >= 'a' and aId <= 'z'))
        return true;
    else
        return false;
}


//This function returns True if variable 'a' is already
//declared in the sym_table. False otherwise.
bool test_sym_table(char a)
{
    for(int i=0; i<var.size(); i++)
    {
    //if a is already declared variable then true else false
        if(var[i].id == a)     
            return true;
    }
    return false;
}


//This function takes char a (the variable), finds it, and returns the index
//which corresponds to the variable and its value in the sym_table represented by a vector.
int get_sym_index(char a)
{
    for(int i=0; i<var.size(); i++)
    {
        if(var[i].id == a)//if found then return index
            return i;   
    }                     //this function guarantees that it will return some index i
}                         //because I already check if char 'a' is in the sym_table before
                          //calling this function


//display all the errors  the interpreter caught.
void display_errors()
{
    cout<<endl<<endl;
    for(int i = 0; i < error_msgs.size(); i++)
        cout<<error_msgs[i]<<endl;
    cout<<"Total Error(s): "<<total_errors<<endl;
    cout<<endl<<endl;
}


//*********************************************************
//This is the solution of Lab2 that was provided by the TA
//with some modifications.
int Exp()
{
    return Exp2(Term());
}

int Term()
{
    return Term2(Pwr());
}

int Exp2(int inp)
{
    int result = inp;
    char a;
    if(programFile >> a)
    {
        if(a == '+')
        {
            result = Exp2(result + Term());
        }
        else if(a == '-')
        {
            result = Exp2(result - Term());
        }
        else if(a == ')' || a == ';')
            programFile.putback(a);
    }
    
    return result;
}

int Term2(int inp)
{
    int result = inp;
    char a;
    if(programFile >> a)
    {
        if (a == '*')
            result = Term2(result * Pwr());
        else if (a == '/')
            result = Term2(result / Pwr());
        else if (a == '+' || a == '-' || a == ')' || a == ';')
            programFile.putback(a);
    }
    
    return result;
}

int Pwr() {
    int p = Fact();
    char a;
    
    if (programFile >> a) {
        if (a == '^') {
            p = pow(p, Pwr());
        }
        if (a == '+' || a == '-' || a == '*' || a == '/' || a == ')' || a == ';')
            programFile.putback(a);
    }
    
    return p;
}

int Fact()
{
    char a;
    programFile.get(a);
    while(a == ' ')         //ignore whitespaces
        programFile.get(a);
    if (a == '(')
    {
        int exp = Exp();
        programFile >> a;
        if (a == ')')
        {
            return exp;
        }
    }
    else if(not(isdigit(a)))//if a is not a number then error
    {
        has_errors = true;
        total_errors++;
        error_msgs.push_back(filename+": error: Not valid <Expression> on 'output' or 'assignment' statement.");
    }
    return (a - '0');
}