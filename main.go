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


func (board *Board) GetPossibleValuesColumn (j int) (map[int][]int) {
	valuesInColumn := map[int][]int{}
	for k := 0; k < board.Width; k++ {
		if board.Content[k][j].Value == 0 && !board.Content[k][j].isBlocked {
			possibleValues := board.GetPossibleValues(k, j)
			for key := range possibleValues {
				valuesInColumn[key] = append(valuesInColumn[key], k)
			}
		}
	}
	return valuesInColumn
}


func (board *Board) GetPossibleValuesRow (i int) (map[int][]int) {
	valuesInRow := map[int][]int{}
	for k := 0; k < board.Width; k++ {
		if board.Content[i][k].Value == 0 && !board.Content[i][k].isBlocked {
			possibleValues := board.GetPossibleValues(i, k)
			for key := range possibleValues {
				valuesInRow[key] = append(valuesInRow[key], k)
			}
		}
	}
	return valuesInRow
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
			} else if board.Content[i][j].isBlocked {
				fmt.Print("\033[35m")
				fmt.Print(strconv.Itoa(value))
				fmt.Print("\033[0m")
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


func (board *Board) GetValuesInSameSquare (i int, j int) (map[int]struct{}) {
	init_row := (i / 3) * 3
	init_column := (j / 3) * 3
	usedValues := map[int]struct{}{}
	for k_i := init_row; k_i < init_row + 3; k_i++ {
        for k_j := init_column; k_j < init_column + 3; k_j++ {
			if _, ok := usedValues[board.Content[k_i][k_j].Value]; !ok && board.Content[k_i][k_j].Value != 0 {
				usedValues[board.Content[k_i][k_j].Value] = struct{}{}
			}
        }
    }
	return usedValues
}


func (board *Board) GetPossibleValues(i int, j int) (map[int]struct{}) {
	usedValues := map[int]struct{}{}
	possibleValues := map[int]struct{}{
		1: {}, 2: {}, 3: {},
		4: {}, 5: {}, 6: {},
		7: {}, 8: {}, 9: {},
	}
	for k := 0; k < 9; k++ {
		if _, ok := usedValues[board.Content[i][k].Value]; !ok && board.Content[i][k].Value != 0 {
			usedValues[board.Content[i][k].Value] = struct{}{}
		}
		if _, ok := usedValues[board.Content[k][j].Value]; !ok && board.Content[k][j].Value != 0 {
			usedValues[board.Content[k][j].Value] = struct{}{}
		}
	}
	inSameSquare := board.GetValuesInSameSquare(i, j)
	for key := range inSameSquare {
        usedValues[key] = struct{}{}
    }
	delete(usedValues, 0)
	for key := range usedValues {
		delete(possibleValues, key)
    }
	return possibleValues
}


func (board *Board) FillWithUniqueValues () {
	for i := 0; i < board.High; i++ {
		for j := 0; j < board.Width; j++ {
			if !board.Content[i][j].isBlocked && board.Content[i][j].Value == 0 {
				possibleValues := board.GetPossibleValues(i, j)
				if len(possibleValues) == 1 {
					for key := range possibleValues {
						_, err := board.SetValue(key, i, j)
						if err != nil {
							panic(err)
						}
					}
				}
			}
		}
		rowValues := board.GetPossibleValuesRow(i)
		for key := range rowValues {
			if len(rowValues[key]) == 1 {
				_, err := board.SetValue(key, i, rowValues[key][0])
				if err != nil {
					panic(err)
				}
			}
		}
		columnValues := board.GetPossibleValuesColumn(i)
		for key := range columnValues {
			if len(columnValues[key]) == 1 {
				_, err := board.SetValue(key, columnValues[key][0], i)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}


func (board *Board) SetWithGuessFiller () {

}


func main() {

	board := Board{
		High:  9,
		Width: 9,
	}
	board.Init()

	valuesToFillBoard := loadboard.GetValuesFromAPI("https://sudoku-game-and-api.netlify.app/api/sudoku")
	board.FillBoard(valuesToFillBoard)
	board.Print()

	for k := 0; k < 5; k++ {
		board.FillWithUniqueValues()

		fmt.Println("============================")
		board.Print()
	}
}
