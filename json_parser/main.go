package main

import (
  "os"
  "fmt"
  "flag"
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

  lexer := NewLexer(file)

  for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}

		fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)
	}
}
