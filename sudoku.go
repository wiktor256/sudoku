package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type Board [81]byte

func main() {
	start := time.Now()

	var prettyArg bool
	flag.BoolVar(&prettyArg, "pretty", false, "Display solution board")
	var puzzleArg string
	flag.StringVar(&puzzleArg, "puzzle", "", "<puzzle input - 81 digits>")
	flag.Parse()

	re := regexp.MustCompile("^[0-9]{81}$")

	if !re.MatchString(puzzleArg) {
		fmt.Println("Puzzle input should be a list of 81 digits!")
		os.Exit(2)
	}

	var puzzle Board

	for i, c := range puzzleArg {
		b, _ := strconv.ParseInt(string(c), 10, 32)
		puzzle[i] = byte(b)
	}

	var puzzles = [4]Board{puzzle}
	puzzles[1] = rotateBoard(puzzles[0])
	puzzles[2] = rotateBoard(puzzles[1])
	puzzles[3] = rotateBoard(puzzles[2])

	var solution Board
	var iter uint64

	var wg sync.WaitGroup
	var once sync.Once
	wg.Add(1)

	for i := range puzzles {
		go func(i int) {
			board, it := solvePuzzle(puzzles[i])
			if i > 0 {
				for r := i; r < 4; r++ {
					board = rotateBoard(board)
				}
			}
			once.Do(func () {
				solution = board
				iter = it
				wg.Done()
			})
		}(i)
	}

	wg.Wait()


	var hasSolution = true
	for i := 0; i < 81; i++ {
		if solution[i] == 0 {
			hasSolution = false
		}
	}

	if !hasSolution {
		fmt.Println("No solution found")
		os.Exit(2)
	}

	end := time.Now()
	duration:= end.Sub(start)

	if prettyArg {
		printBoard(solution, &puzzle)
		fmt.Printf("Iterations: %d\nTime: %v\n", iter, duration)
	} else {
		var buffer bytes.Buffer
		for i := 0; i < len(solution); i++ {
			buffer.WriteString(strconv.Itoa(int(solution[i])))
		}
		fmt.Printf("%s\t%d\t%v\n", buffer.String(), iter, duration)
	}
}

type Bits uint16
type ValidationMask [9]Bits

const (
	D0 Bits = 0
	D1 Bits = 1 << iota
	D2
	D3
	D4
	D5
	D6
	D7
	D8
	D9
)
var digits = [10]Bits{D0, D1, D2, D3, D4, D5, D6, D7, D8, D9}

var squareOf = [81]int{
	0, 0, 0, 1, 1, 1, 2, 2, 2,
	0, 0, 0, 1, 1, 1, 2, 2, 2,
	0, 0, 0, 1, 1, 1, 2, 2, 2,
	3, 3, 3, 4, 4, 4, 5, 5, 5,
	3, 3, 3, 4, 4, 4, 5, 5, 5,
	3, 3, 3, 4, 4, 4, 5, 5, 5,
	6, 6, 6, 7, 7, 7, 8, 8, 8,
	6, 6, 6, 7, 7, 7, 8, 8, 8,
	6, 6, 6, 7, 7, 7, 8, 8, 8,
}

func solvePuzzle(puzzle Board) (Board, uint64) {
	var squares ValidationMask
	var rows ValidationMask
	var cols ValidationMask
	var d Bits

	// Mark preset values for the puzzle
	for i, row := 0, 0; row < 9; row++ {
		for col := 0; col < 9; col, i = col+1, i+1 {
			digit := puzzle[i]
			if digit == byte(0) {
				continue
			}
			d = digits[digit]
			sq := squareOf[i]
			if squares[sq]&d != 0 || rows[row]&d != 0 || cols[col]&d != 0 { // has
				fmt.Printf("Invalid puzzle [%d, %d] = %d\n", row, col, digit)
				os.Exit(2)
			}
			squares[sq] |= d // set
			rows[row] |= d
			cols[col] |= d
		}
	}

	// clone the puzzle board so we can keep updating it while
	// the original puzzle can be used to identify preset cells
	var board Board
	for i, d := range puzzle {
		board[i] = d
	}
	// Backtracking algorithm.
	var backtrack = 1
	var iter uint64 = 0
	var row, col, sq int
	var curr byte

	for i := 0; i >= 0 && i < 81; i+= backtrack {
		col = i % 9
		row = (i - col) / 9

		// skip cell preset in the puzzle
		if puzzle[i] != byte(0) {
			continue
		}

		iter++

		// index of 3x3 square
		sq = squareOf[i]
		// current value in the cell
		curr = board[i]

		if curr != byte(0) {
			// clear the current value since it didn't work out (we got here through backtracking)
			board[i] = byte(0)
			d = digits[curr]
			squares[sq] &= ^d // clear
			rows[row] &= ^d
			cols[col] &= ^d
		}

		for val := curr + 1; val <= 9; val++ {
			d = digits[val]
			if rows[row]&d == 0 && squares[sq]&d == 0 && cols[col]&d == 0 { // !has
				// found possible match
				board[i] = val
				squares[sq] |= d
				rows[row] |= d
				cols[col] |= d
				break
			}
		}

		if board[i] == byte(0) {
			backtrack = -1
		} else {
			backtrack = 1
		}
	}

	return board, iter
}


func rotateBoard(board Board) Board {
	var rotated Board
	for i := 0; i < 81; i++ {
		col := i % 9
		row := (i - col) / 9
		rotated[(8 - col) * 9 + row] = board[i]
	}
	return rotated
}

func printBoard(board Board, mask *Board) {
	for i, row := 0, 0; row < 9; row++ {
		var buffer bytes.Buffer
		for col := 0; col < 9; col, i = col+1, i+1 {
			highlight := false
			if mask == nil {
				highlight = board[i] != 0
			} else {
				highlight = mask[i] != 0
			}
			format := " %d "
			if highlight {
				format = "\x1b[7m %d \x1b[0m"
			}
			buffer.WriteString(fmt.Sprintf(format, board[i]))
		}
		fmt.Println(buffer.String())
	}
}
