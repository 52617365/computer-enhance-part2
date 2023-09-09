package main

type Token struct {
  tokenType TokenType
  contents string
}


type TokenType int64

const (
  CURLYOPEN TokenType = iota
  CURLYCLOSE
  COLON
  SQUAREOPEN
  SQUARECLOSE
  COMMA
)

var tokens = []string {
  CURLYOPEN: "{",
  CURLYCLOSE: "}",
  COLON: ":",
  SQUAREOPEN: "[",
  SQUARECLOSE: "]",
  COMMA: ",",
}

func (t TokenType) String() string {
  return tokens[t]
}

