package main

import (
     "os"
     "fmt"
     "flag"
     "math/rand"
     "encoding/json"
)

type Coordinate struct { 
        X0 float64 `json:"x0"`
        X1 float64 `json:"x1"`
        Y0 float64 `json:"y0"`
        Y1 float64 `json:"y1"`
}

type JsonPairs struct {
        Pairs []Coordinate `json:"pairs"`
}


func main() {
    var jsonNum int
    var seed int
    var outputToFile bool


    flag.IntVar(&jsonNum, "jsonNum", 10, "Number of JSON lines to generate")
    flag.IntVar(&seed, "seed", 929292988, "Seed to generate random numbers with")
    flag.BoolVar(&outputToFile, "outputToFile", true, "A flag that determines if the JSON should be outputted to a file. Defaults to true.")
    flag.Parse()

    fmt.Printf("Generating %d lines of JSON data\n", jsonNum)

    rand.Seed(int64(seed))


    var randomJsonData JsonPairs

    min := -90.0
    max := 90.0

    for i := 0; i < jsonNum; i++ {
        coordinates := Coordinate{
                X0: min + rand.Float64() * (max - min), // Random number between 0 and 180, then shift to range -90 to +90
                X1: min + rand.Float64() * (max - min), // Random number between 0 and 180, then shift to range -90 to +90
                Y0: min + rand.Float64() * (max - min), // Random number between 0 and 180, then shift to range -90 to +90
                Y1: min + rand.Float64() * (max - min), // Random number between 0 and 180, then shift to range -90 to +90
        }

        randomJsonData.Pairs = append(randomJsonData.Pairs, coordinates)
    }

    randomJsonDataEncoded, _ := json.Marshal(randomJsonData) 

    if outputToFile {
            fileName := fmt.Sprintf("Output/JsonOutput_%d.json", jsonNum)
            f, err := os.Create(fileName)
            if err != nil {
                    panic("Error creating a file for some reason.")
            }

            defer f.Close()

            _, err = f.WriteString(string(randomJsonDataEncoded))

            if err != nil {
                    panic("Error writing to file.")
            }
    } else {
            fmt.Printf("Final JSON data: %s\n", string(randomJsonDataEncoded))
    }


}
