//Juan Taruc: CSCI 164 Project Spring 2016
//Tile Sliding Puzzle Solver in Go


package main

import(
	"math"
	"fmt"
	"container/heap"
	"math/rand"
	"time"
)

//8 puzzle
const size = 3
type state [size][size]int //represents the state of the sliding tile puzzle.
var goalState = state{{0,1,2},{3,4,5},{6,7,8}} //goal configuration of sliding tile puzzle.


type node struct{ //used for stack and queue data
	state 	state //the state of the puzzle
	action  string //actions taken to reach the state
	cost    int
	priority int //only used for priority queue
	index    int
	blankPos [2]int
}
 
type queue []*node
func (q *queue) push(x interface{}){
	n := x.(*node)
	*q = append(*q, n)
}
func (q *queue) pop() interface{}{
	if q.isEmpty() {
		return nil
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return item
}
func (q *queue) isEmpty() bool {
	return len(*q) == 0
}



/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// A PriorityQueue implements heap.Interface and holds node.
type PriorityQueue []node
func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(node)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *node, value state, priority int) {
	item.state = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) isEmpty() bool{
	return (*pq).Len() == 0
}
//////////////////////////////////////////////////////////////////////////////////////////////////


//Manhattan distance heuristic
func manhattanDistance(current state, goal state, goalMap map[int][2]int) int {
	totalDist := 0
	for i := 0; i < size; i++{
		for j := 0; j < size; j++{
			if current[i][j] != goal[i][j]{
				if current[i][j] == 0 { continue }
				goalXY := goalMap[current[i][j]]
				totalDist += int(math.Abs(float64(i - goalXY[0] )) + math.Abs(float64(j - goalXY[1])))
			}
		}
	}
	return totalDist
}


//Misplaced Tiles heuristic
func misplacedTiles(current state, goal state) int {
	total := 0
	for i := 0; i < size; i++{
		for j := 0; j < size; j++{
			if current[i][j] != goal[i][j]{
				total++
			}
		}
	}
	return total
}


//returns a list of successor nodes reachable through the current node
func getSuccessors(n node) []node {
	var successors []node
	for _, a := range getLegalMoves(n){
		successors = append(successors, result(a, n))
	}
	return successors
}


//updates the currentState with the move taken producing
//a newNode
func result(move string, current node) node {
	newNode := current
	row := current.blankPos[0]
	col := current.blankPos[1]
	var newrow int
	var newcol int

	//compute the new index of element 0(blank tile)
	if move == "up"{
		newrow = row - 1
		newcol = col
	}else if move == "down"{
		newrow = row + 1
		newcol = col
	}else if move == "right"{
		newrow = row
		newcol = col + 1
	}else if move == "left"{
		newrow = row
		newcol = col - 1
	}

	temp := newNode.state[newrow][newcol]
	newNode.state[newrow][newcol] = newNode.state[row][col]
	newNode.state[row][col] = temp
	newNode.blankPos[0] = newrow
	newNode.blankPos[1] = newcol
	newNode.action = current.action + "|" + move
	newNode.cost++

	return newNode
}


//returns a list of possible moves based on the position of element 0(blank tile)
func getLegalMoves(current node) []string{
	var moves = []string{}
	row := current.blankPos[0]
	col := current.blankPos[1]
	bound := len(current.state) - 1
	
	if row != 0{
		moves = append(moves, "up")
	}
	if row != bound{
		moves = append(moves, "down")
	}
	if col != 0{
		moves = append(moves, "left")
	}
	if col != bound{
		moves = append(moves, "right")
	}

	return moves
}


//returns the x,y position of element 0(blank tile)
func getBlankPos(this state) (int, int){
	var row int
	var col int
	loop:	
	for i := 0; i < size; i++{
		for j := 0; j < size; j++{
			if this[i][j] == 0{
				row = i
				col = j
				break loop
			}
		}
	}
	return row, col
}


//prints the state of tile sliding puzzle
func print(puzzle state){
	for i := 0; i < size; i++{
		fmt.Println(puzzle[i])
	}
	fmt.Println()
}


//checks the solution found by the search algorithms if
//the solution found will indeed lead to goalState
func checkMoves(moves []string, startNode node) {
	for _, a := range moves {
		startNode = result(a, startNode)
	}

	if startNode.state != goalState{
		panic("*** Wrong Solution Step ***")
	}
}


//generates a random 8 puzzle
//it also checks if the random puzzle generated is solvable or unsolvable
//base on the number of inversions.
//if # of inversion is even -> solvable
//else generate a new random puzzle until we generate a solvable puzzle
func randomPuzzle() state {
	var s state
	t := time.Now()
 	rand.Seed(int64(t.Nanosecond()))
 	temp := []int{0,1,2,3,4,5,6,7,8}
 	for i := len(temp) - 1; i > 0; i-- {
 		j := rand.Intn(i)
 		temp[i], temp[j] = temp[j], temp[i]
 	}
 	inversion := 0
 	for i := 0; i < len(temp); i++ {
 		if temp[i] == 0 { continue }
 		for j := i+1; j < len(temp); j++{
 			if temp[j] == 0 { continue }
 			if temp[i] > temp[j]{ inversion++ }
 		}
 	}

 	k := 0
 	for i := 0; i < size; i++{
 		for j := 0; j < size; j++{
 			s[i][j] = temp[k]
 			k++
 		}
 	}

 	if inversion % 2 != 0{ s = randomPuzzle() }
 	return s
}