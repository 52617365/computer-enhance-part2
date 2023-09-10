package main


type Token struct {
  pos Position
  tokenType TokenType
  tokenContents string
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
  CURLYOPEN: "CURLYOPEN",
  CURLYCLOSE: "CURLYCLOSE",
  COLON: "COLON:",
  SQUAREOPEN: "SQUAREOPEN",
  SQUARECLOSE: "SQUARECLOSE",
  COMMA: "COMMA",
  QUOTATION: "QUOTATION",
  IDENT: "IDENT",
  ILLEGAL: "ILLEGAL",
}

func (t TokenType) String() string {
  if t == EOF {
    panic("Can't lookup string of EOF")
  }
  return tokens[t]
}

