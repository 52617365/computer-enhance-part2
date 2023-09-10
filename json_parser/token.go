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
  QUOTATION
  EOF
  IDENT
  ILLEGAL
)

var tokens = []string {
  CURLYOPEN: "{",
  CURLYCLOSE: "}",
  COLON: ":",
  SQUAREOPEN: "[",
  SQUARECLOSE: "]",
  COMMA: ",",
  QUOTATION: "\"",
  IDENT: "IDENT",
  ILLEGAL: "ILLEGAL",
}

func (t TokenType) String() string {
  if t == EOF {
    panic("Can't lookup string of EOF")
  }
  return tokens[t]
}

