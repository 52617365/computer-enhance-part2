package main

import (
  "flag"
  "fmt"
  "os"
  "io"
  "bufio"
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


  file, err := os.Open(arguments.filePath)
  if err != nil {
    panic(err)
  }
  defer file.Close()

  lexer := Lexer{}.New(file)

  results := parseJson(data)
}


 
func parseJson(jsonData string) map[string]any {
  return map[string]interface{}{"hello": "test"}
}


type Position struct {
  line int
  column int
}

type Lexer struct {
  pos Position
  reader *bufio.Reader
}



func (l Lexer) New(reader io.Reader) *Lexer {
  return &Lexer{
    pos: Position{line: 1, column: 0},
    reader: bufio.NewReader(reader),
  }
}


// func lex(jsonData string) Token {
//     
//
// }
