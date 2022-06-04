package main

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
)

const ROW = 11
const COL = 13
const COIN = -45
const BONUS = -120
const MONSTER = 30

type Pair struct {
	first, second int
}

type Cell struct {
	parent_i, parent_j int
	f, g, h, coinCount float64
}

type Entity struct {
	monster                                     []Pair
	bonusExists, daggerExists, mExists, pExists bool
	mDist                                       int
	mExpander                                   int
	player                                      Pair
}

//----------priority-queue-start----------
type Node struct { //for a priority queue
	f     float64 //priority
	pair  *Pair
	index int
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].f < pq[j].f //for the shortest path - '<', for the longest path - '>'
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Node)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

//----------priority-queue-end----------

func Make_Pair(row, col int) *Pair {
	p := Pair{
		first:  row,
		second: col,
	}
	return &p
}

func Make_pPair(fNew float64, p *Pair) *Node {
	pP := Node{
		f:    fNew,
		pair: p,
	}
	return &pP
}

func IsValid(row, col int) bool {
	return row >= 0 && row < ROW && col >= 0 && col < COL
}

func IsUnBlocked(grid [ROW][COL]rune, row, col int) bool {
	if grid[row][col] == '!' {
		return false
	}
	return true
}

func IsDestination(row, col int, dest Pair) bool {
	return row == dest.first && col == dest.second
}

func CalculateHValue(row, col int, dest Pair) float64 {
	return math.Abs(float64(row-dest.first)) + math.Abs(float64(col-dest.second))
}

func DistanceA(grid [ROW][COL]rune, src, dest Pair) int {
	if !IsValid(src.first, src.second) {
		return 1000
	}
	if !IsValid(dest.first, dest.second) {
		return 1000
	}
	if !IsUnBlocked(grid, src.first, src.second) || !IsUnBlocked(grid, dest.first, dest.second) {
		return 1000
	}
	if IsDestination(src.first, src.second, dest) {
		return 0
	}
	var closedList [ROW][COL]bool
	var cellDetails [ROW][COL]Cell
	var i, j int
	for i = 0; i < ROW; i++ {
		for j = 0; j < COL; j++ {
			cellDetails[i][j].f = math.MaxFloat64
		}
	}
	i = src.first
	j = src.second
	cellDetails[i][j].f = 0.0
	cellDetails[i][j].g = 0.0
	openList := make(PriorityQueue, 0, 100)
	heap.Push(&openList, Make_pPair(0.0, Make_Pair(i, j)))
	for openList.Len() > 0 {
		p := heap.Pop(&openList).(*Node)
		i = p.pair.first
		j = p.pair.second
		closedList[i][j] = true
		var gNew, hNew, fNew float64

		if IsValid(i-1, j) { //----------1st Successor (North)----------
			if IsDestination(i-1, j, dest) {
				gNew = cellDetails[i][j].g + 1.0
				return int(gNew)
			} else if !closedList[i-1][j] && IsUnBlocked(grid, i-1, j) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i-1, j, dest)
				fNew = gNew + hNew
				if cellDetails[i-1][j].f == math.MaxFloat64 || cellDetails[i-1][j].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i-1, j)))
					cellDetails[i-1][j].f = fNew
					cellDetails[i-1][j].g = gNew
				}
			}
		}

		if IsValid(i+1, j) { //----------2nd Successor (South)----------
			if IsDestination(i+1, j, dest) {
				gNew = cellDetails[i][j].g + 1.0
				return int(gNew)
			} else if !closedList[i+1][j] && IsUnBlocked(grid, i+1, j) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i+1, j, dest)
				fNew = gNew + hNew
				if cellDetails[i+1][j].f == math.MaxFloat64 || cellDetails[i+1][j].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i+1, j)))
					cellDetails[i+1][j].f = fNew
					cellDetails[i+1][j].g = gNew
				}
			}
		}

		if IsValid(i, j+1) { //----------3rd Successor (East)----------
			if IsDestination(i, j+1, dest) {
				gNew = cellDetails[i][j].g + 1.0
				return int(gNew)
			} else if !closedList[i][j+1] && IsUnBlocked(grid, i, j+1) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i, j+1, dest)
				fNew = gNew + hNew
				if cellDetails[i][j+1].f == math.MaxFloat64 || cellDetails[i][j+1].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i, j+1)))
					cellDetails[i][j+1].f = fNew
					cellDetails[i][j+1].g = gNew
				}
			}
		}

		if IsValid(i, j-1) { //----------4th Successor (West)----------
			if IsDestination(i, j-1, dest) {
				gNew = cellDetails[i][j].g + 1.0
				return int(gNew)
			} else if !closedList[i][j-1] && IsUnBlocked(grid, i, j-1) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i, j-1, dest)
				fNew = gNew + hNew
				if cellDetails[i][j-1].f == math.MaxFloat64 || cellDetails[i][j-1].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i, j-1)))
					cellDetails[i][j-1].f = fNew
					cellDetails[i][j-1].g = gNew
				}
			}
		}
	}
	return 1000
}

func TracePath(cellDetails [ROW][COL]Cell, dest Pair) *Pair {
	row := dest.first
	col := dest.second
	var nextRow, nextCol int

	for !(cellDetails[row][col].parent_i == row && cellDetails[row][col].parent_j == col) {
		tempRow := cellDetails[row][col].parent_i
		tempCol := cellDetails[row][col].parent_j
		nextRow = row
		nextCol = col
		row = tempRow
		col = tempCol
	}
	return Make_Pair(nextRow, nextCol)
}

func gValueSum(grid [ROW][COL]rune, e Entity, g float64, i, j int) float64 {
	result := g
	if e.pExists && e.player.first == i && e.player.second == j {
		result += 55 //
	} else if !e.bonusExists && grid[i][j] == 'b' {
		result += BONUS
	} else if grid[i][j] == '#' {
		result += COIN
	}
	if !e.daggerExists {
		k := e.mExpander
		for l := range e.monster {
			if i >= e.monster[l].first-k && i <= e.monster[l].first+k && j >= e.monster[l].second-k && j <= e.monster[l].second+k {
				result += MONSTER
			}
		}
	}
	return result
}

func fValueChange(grid [ROW][COL]rune, e Entity, i, j int) float64 {
	result := 0.0
	if grid[i][j] == '#' || grid[i][j] == 'd' || grid[i][j] == 'b' {
		result = result - 1
	}
	if !e.daggerExists {
		k := e.mExpander
		for l := range e.monster {
			if i >= e.monster[l].first-k && i <= e.monster[l].first+k && j >= e.monster[l].second-k && j <= e.monster[l].second+k {
				result += MONSTER
			}
		}
	}
	return result
}

//----------A* Search Algorithm start----------
func aStarSearch(grid [ROW][COL]rune, src, dest Pair, e Entity) *Pair {
	if IsDestination(src.first, src.second, dest) {
		return &src
	}
	var closedList [ROW][COL]bool

	var cellDetails [ROW][COL]Cell
	var i, j int

	for i = 0; i < ROW; i++ {
		for j = 0; j < COL; j++ {
			cellDetails[i][j].f = math.MaxFloat64
			cellDetails[i][j].g = math.MaxFloat64
			cellDetails[i][j].h = math.MaxFloat64
			cellDetails[i][j].parent_i = -1
			cellDetails[i][j].parent_j = -1
		}
	}
	i = src.first
	j = src.second
	cellDetails[i][j].f = 0.0
	cellDetails[i][j].g = 0.0
	cellDetails[i][j].h = 0.0
	cellDetails[i][j].parent_i = i
	cellDetails[i][j].parent_j = j

	openList := make(PriorityQueue, 0, 100)
	heap.Push(&openList, Make_pPair(0.0, Make_Pair(i, j)))
	for openList.Len() > 0 {
		p := heap.Pop(&openList).(*Node)
		i = p.pair.first
		j = p.pair.second
		closedList[i][j] = true

		var gNew, hNew, fNew float64

		if IsValid(i-1, j) {
			if IsDestination(i-1, j, dest) {
				cellDetails[i-1][j].parent_i = i
				cellDetails[i-1][j].parent_j = j
				nextMove := TracePath(cellDetails, dest)
				return nextMove
			} else if !closedList[i-1][j] && IsUnBlocked(grid, i-1, j) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i-1, j, dest)
				fNew = gNew + hNew + fValueChange(grid, e, i-1, j)

				if cellDetails[i-1][j].f == math.MaxFloat64 || cellDetails[i-1][j].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i-1, j)))
					cellDetails[i-1][j].f = fNew
					cellDetails[i-1][j].g = gNew
					cellDetails[i-1][j].h = hNew
					cellDetails[i-1][j].parent_i = i
					cellDetails[i-1][j].parent_j = j
				}
			}
		}

		if IsValid(i+1, j) {
			if IsDestination(i+1, j, dest) {
				cellDetails[i+1][j].parent_i = i
				cellDetails[i+1][j].parent_j = j
				nextMove := TracePath(cellDetails, dest)
				return nextMove
			} else if !closedList[i+1][j] && IsUnBlocked(grid, i+1, j) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i+1, j, dest)
				fNew = gNew + hNew + fValueChange(grid, e, i+1, j)

				if cellDetails[i+1][j].f == math.MaxFloat64 || cellDetails[i+1][j].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i+1, j)))
					cellDetails[i+1][j].f = fNew
					cellDetails[i+1][j].g = gNew
					cellDetails[i+1][j].h = hNew
					cellDetails[i+1][j].parent_i = i
					cellDetails[i+1][j].parent_j = j
				}
			}
		}

		if IsValid(i, j+1) {
			if IsDestination(i, j+1, dest) {
				cellDetails[i][j+1].parent_i = i
				cellDetails[i][j+1].parent_j = j
				nextMove := TracePath(cellDetails, dest)
				return nextMove
			} else if !closedList[i][j+1] && IsUnBlocked(grid, i, j+1) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i, j+1, dest)
				fNew = gNew + hNew + fValueChange(grid, e, i, j+1)

				if cellDetails[i][j+1].f == math.MaxFloat64 || cellDetails[i][j+1].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i, j+1)))
					cellDetails[i][j+1].f = fNew
					cellDetails[i][j+1].g = gNew
					cellDetails[i][j+1].h = hNew
					cellDetails[i][j+1].parent_i = i
					cellDetails[i][j+1].parent_j = j
				}
			}
		}

		if IsValid(i, j-1) {
			if IsDestination(i, j-1, dest) {
				cellDetails[i][j-1].parent_i = i
				cellDetails[i][j-1].parent_j = j
				nextMove := TracePath(cellDetails, dest)
				return nextMove
			} else if !closedList[i][j-1] && IsUnBlocked(grid, i, j-1) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i, j-1, dest)
				fNew = gNew + hNew + fValueChange(grid, e, i, j-1)

				if cellDetails[i][j-1].f == math.MaxFloat64 || cellDetails[i][j-1].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i, j-1)))
					cellDetails[i][j-1].f = fNew
					cellDetails[i][j-1].g = gNew
					cellDetails[i][j-1].h = hNew
					cellDetails[i][j-1].parent_i = i
					cellDetails[i][j-1].parent_j = j
				}
			}
		}
	}
	return &src
}

//----------Distance - A* start----------
func GoalCost(grid [ROW][COL]rune, e Entity, src, dest Pair) float64 {
	if !IsValid(src.first, src.second) {
		return 1000
	}
	if !IsValid(dest.first, dest.second) {
		return 1000
	}
	if !IsUnBlocked(grid, src.first, src.second) || !IsUnBlocked(grid, dest.first, dest.second) {
		return 1000
	}
	if IsDestination(src.first, src.second, dest) {
		return 0
	}
	var closedList [ROW][COL]bool
	var cellDetails [ROW][COL]Cell
	var i, j int
	for i = 0; i < ROW; i++ {
		for j = 0; j < COL; j++ {
			cellDetails[i][j].f = math.MaxFloat64
			cellDetails[i][j].g = math.MaxFloat64
			cellDetails[i][j].h = math.MaxFloat64
		}
	}
	i = src.first
	j = src.second
	cellDetails[i][j].f = 0.0
	cellDetails[i][j].g = 0.0
	cellDetails[i][j].h = 0.0
	cellDetails[i][j].coinCount = 0.0
	openList := make(PriorityQueue, 0, 100)
	heap.Push(&openList, Make_pPair(0.0, Make_Pair(i, j)))
	for openList.Len() > 0 {
		p := heap.Pop(&openList).(*Node)
		i = p.pair.first
		j = p.pair.second
		closedList[i][j] = true
		var gNew, hNew, fNew, cCount float64

		if IsValid(i-1, j) {
			if IsDestination(i-1, j, dest) {
				gNew = cellDetails[i][j].g + 1.0
				cCount = cellDetails[i][j].coinCount + gValueSum(grid, e, gNew, i-1, j)
				return cCount
			} else if !closedList[i-1][j] && IsUnBlocked(grid, i-1, j) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i-1, j, dest)
				fNew = gNew + hNew + fValueChange(grid, e, i-1, j)
				cCount = cellDetails[i][j].coinCount + gValueSum(grid, e, gNew, i-1, j)

				if cellDetails[i-1][j].f == math.MaxFloat64 || cellDetails[i-1][j].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i-1, j)))
					cellDetails[i-1][j].f = fNew
					cellDetails[i-1][j].g = gNew
					cellDetails[i-1][j].h = hNew
					cellDetails[i-1][j].coinCount = cCount
				}
			}
		}

		if IsValid(i+1, j) {
			if IsDestination(i+1, j, dest) {
				gNew = cellDetails[i][j].g + 1.0
				cCount = cellDetails[i][j].coinCount + gValueSum(grid, e, gNew, i+1, j)
				return cCount
			} else if !closedList[i+1][j] && IsUnBlocked(grid, i+1, j) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i+1, j, dest)
				fNew = gNew + hNew + fValueChange(grid, e, i+1, j)
				cCount = cellDetails[i][j].coinCount + gValueSum(grid, e, gNew, i+1, j)

				if cellDetails[i+1][j].f == math.MaxFloat64 || cellDetails[i+1][j].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i+1, j)))
					cellDetails[i+1][j].f = fNew
					cellDetails[i+1][j].g = gNew
					cellDetails[i+1][j].h = hNew
					cellDetails[i+1][j].coinCount = cCount
				}
			}
		}

		if IsValid(i, j+1) {
			if IsDestination(i, j+1, dest) {
				gNew = cellDetails[i][j].g + 1.0
				cCount = cellDetails[i][j].coinCount + gValueSum(grid, e, gNew, i, j+1)
				return cCount
			} else if !closedList[i][j+1] && IsUnBlocked(grid, i, j+1) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i, j+1, dest)
				fNew = gNew + hNew + fValueChange(grid, e, i, j+1)
				cCount = cellDetails[i][j].coinCount + gValueSum(grid, e, gNew, i, j+1)

				if cellDetails[i][j+1].f == math.MaxFloat64 || cellDetails[i][j+1].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i, j+1)))
					cellDetails[i][j+1].f = fNew
					cellDetails[i][j+1].g = gNew
					cellDetails[i][j+1].h = hNew
					cellDetails[i][j+1].coinCount = cCount
				}
			}
		}

		if IsValid(i, j-1) {
			if IsDestination(i, j-1, dest) {
				gNew = cellDetails[i][j].g + 1.0
				cCount = cellDetails[i][j].coinCount + gValueSum(grid, e, gNew, i, j-1)
				return cCount
			} else if !closedList[i][j-1] && IsUnBlocked(grid, i, j-1) {
				gNew = cellDetails[i][j].g + 1.0
				hNew = CalculateHValue(i, j-1, dest)
				fNew = gNew + hNew + fValueChange(grid, e, i, j-1)
				cCount = cellDetails[i][j].coinCount + gValueSum(grid, e, gNew, i, j-1)

				if cellDetails[i][j-1].f == math.MaxFloat64 || cellDetails[i][j-1].f > fNew {
					heap.Push(&openList, Make_pPair(fNew, Make_Pair(i, j-1)))
					cellDetails[i][j-1].f = fNew
					cellDetails[i][j-1].g = gNew
					cellDetails[i][j-1].h = hNew
					cellDetails[i][j-1].coinCount = cCount
				}
			}
		}
	}
	return 1000
}

//----------Next Move start----------
func NextMove(src, dest Pair) {
	if src.first-1 == dest.first {
		fmt.Println("up")
	} else if src.first+1 == dest.first {
		fmt.Println("down")
	} else if src.second-1 == dest.second {
		fmt.Println("left")
	} else if src.second+1 == dest.second {
		fmt.Println("right")
	} else {
		fmt.Println("stay")
	}
}

//----------Next Move end----------

//----------Check for Safety start----------
func IsSafe(grid *[ROW][COL]rune, me Pair, e Entity) bool {
	for i := range e.monster {
		dist := DistanceA(*grid, me, e.monster[i])
		if dist <= e.mDist || dist == 1000 {
			return false
		}
	}
	return true
}

//----------Check for Safety end----------

func IsAllowed(grid *[ROW][COL]rune, e Entity, row, col int) bool {
	if IsSafe(grid, *Make_Pair(row+1, col), e) {
		return true
	}
	if IsSafe(grid, *Make_Pair(row-1, col), e) {
		return true
	}
	if IsSafe(grid, *Make_Pair(row, col+1), e) {
		return true
	}
	if IsSafe(grid, *Make_Pair(row, col-1), e) {
		return true
	}
	return false
}

func LookForSafety(grid [ROW][COL]rune, e Entity, i, j int) *Pair {
	if IsValid(i, j+1) && IsUnBlocked(grid, i, j+1) && IsAllowed(&grid, e, i, j+1) {
		return Make_Pair(i, j+1)
	}
	if IsValid(i, j-1) && IsUnBlocked(grid, i, j-1) && IsAllowed(&grid, e, i, j-1) {
		return Make_Pair(i, j-1)
	}
	if IsValid(i+1, j) && IsUnBlocked(grid, i+1, j) && IsAllowed(&grid, e, i+1, j) {
		return Make_Pair(i+1, j)
	}
	if IsValid(i-1, j) && IsUnBlocked(grid, i-1, j) && IsAllowed(&grid, e, i-1, j) {
		return Make_Pair(i-1, j)
	}
	return nil
}

func anotherLFS(grid [ROW][COL]rune, e Entity, i, j int) *Pair { //use it after changing mDist to 2
	right := Make_Pair(i, j+1)
	if IsSafe(&grid, *right, e) {
		return right
	}
	left := Make_Pair(i, j-1)
	if IsSafe(&grid, *left, e) {
		return left
	}
	down := Make_Pair(i+1, j)
	if IsSafe(&grid, *down, e) {
		return down
	}
	up := Make_Pair(i-1, j)
	if IsSafe(&grid, *up, e) {
		return up
	}
	if IsValid(i, j+1) && IsUnBlocked(grid, i, j+1) && j+1 == 6 {
		return right
	}
	if IsValid(i, j-1) && IsUnBlocked(grid, i, j-1) && j-1 == 6 {
		return left
	}
	return nil
}

func main() {
	daggerExpire, bonusExpire := 0, 0
	for true {
		var w, h, playerID, tick int
		fmt.Scan(&w, &h, &playerID, &tick)
		var grid [11][13]rune
		var e Entity
		goals := make([]Pair, 0, 50)

		for i := 0; i < h; i++ {
			line := ""
			fmt.Scan(&line)
			aRune := []rune(line)
			for j, v := range aRune {
				grid[i][j] = v
				if v == '#' || v == 'b' || v == 'd' {
					goals = append(goals, *Make_Pair(i, j))
				}
			}
		}
		e.mDist = 3 //
		e.mExpander = 1

		var n, myX, myY int // number of entities
		fmt.Scan(&n)
		e.mExists, e.pExists = false, false
		pBonus := false
		for i := 0; i < n; i++ {
			var entType string
			var pID, x, y, param1, param2 int
			fmt.Scan(&entType, &pID, &x, &y, &param1, &param2)
			if entType == "p" && playerID == pID {
				myX = x
				myY = y
				if param1 == 1 && daggerExpire < 12 {
					daggerExpire++
					e.daggerExists = true
				} else if param1 != 1 {
					daggerExpire = 0
				}
				if param2 > 0 && bonusExpire < 12 {
					if param2 == 2 {
						e.mDist = 0
						e.mExpander = 0
					}
					bonusExpire++
					e.bonusExists = true
				} else if param2 == 0 {
					bonusExpire = 0
				}
			} else if entType == "m" {
				e.mExists = true
				mp := Make_Pair(y, x)
				e.monster = append(e.monster, *mp)
			} else if entType == "p" {
				e.pExists = true
				e.player = *Make_Pair(y, x)
				if param2 == 1 {
					pBonus = true
				}
			}
		}
		src := *Make_Pair(myY, myX)
		dest := *Make_Pair(myY, myX)

		//sort goals by dist from me
		sort.SliceStable(goals, func(i, j int) bool {
			d1 := DistanceA(grid, src, goals[i])
			d2 := DistanceA(grid, src, goals[j])
			return d1 < d2
		})

		min := math.MaxFloat64
		for i := range goals {
			gc1 := GoalCost(grid, e, src, goals[i]) //if a goal isn't reachable, a*Search returns src
			for j := i; j < len(goals); j++ {
				gc2 := GoalCost(grid, e, goals[i], goals[j])
				destValue := gc1 + gc2
				if destValue < min {
					min = destValue
					dest = goals[i]
				}
			}
		}

		if len(goals) <= 2 && pBonus && e.mExists { //don't let your opponent use his bonus
			dest = src
		}
		if tick >= 290 && len(goals) > 0 {
			dest = goals[0]
		}
		//fmt.Fprint(os.Stderr, dest, "<-dest\n")
		nextPath := aStarSearch(grid, src, dest, e)

		if !e.daggerExists && !IsAllowed(&grid, e, nextPath.first, nextPath.second) {
			nextPath = LookForSafety(grid, e, src.first, src.second)
			if nextPath == nil {
				e.mDist = 2
				nextPath = anotherLFS(grid, e, src.first, src.second)
			}
		}
		if nextPath == nil {
			fmt.Println("stay")
		} else {
			NextMove(src, *nextPath)
		}
	}
}
