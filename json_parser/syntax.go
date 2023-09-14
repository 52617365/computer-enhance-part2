package main

// NOTE: We probably don't need to build an AST to parse JSON but I'll do it anyway just in case and to learn.

type ObjectNode struct {
	Pairs map[string]Node
  startPos int
  endPos int
}

type ArrayNode struct {
	Elements []Node
  startPos int
  endPos int
}

type StringNode struct {
  Value string
  startPos int
  endPos int
}

type Parser struct {
	tokens []Token
	pos    int
}


// Node interface to represent AST nodes
type Node interface{}

func (p *Parser) parseString() StringNode {
  var parsedString string

  p.pos++ // Skipping the " character.

  start := p.pos

  for p.tokens[p.pos].tokenType != QUOTATION {
    parsedString = parsedString + p.tokens[p.pos].tokenContents

    p.pos++
  }

  end := p.pos // Capturing the end of the string

  p.pos++ // Skipping the closing " character.

  return StringNode{
    Value: parsedString,
    startPos: start,
    endPos: end,
  }

}

func (p *Parser) parseArray() ArrayNode {
  if p.tokens[p.pos].tokenType != SQUAREOPEN {
    panic("Expected the current token type to be SQUAREOPEN.")
  }

  var elements []Node

  p.pos++ // Skipping the [ character.

  start := p.pos

  for p.tokens[p.pos].tokenType != SQUARECLOSE {
    
    node := p.parse()

    elements = append(elements, node)

    p.pos++
  }

  end := p.pos

  if p.tokens[p.pos].tokenType != SQUARECLOSE {
    panic("Expected the current token type to be SQUARECLOSE.")
  }

  p.pos++ // Skipping the ] character at the end of the array.

  return ArrayNode{
    Elements: elements,
    startPos: start,
    endPos: end,
  }
}

func (p *Parser) parseObject() StringNode {
  if p.tokens[p.pos].tokenType != CURLYOPEN {
    panic("Expected the current token type to be CURLYOPEN")
  }

  var object map[string]Node

  p.pos++ // Skipping the { character.

  start := p.pos

  for p.tokens[p.pos].tokenType != SQUARECLOSE {
    node := p.parse()
    // TODO: We have to parse the key and value in here to produce the correct object.
    parsedString = parsedString + p.tokens[p.pos].tokenContents

    p.pos++
  }

  end := p.pos

  if p.tokens[p.pos].tokenType != CURLYCLOSE {
    panic("Expected the current token type to be CURLYCLOSE")
  }

  p.pos++ // Skipping the } character at the end of the object.

  return ObjectNode{
    Value: object,
    startPos: start,
    endPos: end,
  }
}
// TODO: define all the nodes. E.g. string, number, array, etc.
func (p *Parser) parse() Node {
	token := p.tokens[p.pos]
	p.pos++

	switch token.tokenType {
	case CURLYOPEN:
		return p.parseObject()
	case SQUAREOPEN:
		return p.parseArray()
  case QUOTATION: // FIXME: this one right here will not work because the value of the token is '"', 
                  // instead we will have to make a function that loops until it finds another " and captures that string instead.
    return p.parseString()
	default:
		return nil
	}
}
// type Node struct {
//     data Token
//     prev *Node
//     succ *Node
// }
//
// func buildSyntaxTreeFromTokens(tokens []Token) {
//
//   nodes := []Node
//
//
//   for _, token := range tokens {
//   }
//
// }
