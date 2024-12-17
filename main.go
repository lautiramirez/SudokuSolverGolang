package main

import (
	"fmt"
	"strconv"
	"SudokuSolverGolang/loadboard"
)

type Board struct {
	High int
	Width int 
	Content [][]int
}

func (board *Board) Init() {
	board.Content = make([][]int, board.High)
	for i := range board.Content {
		board.Content[i] = make([]int, board.Width)
	}
}

func (board *Board) Print() {
	fmt.Println("----------------------------")
	for i := 0; i < board.High; i++ {
		for j := 0; j < board.Width; j++ {
			value := board.Content[i][j]
			if j == 0 {
				fmt.Print("| ")
			}
			if value == 0 {
				fmt.Print("-")
			} else {
				fmt.Print(strconv.Itoa(value))
			}
			if (j + 1) % 3 == 0 {
				fmt.Print("  | ")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
		if (i + 1) % 3 == 0 {
			fmt.Println("----------------------------")
		}
	}
}

func (board *Board) SetValue(value int, x int, y int)(int, error){
	if !(x < board.High && y < board.Width && 0 < value && value < board.Width) {
		err := fmt.Errorf("invalid value or position: value=%d, x=%d, y=%d", value, x, y)
		panic(err)
	} 
	board.Content[x][y] = value
	return value, nil
}

func (board *Board) FillBoard (valuesToFill [][]int) {
	for i := 0; i < board.High; i++ {
		for j := 0; j < board.Width; j++ {
			board.Content[i][j] = valuesToFill[i][j]
		}
	}
}

func main() {

	board := Board{
		High:  9,
		Width: 9,
	}
	board.Init()
	
	valuesToFillBoard := loadboard.GetValuesFromAPI("https://sudoku-api.vercel.app/api/dosuku")
	board.FillBoard(valuesToFillBoard)
	board.Print()
	
	var x_position int
	var y_position int
	var new_value int

	for {
		fmt.Printf("Enter the position (y,x) and a value to fill: ")
		fmt.Scanln(&x_position, &y_position, &new_value)
		board.SetValue(new_value, x_position, y_position)
		fmt.Println()
		board.Print()
	}

}
