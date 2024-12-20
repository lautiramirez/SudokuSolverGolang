package loadboard

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
	"encoding/json"
)

type ResponseAPI struct {
    Values [][]int  `json:"medium"`
}

func GetValuesFromAPI (urlAPI string) (values [][]int) {
    response, err := http.Get(urlAPI)

    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }
    responseData, err := io.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    var sudokuResponse ResponseAPI
    err = json.Unmarshal(responseData, &sudokuResponse)
    if err != nil {
        fmt.Println("Error al parsear el JSON:", err)
        return
    }
    return sudokuResponse.Values
}