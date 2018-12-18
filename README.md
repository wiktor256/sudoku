# Sudoku solver written in Go

This is a sample Go implementation of Sudoku solver using backtracking algorithm.

The simple backtracking algorithm is very sensitive to the initial board orientation.
When running the backtracking algorithm with the top 500 most difficult sudoku puzzles,
the difference in number of iterations is enormous. 

For example the puzzle "#	2097	FNBTWXY	S8.f	32160" from http://www.sfsudoku.com/su17ExtremeDiff500.txt, 
in the original orientation requires 328,454,800 iterations, but rotated clockwise by 90 degrees, it only needs 666,106 iterations.
This is almost 500x improvement.

Because of the above observation, this program runs the backtracking algorithm in four parallel goroutines using different orientations of the puzzle board.
As soon any of the four goroutines solves the puzzle the program ends.

## Usage
```
sudoku -puzzle=000000051070030000800000000000501040030000600000800000500420000001000300000000700 -pretty
```

Output:
```
 3  6  9  7  8  4  2  5  1
 1  7  4  2  3  5  8  6  9
 8  5  2  1  9  6  4  3  7
 2  8  7  5  6  1  9  4  3
 4  3  5  9  7  2  6  1  8
 9  1  6  8  4  3  5  7  2
 5  9  3  4  2  7  1  8  6
 7  2  1  6  5  8  3  9  4
 6  4  8  3  1  9  7  2  5
Iterations: 2681344
Time: 41.18399ms
```
