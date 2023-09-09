package main

import (
  "flag"
  "fmt"
  "os"
)

type Arguments struct {
  filePath string
}

func parseArguments() Arguments {
  var filePath string

  flag.StringVar(&filePath, "filePath", "", "Path to the parsable JSON file.")
  flag.Parse()


  if filePath == "" {
    panic("file path is not provided")
  }

  return Arguments {
    filePath: filePath,
  }
}

func readFileToString(filePath string) string {
  data, err := os.ReadFile(filePath)

  if err != nil {
    panic("Error reading the JSON file contents")
  }

  dataString := string(data)

  return dataString
}

func main() {
  arguments := parseArguments()

  data := readFileToString(arguments.filePath)

  fmt.Printf("-------------\n")
  fmt.Printf("File path: %s\n", arguments.filePath)
  fmt.Printf("Data:\n%s\n", data)
  fmt.Printf("-------------\n")
}


 
func parseJson(jsonData string) map[string]any {


  return map[string]interface{}{"hello": "test"}
}
