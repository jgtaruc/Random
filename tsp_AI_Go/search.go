//Juan Taruc: CSCI 164 Project Spring 2016
//Tile Sliding Puzzle Solver in Go

/*
State-Space Search Algorithms used to solve the 8 puzzle problems
	Breadth-First Search - optimal but has exponential space complexity
	Greedy Best-First Search - not optimal but efficient
	A* Search - optimal but space complexity is still prohibitive
	IDA* Search - optimal and has linear space complexity but takes more time to find a solution
*/

package main

import(
	"strings"
	"time"
	"container/heap"
)


func breadthFirstSearch(initial state) (string, int, int, time.Duration){
	explored := make(map[state]bool)
	frontier := make(queue, 0)
	var temp [2]int
	temp[0], temp[1] = getBlankPos(initial)

	//push the initial node to the queue
	n := &node{
		state: initial,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}

	start := time.Now()
	frontier.push(n)
	for{
		if frontier.isEmpty() { return "", -1, -1, time.Since(start) }//fail
		x := frontier.pop().(*node)
		

		if x.state == goalState{
			//return action without leading "|" for easy splitting 
			return strings.TrimPrefix(x.action, "|"), x.cost, len(explored), time.Since(start)
		}
		
		//expand the node generating successor nodes and push successor nodes
		//to queue. 
		successors := getSuccessors(*x)
		for i := 0; i < len(successors); i++{
			//if successor is already explored get another successor
			if _, ok := explored[successors[i].state]; ok {
				continue
			}
			explored[x.state] = true
			frontier.push(&successors[i])
		}
	}
}


func greedyBestFirstSearchManhattan(initial state) (string, int, int, time.Duration){
	goalMap := make(map[int][2]int)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			goalMap[goalState[i][j]] = [2]int{i, j}
		}
	}	
	explored := make(map[state]bool)
	var frontier PriorityQueue
	heap.Init(&frontier)

	var temp [2]int
	temp[0], temp[1] = getBlankPos(initial)
	//push the initial node to the queue
	n := &node{
		state: initial,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}

	start := time.Now()
	heap.Push(&frontier, *n)
	for{
		if frontier.isEmpty() { return "", -1, -1, time.Since(start) }//fail
		x := heap.Pop(&frontier).(node)

		if x.state == goalState{
			//return action without leading "|" for easy splitting 
			return strings.TrimPrefix(x.action, "|"), x.cost, len(explored), time.Since(start) 
		}
			
		//expand the node generating successor nodes and push successor nodes
		//to queue. 
		successors := getSuccessors(x)
		for i := 0; i < len(successors); i++{
			if _, ok := explored[successors[i].state]; ok {
				continue
			}
			explored[x.state] = true
			successors[i].priority = manhattanDistance(successors[i].state, goalState, goalMap)
			heap.Push(&frontier, successors[i])
		}
	}
}


func greedyBestFirstSearchMisplaced(initial state) (string, int, int, time.Duration){
	goalMap := make(map[int][2]int)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			goalMap[goalState[i][j]] = [2]int{i, j}
		}
	}	
	explored := make(map[state]bool)
	var frontier PriorityQueue
	heap.Init(&frontier)

	var temp [2]int
	temp[0], temp[1] = getBlankPos(initial)
	//push the initial node to the queue
	n := &node{
		state: initial,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}

	start := time.Now()
	heap.Push(&frontier, *n)
	for{
		if frontier.isEmpty() { return "", -1, -1, time.Since(start) }//fail
		x := heap.Pop(&frontier).(node)

		if x.state == goalState{
			//return action without leading "|" for easy splitting 
			return strings.TrimPrefix(x.action, "|"), x.cost, len(explored), time.Since(start) 
		}
		
		//expand the node generating successor nodes and push successor nodes
		//to queue. 
		successors := getSuccessors(x)
		for i := 0; i < len(successors); i++{
			if _, ok := explored[successors[i].state]; ok {
				continue
			}
			explored[x.state] = true
			successors[i].priority = misplacedTiles(successors[i].state, goalState)
			heap.Push(&frontier, successors[i])
		}
	}
}

	
func aStarSearchManhattan(initial state) (string, int, int, time.Duration){
	goalMap := make(map[int][2]int)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			goalMap[goalState[i][j]] = [2]int{i, j}
		}
	}	
	explored := make(map[state]bool)
	var frontier PriorityQueue
	heap.Init(&frontier)

	var temp [2]int
	temp[0], temp[1] = getBlankPos(initial)
	//push the initial node to the queue
	n := &node{
		state: initial,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}
	
	start := time.Now()
	heap.Push(&frontier, *n)
	for{
		if frontier.isEmpty() { return "", -1, -1, time.Since(start)  }//fail
		x := heap.Pop(&frontier).(node)

		if x.state == goalState{
			//return action without leading "|" for easy splitting 
			return strings.TrimPrefix(x.action, "|"), x.cost, len(explored), time.Since(start) 
		}
		
		//expand the node generating successor nodes and push successor nodes
		//to queue. 
		successors := getSuccessors(x)
		for i := 0; i < len(successors); i++{
			//if successor is already explored get another successor
			if _, ok := explored[successors[i].state]; ok {
				continue
			}
			explored[x.state] = true
			successors[i].priority = successors[i].cost + manhattanDistance(successors[i].state, goalState, goalMap)
			//manhattanDistance(successors[i].state, goalState, goalMap)
			//misplacedTiles(successors[i].state, goalState)
			heap.Push(&frontier, successors[i])
		}
	}
}


func aStarSearchMisplaced(initial state) (string, int, int, time.Duration){
	goalMap := make(map[int][2]int)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			goalMap[goalState[i][j]] = [2]int{i, j}
		}
	}	
	explored := make(map[state]bool)
	var frontier PriorityQueue
	heap.Init(&frontier)

	var temp [2]int
	temp[0], temp[1] = getBlankPos(initial)
	//push the initial node to the queue
	n := &node{
		state: initial,
		action: "",
		cost: 0,
		priority: 0,
		blankPos: temp,
	}
	
	start := time.Now()
	heap.Push(&frontier, *n)
	for{
		if frontier.isEmpty() { return "", -1, -1, time.Since(start)  }//fail
		x := heap.Pop(&frontier).(node)

		if x.state == goalState{
			//return action without leading "|" for easy splitting 
			return strings.TrimPrefix(x.action, "|"), x.cost, len(explored), time.Since(start) 
		}
				
		//expand the node generating successor nodes and push successor nodes
		//to queue. 
		successors := getSuccessors(x)
		for i := 0; i < len(successors); i++{
			//if successor is already explored get another successor
			if _, ok := explored[successors[i].state]; ok {
				continue
			}
			explored[x.state] = true
			successors[i].priority = successors[i].cost + misplacedTiles(successors[i].state, goalState)
			heap.Push(&frontier, successors[i])
		}
	}
}


var e = 0 //keep track of expanded nodes

//IDA* often has a lot more nodes expanded compared to breadth-first search
//because it uses recursive depth-first search. It does not keep track of expanded
//nodes resulting into nodes getting expanded all over again.
func IDaStarSearchManhattan(initial state) (string, int, int, time.Duration) {
	e = 0
	goalMap := make(map[int][2]int)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			goalMap[goalState[i][j]] = [2]int{i, j}
		}
	}	
	var temp [2]int
	temp[0], temp[1] = getBlankPos(initial)
	n := &node{
		state: initial,
		action: "",
		cost: 0,
		priority: manhattanDistance(initial, goalState, goalMap),
		blankPos: temp,
	}

	start := time.Now()
	bound := n.priority
	for{
		x := searchManhattan(*n, 0, bound, goalMap)
		if x.state == goalState{ return strings.TrimPrefix(x.action, "|"), x.cost, e, time.Since(start) }
		bound++
	}
}

func searchManhattan(n node, g int, bound int, goalMap map[int][2]int) (node){
	
	h := manhattanDistance(n.state, goalState, goalMap)
	if n.state == goalState{ return n }
	f := n.cost + h
	if f > bound{ return n }

	successors := getSuccessors(n)
	for i := 0; i < len(successors); i++{
		e++
		x := searchManhattan(successors[i], successors[i].cost, bound, goalMap)
		if x.state == goalState{ return x }
	}
	return n
}


//IDA* with misplaced tile heuristic can sometimes find a solution and sometimes
//needs a lot of time compared to IDA* with manhattan distance heuristic
func IDaStarSearchMisplaced(initial state) (string, int, int, time.Duration) {
	e = 0
	var temp [2]int
	temp[0], temp[1] = getBlankPos(initial)
	n := &node{
		state: initial,
		action: "",
		cost: 0,
		priority: misplacedTiles(initial, goalState),
		blankPos: temp,
	}

	start := time.Now()
	bound := n.priority
	for{
		x := searchMisplaced(*n, 0, bound)
		if x.state == goalState{ return strings.TrimPrefix(x.action, "|"), x.cost, e, time.Since(start) }
		bound++
	}
}

func searchMisplaced(n node, g int, bound int) (node){
	
	h := misplacedTiles(n.state, goalState)
	if n.state == goalState{ return n }
	f := n.cost + h
	if f > bound{ return n }

	successors := getSuccessors(n)
	for i := 0; i < len(successors); i++{
		e++
		x := searchMisplaced(successors[i], successors[i].cost, bound)
		if x.state == goalState{ return x }
	}
	return n
}