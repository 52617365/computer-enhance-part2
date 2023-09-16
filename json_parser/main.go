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

	parser.parse()

  // FIXME: with objects, for example pairs: {...} we have not properly implemented the fact that
  //  the object is actually the value of pairs.
	for _, node := range parser.syntax {
		if eof, ok := node.(EndOfFile); ok {
			fmt.Printf("EOF - File end at position %d\n", eof.endPos)
			break
		} else if contents, ok := node.(ObjectNode); ok {
			fmt.Printf("Object - start: %d, end: %d, objects: %s\n", contents.startPos, contents.endPos, contents.Objects)
		} else if contents, ok := node.(ArrayNode); ok {
			fmt.Printf("Array - start: %d, end: %d, array: %s\n", contents.startPos, contents.endPos, contents.Elements)
		} else if contents, ok := node.(StringNode); ok {
			fmt.Printf("String - start: %d, end: %d, string: %s\n", contents.startPos, contents.endPos, contents.Value)
		} else if contents, ok := node.(NumberNode); ok {
			fmt.Printf("Number - start: %d, end: %d, number: %.15f\n", contents.startPos, contents.endPos, contents.Value)
		}
	}
}
