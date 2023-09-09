package main

import (
     "os"
     "fmt"
     "flag"
     "math/rand"
     "encoding/json"
     "math"
)


type DistanceResults struct {
  
}

// this is outputted to a different file where we actually calculate the results.
type CoordinateDistance struct {
  Distance float64`json:"distance"`
}

type Coordinate struct { 
  X0 float64 `json:"x0"`
  X1 float64 `json:"x1"`
  Y0 float64 `json:"y0"`
  Y1 float64 `json:"y1"`
}

type JsonPairs struct {
  Pairs []Coordinate `json:"pairs"`
}


func Square(number float64) float64 {
  return number * number
}
func RadiansFromDegrees(degrees float64) float64  {
  return 0.01745329251994329577 * degrees
}

// earth radius 6372.8
func calculateDistance(x0 float64, y0 float64, x1 float64, y1 float64, earthRadius float64) float64 {
  lat1 := y0
  lat2 := y1

  lon1 := x0
  lon2 := x1


  distanceLat := RadiansFromDegrees(lat2 - lat1)
  distanceLon := RadiansFromDegrees(lon2 - lon1)

  lat1 = RadiansFromDegrees(lat1)
  lat2 = RadiansFromDegrees(lat2)

  a := Square(math.Sin(distanceLat/2.0)) + math.Cos(lat1)*math.Cos(lat2)*Square(math.Sin(distanceLon/2));
  c := 2.0 * math.Sin(math.Sqrt(a))

  result := earthRadius * c

  return result
}



func calculateAverageDistanceForHaversines(coordinates []Coordinate) float64 {
  countOfCoordinates := len(coordinates)
  
  var calculatedDistanceSum float64

  for _, coordinate := range coordinates {

    calculatedDistance := calculateDistance(coordinate.X0, coordinate.Y0, coordinate.X1, coordinate.Y1, 6372.8)
   
    calculatedDistanceSum += calculatedDistance
  }


  return calculatedDistanceSum / float64(countOfCoordinates)
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

    haversineAverage := calculateAverageDistanceForHaversines(randomJsonData.Pairs)

    if outputToFile {
      fileName := fmt.Sprintf("Output/JsonOutput_%d.json", jsonNum)
      writeToFile(fileName, string(randomJsonDataEncoded))

      fmt.Printf("Average haversine amount: %f\n", haversineAverage)
    } else {
      fmt.Printf("Average haversine amount: %f\n", haversineAverage)
      fmt.Printf("Final JSON data: %s\n", string(randomJsonDataEncoded))
    }

}

func writeToFile(fileName string, contents string) {
      f, err := os.Create(fileName)
      if err != nil {
              panic("Error creating json output file for some reason.")
      }

      defer f.Close()

      _, err = f.WriteString(contents)

      if err != nil {
          panic("Error writing to file.")
      }
}
