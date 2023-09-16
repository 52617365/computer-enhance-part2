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

	return Arguments{
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

	// data := readFileToString(arguments.filePath)

	fmt.Printf("-------------\n")
	fmt.Printf("File path: %s\n", arguments.filePath)
	// fmt.Printf("Data:\n%s\n", data)
	fmt.Printf("-------------\n")

	file, err := os.Open(arguments.filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lexer := NewLexer(file)

	tokens := lexer.GetTokensFromLexer()

	parser := GetParser(tokens)

  // FIXME: right now we're expecting the parse function to be called multiple times but it's only called once since it's recursive.
  // POTENTIAL FIX: Maybe we should add an AST member into the Parser struct that gets updated?
  // UPDATE: Added syntax member that is []Node
  parser.parse()

  for _, node := range parser.syntax {
    if eof, ok := node.(EndOfFile); ok {
      fmt.Printf("File end at position %d\n", eof.endPos)
      break
    } else if contents, ok := node.(ObjectNode); ok {
      fmt.Printf("start: %d, end: %d, objects: %s\n", contents.startPos, contents.endPos, contents.Objects)
    } else if contents, ok := node.(ArrayNode); ok {
      fmt.Printf("start: %d, end: %d, array: %s\n", contents.startPos, contents.endPos, contents.Elements)
    } else if contents, ok := node.(StringNode); ok {
      fmt.Printf("start: %d, end: %d, string: %s\n", contents.startPos, contents.endPos, contents.Value)
    } else if contents, ok := node.(NumberNode); ok {
      fmt.Printf("start: %d, end: %d, number: %s\n", contents.startPos, contents.endPos, contents.Value)
    }
  }
}
