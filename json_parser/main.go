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

	fmt.Printf("-------------\n")
	fmt.Printf("File path: %s\n", arguments.filePath)
	fmt.Printf("-------------\n")

	file, err := os.Open(arguments.filePath)
	if err != nil {
		panic(err)
	}

	// TODO: this appending and returning thing is causing hella issues for us.
	//  It's causing the items to be appended twice..
	defer file.Close()

	lexer := NewLexer(file)

	tokens := lexer.GetTokensFromLexer()

	parser := GetParser(tokens)

	parser.parse()

	// FIXME: with objects, for example pairs: {...} we have not properly implemented the fact that
	//  the object is actually the value of pairs.
	for _, node := range parser.syntax {
    if eof, ok := node.(EndOfFile); ok { 
      fmt.Printf("EOF - File end at position %d\n", eof.endPos) 
    break 
    } else if contents, ok := node.(ObjectNode); ok {
      fmt.Printf("%s - start: %d:%d, end: %d:%d, token index range: %d:%d, %s\n", contents.nodeType, contents.startPos.column, contents.endPos.line, contents.endPos.column, contents.tokenIndexStart, contents.tokenIndexEnd, contents.Objects)
    } else if contents, ok := node.(ArrayNode); ok {
      fmt.Printf("%s - start: %d:%d, end: %d:%d, token index range: %d:%d, %s\n", contents.nodeType, contents.startPos.column, contents.endPos.line, contents.endPos.column, contents.tokenIndexStart, contents.tokenIndexEnd, contents.Elements)
    } else if contents, ok := node.(StringNode); ok {
      fmt.Printf("%s - start: %d:%d, end: %d:%d, token index range: %d:%d, %s\n", contents.nodeType, contents.startPos.line, contents.startPos.column, contents.endPos.line, contents.endPos.column, contents.tokenIndexStart, contents.tokenIndexEnd, contents.Value)
    } else if contents, ok := node.(NumberNode); ok {
      fmt.Printf("%s - start: %d:%d, end: %d:%d, token index range: %d:%d, %.15f\n", contents.nodeType, contents.startPos.line, contents.startPos.column, contents.endPos.line, contents.endPos.column, contents.tokenIndexStart, contents.tokenIndexEnd, contents.Value)
    }
  }

}
