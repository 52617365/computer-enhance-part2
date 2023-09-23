package main

import (
	"strconv"
	"strings"
)

// NOTE: We probably don't need to build an AST to parse JSON but I'll do it anyway just in case and to learn.
// TODO: Why are we never getting the object node in our main calling code? Makes zero sense.

type ObjectNode struct {
  nodeType string
	Objects  map[string]Node
  startPos Position
  endPos Position
  tokenIndexStart int
  tokenIndexEnd int
}

type ArrayNode struct {
  nodeType string
	Elements []Node
  startPos Position
  endPos Position
  tokenIndexStart int
  tokenIndexEnd int
}

type StringNode struct {
  nodeType string
	Value    string
  startPos Position
  endPos Position
  tokenIndexStart int
  tokenIndexEnd int
}

type NumberNode struct {
  nodeType string
	Value    float64
  startPos Position
  endPos Position
  tokenIndexStart int
  tokenIndexEnd int
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

func GetParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) parseString() StringNode {
	var parsedString string

	p.IncrementPos() // Skipping the " character.

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

	p.IncrementPos() // Skipping the closing " character.

	return StringNode{
    nodeType: "string",
		Value:    parsedString,
    startPos: Position{line: startLine, column: startColumn},
    endPos: Position{line: endLine, column: endColumn},
    tokenIndexStart: startIndex,
    tokenIndexEnd: endIndex,
	}

}

func (p *Parser) parseArray() Node {
	if p.tokens[p.pos].tokenType != SQUAREOPEN {
		panic("Expected the current token type to be SQUAREOPEN.")
	}

	var elements []Node

	p.IncrementPos() // Skipping the [ character.

	startLine := p.tokens[p.pos].pos.line
	startColumn := p.tokens[p.pos].pos.column
  startIndex := p.pos

	for p.tokens[p.pos].tokenType != SQUARECLOSE {

		node := p.parse()

		// return early if we hit end of file with parse
		if _, ok := node.(EndOfFile); ok {
			return EndOfFile{endPos: p.pos}
		}

		elements = append(elements, node)

		p.IncrementPos()
	}

	if p.tokens[p.pos].tokenType != SQUARECLOSE {
		panic("Expected the current token type to be SQUARECLOSE.")
	}

	endLine := p.tokens[p.pos].pos.line
	endColumn := p.tokens[p.pos].pos.column
  endIndex := p.pos

	return ArrayNode{
    nodeType: "array",
		Elements: elements,
    startPos: Position{line: startLine, column: startColumn},
    endPos: Position{line: endLine, column: endColumn},
    tokenIndexStart: startIndex,
    tokenIndexEnd: endIndex,
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
    nodeType: "number",
		Value:    castedFloat,
    startPos: Position{line: startLine, column: startColumn},
    endPos: Position{line: endLine, column: endColumn},
    tokenIndexStart: startIndex,
    tokenIndexEnd: endIndex,
	}

}

func (p *Parser) parseObject() Node {
	if p.tokens[p.pos].tokenType != CURLYOPEN {
		panic("Expected the current token type to be CURLYOPEN")
	}

	pairs := make(map[string]Node)

	p.IncrementPos() // Skipping the { character.

	startLine := p.tokens[p.pos].pos.line
	startColumn := p.tokens[p.pos].pos.column
  startIndex := p.pos

	for p.tokens[p.pos].tokenType != CURLYCLOSE {

		key := p.parse()
		value := p.parse()

		// return early if we hit end of file with key parse.
		if _, ok := key.(EndOfFile); ok {
			return EndOfFile{endPos: p.pos}
		}

		// return early if we hit end of file with value parse.
		if _, ok := value.(EndOfFile); ok {
			return EndOfFile{endPos: p.pos}
		}

    // NOTE: This
	  keyCast, found := key.(StringNode)
    if !found {
      panic("Expected a string here.")
    }

    pairs[keyCast.Value] = value
    // NOTE: To this needs testing, I'm not sure if it's ok. It's the morning as I'm writing this.

		if p.tokens[p.pos].tokenType == COMMA {
			key = p.parse()

			// return early if we hit end of file with key parse.
			if _, ok := key.(EndOfFile); ok {
				return EndOfFile{endPos: p.pos}
			}

			keyCast, found := key.(StringNode)

			if !found {
				panic("Expected a string here.")
			}

			node := p.parse()

			// return early if we hit end of file
			if _, ok := node.(EndOfFile); ok {
				return EndOfFile{endPos: p.pos}
			}

			pairs[keyCast.Value] = node
		}

		p.IncrementPos()
	}

	endLine := p.tokens[p.pos].pos.line
	endColumn := p.tokens[p.pos].pos.column
  endIndex := p.pos

	if p.tokens[p.pos].tokenType != CURLYCLOSE {
		panic("Expected the current token type to be CURLYCLOSE")
	}

	return ObjectNode{
    nodeType: "object",
		Objects:  pairs,
    startPos: Position{line: startLine, column: startColumn},
    endPos: Position{line: endLine, column: endColumn},
    tokenIndexStart: startIndex,
    tokenIndexEnd: endIndex,
	}
}

func (p *Parser) parse() Node {
	if p.pos >= len(p.tokens) {
		p.syntax = append(p.syntax, EndOfFile{endPos: p.pos})
		return EndOfFile{endPos: p.pos}
	}

  // NOTE: We have checked that the results seem to be ok. Now we have to build the AST.

	token := p.tokens[p.pos]

	switch token.tokenType {
	case CURLYOPEN:
    parsedObject := p.parseObject()
    p.syntax = append(p.syntax, parsedObject)
		return parsedObject
	case CURLYCLOSE:
		p.IncrementPos()
		return p.parse()
	case SQUAREOPEN:
    parsedArray := p.parseArray()
    p.syntax = append(p.syntax, parsedArray)
		return parsedArray
	case SQUARECLOSE:
		p.IncrementPos()
		return p.parse()
	case QUOTATION:
    parsedString := p.parseString()
    p.syntax = append(p.syntax, parsedString)
    return parsedString
	case IDENT:
    parsedString := p.parseString()
    p.syntax = append(p.syntax, parsedString)
    return parsedString
	case NUMBER:
    parsedNumber := p.parseNumber()
    p.syntax = append(p.syntax, parsedNumber)
    return parsedNumber
	case COLON:
		p.IncrementPos() // Skipping the colon because we don't actually care about it.
		return p.parse()
	case COMMA:
		p.IncrementPos() // Skipping the comma because we don't actually care about it in this context.
		return p.parse()
	default:
		panic("Why did we get here?")
	}
}
