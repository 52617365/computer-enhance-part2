package main

import (
	"fmt"
	"strconv"
	"strings"
)

// NOTE: We probably don't need to build an AST to parse JSON but I'll do it anyway just in case and to learn.
// TODO: Why are we never getting the object node in our main calling code? Makes zero sense.

type ObjectNode struct {
	nodeType        string
	Objects         map[string]Node
	startPos        Position
	endPos          Position
	tokenIndexStart int
	tokenIndexEnd   int
}

type ArrayNode struct {
	nodeType        string
	Elements        []Node
	startPos        Position
	endPos          Position
	tokenIndexStart int
	tokenIndexEnd   int
}

type StringNode struct {
	nodeType        string
	Value           string
	startPos        Position
	endPos          Position
	tokenIndexStart int
	tokenIndexEnd   int
}

type NumberNode struct {
	nodeType        string
	Value           float64
	startPos        Position
	endPos          Position
	tokenIndexStart int
	tokenIndexEnd   int
}

type Parser struct {
	tokens []Token
	syntax []Node
	pos    int
}
type EndOfFile struct {
	endPos int
}

func (p *Parser) IncrementPos() {
	if len(p.tokens) > p.pos {
		p.pos++
	}
}

// Node interface to represent AST nodes
type Node interface{}

func printContents(n Node) {
	if eof, ok := n.(EndOfFile); ok {
		fmt.Printf("EOF - File end at position %d\n", eof.endPos)
		return
	} else if contents, ok := n.(ObjectNode); ok {
		fmt.Printf("%s - start: %d:%d, end: %d:%d, token index range: %d:%d\n", contents.nodeType, contents.startPos.column, contents.startPos.line, contents.endPos.column, contents.endPos.line, contents.tokenIndexStart, contents.tokenIndexEnd)

		for k, v := range contents.Objects {
			fmt.Printf("%s: ", k)
			printContents(v)
		}

	} else if contents, ok := n.(ArrayNode); ok {
		fmt.Printf("%s - start: %d:%d, end: %d:%d, token index range: %d:%d\n", contents.nodeType, contents.startPos.column, contents.startPos.line, contents.endPos.column, contents.endPos.line, contents.tokenIndexStart, contents.tokenIndexEnd)

		for k, v := range contents.Elements {
			fmt.Printf("\nArray [%d]: ", k)
			printContents(v)
			fmt.Printf("----")
		}

	} else if contents, ok := n.(StringNode); ok {
		fmt.Printf("%s\t", contents.Value)
		fmt.Printf("%s - start: %d:%d, end: %d:%d, token index range: %d:%d\n", contents.nodeType, contents.startPos.column, contents.startPos.line, contents.endPos.column, contents.endPos.line, contents.tokenIndexStart, contents.tokenIndexEnd)
	} else if contents, ok := n.(NumberNode); ok {
		fmt.Printf("%.15f\t", contents.Value)
		fmt.Printf("%s - start: %d:%d, end: %d:%d, token index range: %d:%d\n", contents.nodeType, contents.startPos.column, contents.startPos.line, contents.endPos.column, contents.endPos.line, contents.tokenIndexStart, contents.tokenIndexEnd)
	}
}

func GetParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) parseString() StringNode {
	var parsedString string

	p.IncrementPos() // Getting rid of the opening " character.

	startLine := p.tokens[p.pos].pos.line
	startColumn := p.tokens[p.pos].pos.column
	startIndex := p.pos

	for p.tokens[p.pos].tokenType != QUOTATION { // FIXME: p.pos out of range here.
		parsedString = parsedString + p.tokens[p.pos].tokenContents

		p.IncrementPos()
	}

	if p.tokens[p.pos].tokenType != QUOTATION {
		panic("Expected the end of a string (\") here.")
	}

	endLine := p.tokens[p.pos].pos.line
	endColumn := p.tokens[p.pos].pos.column
	endIndex := p.pos

	p.IncrementPos() // Getting rid of the closing " character.

	return StringNode{
		nodeType:        "string",
		Value:           parsedString,
		startPos:        Position{line: startLine, column: startColumn},
		endPos:          Position{line: endLine, column: endColumn},
		tokenIndexStart: startIndex,
		tokenIndexEnd:   endIndex,
	}

}

func (p *Parser) parseArray() Node {
	if p.tokens[p.pos].tokenType != SQUAREOPEN {
		panic("Expected the current token type to be SQUAREOPEN.")
	}

	var elements []Node

	p.IncrementPos() // Getting rid of the [ character.

	startLine := p.tokens[p.pos].pos.line
	startColumn := p.tokens[p.pos].pos.column
	startIndex := p.pos

	for p.tokens[p.pos].tokenType != SQUARECLOSE {

		node := p.parse()

		// return early if we hit end of file with parse
		if _, ok := node.(EndOfFile); ok {
			panic(fmt.Sprintf("endPos: %d, error: %s", p.pos, "Expected a closing square bracket but got EOF."))
		}

		elements = append(elements, node)

		if p.tokens[p.pos].tokenType == COMMA {
			p.IncrementPos()
		}

	}

	if p.tokens[p.pos].tokenType != SQUARECLOSE {
		panic(fmt.Sprintf("endPos: %d, error: %s", p.pos, "Expected the current token type to be SQUARECLOSE."))
	}

	endLine := p.tokens[p.pos].pos.line
	endColumn := p.tokens[p.pos].pos.column
	endIndex := p.pos

	p.IncrementPos() // Getting rid of the ] character.
	return ArrayNode{
		nodeType:        "array",
		Elements:        elements,
		startPos:        Position{line: startLine, column: startColumn},
		endPos:          Position{line: endLine, column: endColumn},
		tokenIndexStart: startIndex,
		tokenIndexEnd:   endIndex,
	}
}

func (p *Parser) parseNumber() NumberNode {
	var parsedNumber string

	startLine := p.tokens[p.pos].pos.line
	startColumn := p.tokens[p.pos].pos.column
	startIndex := p.pos

	for p.tokens[p.pos].tokenType != COMMA && p.tokens[p.pos].tokenType != CURLYCLOSE && p.tokens[p.pos].tokenType != SQUARECLOSE {
		parsedNumber = p.tokens[p.pos].tokenContents

		p.IncrementPos()
	}

	if p.tokens[p.pos].tokenType != COMMA && p.tokens[p.pos].tokenType != CURLYCLOSE && p.tokens[p.pos].tokenType != SQUARECLOSE {
		panic("Expected the end of a number (, or } or ]) here.")
	}

	endLine := p.tokens[p.pos].pos.line
	endColumn := p.tokens[p.pos].pos.column
	endIndex := p.pos

	castedFloat, _ := strconv.ParseFloat(strings.TrimSpace(parsedNumber), 64)

	return NumberNode{
		nodeType:        "number",
		Value:           castedFloat,
		startPos:        Position{line: startLine, column: startColumn},
		endPos:          Position{line: endLine, column: endColumn},
		tokenIndexStart: startIndex,
		tokenIndexEnd:   endIndex,
	}

}

func (p *Parser) parseObject() Node {
	if p.tokens[p.pos].tokenType != CURLYOPEN {
		panic("Expected the current token type to be CURLYOPEN")
	}

	pairs := make(map[string]Node)

	p.IncrementPos() // Getting rid of the { character.

	startLine := p.tokens[p.pos].pos.line
	startColumn := p.tokens[p.pos].pos.column
	startIndex := p.pos

	for p.tokens[p.pos].tokenType != CURLYCLOSE {

		key := p.parse()
		value := p.parse()

		if _, ok := key.(EndOfFile); ok {
			panic("Reached EOF before closing curly bracket")
		}

		if _, ok := value.(EndOfFile); ok {
			panic("Reached EOF before closing curly bracket")
		}

		keyCast, found := key.(StringNode)
		if !found {
			panic("Expected key to be a string.")
		}

		pairs[keyCast.Value] = value

		if p.tokens[p.pos].tokenType == COMMA {
			p.IncrementPos()
		}
	}

	endLine := p.tokens[p.pos].pos.line
	endColumn := p.tokens[p.pos].pos.column
	endIndex := p.pos

	if p.tokens[p.pos].tokenType != CURLYCLOSE {
		panic("Expected the current token type to be CURLYCLOSE")
	}

	p.IncrementPos() // Getting rid of the } character.

	return ObjectNode{
		nodeType:        "object",
		Objects:         pairs,
		startPos:        Position{line: startLine, column: startColumn},
		endPos:          Position{line: endLine, column: endColumn},
		tokenIndexStart: startIndex,
		tokenIndexEnd:   endIndex,
	}
}

func (p *Parser) parse() Node {
	token := p.tokens[p.pos]

	switch token.tokenType {
	case CURLYOPEN:
		parsedObject := p.parseObject()
		return parsedObject
	case SQUAREOPEN:
		parsedArray := p.parseArray()
		return parsedArray
	case QUOTATION:
		parsedString := p.parseString()
		return parsedString
	case IDENT:
		parsedString := p.parseString()
		return parsedString
	case NUMBER:
		parsedNumber := p.parseNumber()
		return parsedNumber
	case COLON:
		p.IncrementPos() // Skipping the colon because we don't actually care about it.
		return p.parse()
	default:
		panic("Why did we get here?")
	}
}
