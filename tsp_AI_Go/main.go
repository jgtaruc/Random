//Juan Taruc: CSCI 164 Project Spring 2016
//Tile Sliding Puzzle Solver in Go

//driver program

package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"strconv"
)

var problem = randomPuzzle()

func main(){
	var input int
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("======== Enter Command ========")
	fmt.Println("0 - breadthFirstSearch")
	fmt.Println("1 - greedyBestFirstSearch w/ manhattan")
	fmt.Println("2 - greedyBestFirstSearch w/ misplaced")
	fmt.Println("3 - aStarSearch w/ manhattan")
	fmt.Println("4 - aStarSearch w/ misplaced")
	fmt.Println("5 - IDaStarSearch w/ manhattan")
	fmt.Println("6 - Run all")
	fmt.Println("7 - generate own puzzle")
	fmt.Println("8 - generate randomPuzzle")
	fmt.Println("Ctrl-c - Exit")
	fmt.Println("===============================")

	fmt.Println("\nRandom 8 Puzzle")
	print(problem)

	fmt.Println("Goal State")
	print(goalState)

	fmt.Print("> ")
	fmt.Scan(&input)

   	switch input {
   		case 0: 
   			bfs()
   			fmt.Print("Press [Enter] to continue or Ctrl-c to exit...")
   			reader.ReadString('\n')
   			main()
   		case 1: 
   			greedy1()
   			fmt.Print("Press [Enter] to continue or Ctrl-c to exit...")
   			reader.ReadString('\n')
   			main()
   		case 2: 
   			greedy2()
   			fmt.Print("Press [Enter] to continue or Ctrl-c to exit...")
   			reader.ReadString('\n')
   			main()
   		case 3: 
   			aStar1()
   			fmt.Print("Press [Enter] to continue or Ctrl-c to exit...")
   			reader.ReadString('\n')
   			main()
   		case 4: 
   			aStar2()
   			fmt.Print("Press [Enter] to continue or Ctrl-c to exit...")
   			reader.ReadString('\n')
   			main()
   		case 5: 
   			IDaStar()
   			fmt.Print("Press [Enter] to continue or Ctrl-c to exit...")
   			reader.ReadString('\n')
   			main()
   		case 6: 
   			bfs()
   			greedy1()
   			greedy2()
   			aStar1()
   			aStar2()
   			IDaStar()
   			fmt.Print("Press [Enter] to continue or Ctrl-c to exit...")
   			reader.ReadString('\n')
   			main()
   		case 7:
   			fmt.Println("Enter puzzle configuration in the form: 0 1 2 3 4 5 6 7 8")
   			fmt.Print("> ")
   			rawState, _ := reader.ReadString('\n')
   			rawState = strings.TrimSuffix(rawState, "\n")

   			stringState := strings.Split(rawState, " ")
   			k := 0
   			for i := 0; i < size; i++{
   				for j := 0; j < size; j++{
   					num, err := strconv.Atoi(stringState[k])
   					if err != nil {
   						panic(err)
   					}
   					if num < 0 || num > 8{
   						panic("*** Invalid configuration ***")
   					}
   					problem[i][j] = num
   					k++
   				}
   			}
   			print(problem)
   			fmt.Print("Press [Enter] to continue or Ctrl-c to exit...")
   			reader.ReadString('\n')
   			main()
   		case 8: 
   			problem = randomPuzzle()
   			main()
   		default: main()
   	}
}

func bfs() {
	var temp [2]int
	temp[0], temp[1] = getBlankPos(problem)
	startNode := node{
		state: problem,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}

	fmt.Println("===== breadthFirstSearch =====")
	moves, cost, expanded, t := breadthFirstSearch(problem)
	fmt.Println("Solution   Step:", strings.Split(moves, "|"))
	fmt.Println("Number of Moves:", cost)
	fmt.Println("Expanded  Nodes:", expanded)
	fmt.Println("Solution   Time:", t)
	checkMoves(strings.Split(moves, "|"), startNode)
	fmt.Println("==============================")
	fmt.Println()
}

func greedy1() {
	var temp [2]int
	temp[0], temp[1] = getBlankPos(problem)
	startNode := node{
		state: problem,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}
	fmt.Println("===== greedyBestFirstSearch w/ manhattan heuristic =====")
	moves, cost, expanded, t := greedyBestFirstSearchManhattan(problem)
	fmt.Println("Solution   Step:", strings.Split(moves, "|"))
	fmt.Println("Number of Moves:", cost)
	fmt.Println("Expanded  Nodes:", expanded)
	fmt.Println("Solution   Time:", t)
	checkMoves(strings.Split(moves, "|"), startNode)
	fmt.Println("========================================================")
	fmt.Println()
}

func greedy2() {
	var temp [2]int
	temp[0], temp[1] = getBlankPos(problem)
	startNode := node{
		state: problem,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}
	fmt.Println("===== greedyBestFirstSearch w/ misplaced heuristic =====")
	moves, cost, expanded, t := greedyBestFirstSearchMisplaced(problem)
	fmt.Println("Solution   Step:", strings.Split(moves, "|"))
	fmt.Println("Number of Moves:", cost)
	fmt.Println("Expanded  Nodes:", expanded)
	fmt.Println("Solution   Time:", t)
	checkMoves(strings.Split(moves, "|"), startNode)
	fmt.Println("========================================================")
	fmt.Println()
}

func aStar1() {
	var temp [2]int
	temp[0], temp[1] = getBlankPos(problem)
	startNode := node{
		state: problem,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}
	fmt.Println("===== aStarSearch w/ manhattan heuristic =====")
	moves, cost, expanded, t := aStarSearchManhattan(problem)
	fmt.Println("Solution   Step:", strings.Split(moves, "|"))
	fmt.Println("Number of Moves:", cost)
	fmt.Println("Expanded  Nodes:", expanded)
	fmt.Println("Solution   Time:", t)
	checkMoves(strings.Split(moves, "|"), startNode)
	fmt.Println("==============================================")
	fmt.Println()
}

func aStar2() {
	var temp [2]int
	temp[0], temp[1] = getBlankPos(problem)
	startNode := node{
		state: problem,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}
	fmt.Println("===== aStarSearch w/ misplaced heuristic =====")
	moves, cost, expanded, t := aStarSearchMisplaced(problem)
	fmt.Println("Solution   Step:", strings.Split(moves, "|"))
	fmt.Println("Number of Moves:", cost)
	fmt.Println("Expanded  Nodes:", expanded)
	fmt.Println("Solution   Time:", t)
	checkMoves(strings.Split(moves, "|"), startNode)
	fmt.Println("==============================================")
	fmt.Println()
}

func IDaStar() {
	var temp [2]int
	temp[0], temp[1] = getBlankPos(problem)
	startNode := node{
		state: problem,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}
	fmt.Println("===== IDaStarSearch w/ manhattan heuristic =====")
	moves, cost, expanded, t := IDaStarSearchManhattan(problem)
	fmt.Println("Solution   Step:", strings.Split(moves, "|"))
	fmt.Println("Number of Moves:", cost)
	fmt.Println("Expanded  Nodes:", expanded)
	fmt.Println("Solution   Time:", t)
	checkMoves(strings.Split(moves, "|"), startNode)
	fmt.Println("================================================")
	fmt.Println()
}