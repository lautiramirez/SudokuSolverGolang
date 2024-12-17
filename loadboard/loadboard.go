package loadboard

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
	"encoding/json"
)

type ResponseAPI struct {
	NewBoard NewBoard `json:"newboard"`
}

type NewBoard struct {
	Grids []Grid `json:"grids"`
}

type Grid struct {
	Values [][]int `json:"value"`
}

func GetValuesFromAPI (urlAPI string) (values [][]int) {
    response, err := http.Get(urlAPI)

    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }
    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    var sudokuResponse ResponseAPI
    err = json.Unmarshal(responseData, &sudokuResponse)
    if err != nil {
        fmt.Println("Error al parsear el JSON:", err)
        return
    }
    return sudokuResponse.NewBoard.Grids[0].Values
}