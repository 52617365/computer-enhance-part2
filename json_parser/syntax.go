package main

import (
	"strconv"
	"strings"
)

// NOTE: We probably don't need to build an AST to parse JSON but I'll do it anyway just in case and to learn.

type ObjectNode struct {
	Objects  map[string]Node
	startPos int
	endPos   int
}

type ArrayNode struct {
	Elements []Node
	startPos int
	endPos   int
}

type StringNode struct {
	Value    string
	startPos int
	endPos   int
}

type NumberNode struct {
	Value    float64
	startPos int
	endPos   int
}

type Parser struct {
	tokens []Token
	pos    int
}

type EndOfFile struct {
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

	start := p.pos

	for p.tokens[p.pos].tokenType != QUOTATION { // FIXME: p.pos out of range here.
		parsedString = parsedString + p.tokens[p.pos].tokenContents

		p.IncrementPos()
		// p.pos++
	}

	if p.tokens[p.pos].tokenType != QUOTATION {
		panic("Expected the end of a string (\") here.")
	}

	end := p.pos // Capturing the end of the string

	p.IncrementPos() // Skipping the closing " character.

	return StringNode{
		Value:    parsedString,
		startPos: start,
		endPos:   end,
	}

}

func (p *Parser) parseArray() ArrayNode {
	if p.tokens[p.pos].tokenType != SQUAREOPEN {
		panic("Expected the current token type to be SQUAREOPEN.")
	}

	var elements []Node

	p.IncrementPos() // Skipping the [ character.
	// p.pos++ // Skipping the [ character.

	start := p.pos

	for p.tokens[p.pos].tokenType != SQUARECLOSE {

		node := p.parse()

		elements = append(elements, node)

		p.IncrementPos()
		// p.pos++
	}

	if p.tokens[p.pos].tokenType != SQUARECLOSE {
		panic("Expected the current token type to be SQUARECLOSE.")
	}

	end := p.pos

	return ArrayNode{
		Elements: elements,
		startPos: start,
		endPos:   end,
	}
}

func (p *Parser) parseNumber() NumberNode {
	var parsedNumber string

	start := p.pos

	for p.tokens[p.pos].tokenType != COMMA && p.tokens[p.pos].tokenType != CURLYCLOSE {
		parsedNumber = p.tokens[p.pos].tokenContents
		// parsedNumber = parsedNumber + p.tokens[p.pos].tokenContents

		p.IncrementPos()
	}

	if p.tokens[p.pos].tokenType != COMMA && p.tokens[p.pos].tokenType != CURLYCLOSE {
		panic("Expected the end of a number (, or }) here.")
	}

	end := p.pos // Capturing the end of the number

	castedFloat, _ := strconv.ParseFloat(strings.TrimSpace(parsedNumber), 64)

	return NumberNode{
		Value:    castedFloat,
		startPos: start,
		endPos:   end,
	}

}

func (p *Parser) parseObject() ObjectNode {
	if p.tokens[p.pos].tokenType != CURLYOPEN {
		panic("Expected the current token type to be CURLYOPEN")
	}
	pairs := make(map[string]Node)

	p.IncrementPos() // Skipping the { character.

	start := p.pos

	for p.tokens[p.pos].tokenType != CURLYCLOSE {

		key := p.parse()
		value := p.parse()

		if p.tokens[p.pos].tokenType == COMMA {
			key := p.parse()
			keyCast, found := key.(StringNode)

			if !found {
				panic("Expected a string here.")
			}

			node := p.parse()

			pairs[keyCast.Value] = node
		}

		keyCast, found := key.(StringNode)
		if !found {
			panic("Expected a string here.")
		}
		pairs[keyCast.Value] = value

		p.IncrementPos()
	}

	end := p.pos

	if p.tokens[p.pos].tokenType != CURLYCLOSE {
		panic("Expected the current token type to be CURLYCLOSE")
	}

	return ObjectNode{
		Objects:  pairs,
		startPos: start,
		endPos:   end,
	}
}

func (p *Parser) parse() Node {
	if p.pos >= len(p.tokens) {
		return EndOfFile{}
	}
	token := p.tokens[p.pos]

	switch token.tokenType {
	case CURLYOPEN:
		return p.parseObject()
	case CURLYCLOSE:
		p.IncrementPos()
		return p.parse()
	case SQUAREOPEN:
		return p.parseArray()
	case SQUARECLOSE:
		p.IncrementPos()
		return p.parse()
	case QUOTATION:
		return p.parseString()
	case IDENT:
		return p.parseString()
	case NUMBER:
		return p.parseNumber()
	case COLON:
		p.IncrementPos() // Skipping the colon because we don't actually care about it.
		return p.parse()
	case COMMA:
		p.IncrementPos() // Skipping the comma because we don't actually care about it in this context.
		return p.parse()
	default:
		return nil
	}
}
