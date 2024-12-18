package main

import (
	"SudokuSolverGolang/loadboard"
	"fmt"
	"strconv"
)


type Board struct {
	High int
	Width int 
	Content [][]Cell
}


type Cell struct {
	Value int
	isBlocked bool
}


func (board *Board) Init() {
	board.Content = make([][]Cell, board.High)
	for i := range board.Content {
		board.Content[i] = make([]Cell, board.Width)
	}
}


func (board *Board) Print() {
	fmt.Println("----------------------------")
	for i := 0; i < board.High; i++ {
		for j := 0; j < board.Width; j++ {
			value := board.Content[i][j].Value
			if j == 0 {
				fmt.Print("| ")
			}
			if value == 0 {
				fmt.Print("-")
			} else {
				if (board.Content[i][j].isBlocked) {
					fmt.Print("\033[35m")
					fmt.Print(strconv.Itoa(value))
					fmt.Print("\033[0m")
				} else {
					fmt.Print(strconv.Itoa(value))
				}
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
	if !(x < board.High && y < board.Width && 0 < value && value <= board.Width) {
		err := fmt.Errorf("invalid value or position: value=%d, x=%d, y=%d", value, x, y)
		panic(err)
	} else if (board.Content[x][y].isBlocked) { 
		err := fmt.Errorf("invalid position: x=%d, y=%d. Cell is blocked", x, y)
		panic(err)
	}
	board.Content[x][y].Value = value
	return value, nil
}


func (board *Board) FillBoard (valuesToFill [][]int) {
	for i := 0; i < board.High; i++ {
		for j := 0; j < board.Width; j++ {
			board.Content[i][j].Value = valuesToFill[i][j]
			if board.Content[i][j].Value == 0 {
				board.Content[i][j].isBlocked = false
			} else {
				board.Content[i][j].isBlocked = true
			}
		}
	}
}
 
func (board *Board) GetValuesInSameSquare (y int, x int) (map[int]struct{}) {
	init_row := (x / 3) * 3
	init_column := (y / 3) * 3
	usedValues := map[int]struct{}{}
	for i := init_row; i < init_row + 3; i++ {
        for j := init_column; j < init_column + 3; j++ {
			if _, ok := usedValues[board.Content[i][j].Value]; !ok {
				if board.Content[i][j].Value != 0 {
					usedValues[board.Content[i][j].Value] = struct{}{}
				}
			}
        }
    }
	return usedValues
}


func (board *Board) GetPossibleValues(y int, x int) (map[int]struct{}) {
	usedValues := map[int]struct{}{}
	possibleValues := map[int]struct{}{
		1: {}, 2: {}, 3: {},
		4: {}, 5: {}, 6: {},
		7: {}, 8: {}, 9: {},
	}
	for i := 0; i < 9; i++ {
		if board.Content[x][i].Value != 0 {
			if _, ok := usedValues[board.Content[x][i].Value]; !ok {
				usedValues[board.Content[x][i].Value] = struct{}{}
			}
		}
		if board.Content[i][y].Value != 0 {
			if _, ok := usedValues[board.Content[i][y].Value]; !ok {
				usedValues[board.Content[i][y].Value] = struct{}{}
			}
		}
	}
	inSameSquare := board.GetValuesInSameSquare(y, x)
	for key := range inSameSquare {
        usedValues[key] = struct{}{}
    }
	for key := range usedValues {
		if _, ok := possibleValues[key]; ok {
			delete(possibleValues, key)
		}
    }
	delete(usedValues, 0)
	return possibleValues
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
	for k := 0; k < 10; k++ {
		for i := 0; i < board.High; i++ {
			for j := 0; j < board.Width; j++ {
				if !board.Content[i][j].isBlocked && board.Content[i][j].Value == 0 {
					possibleValues := board.GetPossibleValues(j, i)
					if len(possibleValues) == 1 {
						target_value := 0
						for key := range possibleValues {
							target_value = key
							break
						}
						board.SetValue(target_value, i, j)
					}
				}
			}
		}
		fmt.Println("============================")
		board.Print()
	}
}
