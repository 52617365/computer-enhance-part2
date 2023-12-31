package main

import (
	"bufio"
	"io"
	"unicode"
)

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func (l *Lexer) GetTokensFromLexer() []Token {
	var tokens []Token
	for {
		token := l.Lex()
		if token.tokenType == EOF {
			break
		}

		tokens = append(tokens, token)
	}
	return tokens
}

func (l *Lexer) Lex() Token {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return Token{
					pos:           l.pos,
					tokenType:     EOF,
					tokenContents: "",
				}
			}

			// at this point there isn't much we can do, and the lexer
			// should just return the raw error to the user
			panic(err)
		}
		l.pos.column++

		switch r {
		case '\n':
			l.resetPosition()
		case '{':
			return Token{
				pos:           l.pos,
				tokenType:     CURLYOPEN,
				tokenContents: "{",
			}
		case '}':
			return Token{
				pos:           l.pos,
				tokenType:     CURLYCLOSE,
				tokenContents: "}",
			}
		case '[':
			return Token{
				pos:           l.pos,
				tokenType:     SQUAREOPEN,
				tokenContents: "[",
			}
		case ']':
			return Token{
				pos:           l.pos,
				tokenType:     SQUARECLOSE,
				tokenContents: "]",
			}
		case ':':
			return Token{
				pos:           l.pos,
				tokenType:     COLON,
				tokenContents: ":",
			}
		case ',':
			return Token{
				pos:           l.pos,
				tokenType:     COMMA,
				tokenContents: ",",
			}
		case '"':
			return Token{
				pos:           l.pos,
				tokenType:     QUOTATION,
				tokenContents: "\"",
			}
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsDigit(r) || unicode.IsLetter(r) || l.isIdentSymbol(r) {

				// backup and let lexIdent rescan the beginning of the ident
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()

				isNumber := true

				for _, r := range lit {
					if !unicode.IsNumber(r) && r != '.' && r != '-' {
						isNumber = false
						break
					}
				}

				if isNumber {
					return Token{
						pos:           startPos,
						tokenType:     NUMBER,
						tokenContents: lit,
					}
				} else if lit == "true" || lit == "false" {
					return Token{
						pos:           startPos,
						tokenType:     BOOLEAN,
						tokenContents: lit,
					}
				} else {
					return Token{
						pos:           startPos,
						tokenType:     IDENT,
						tokenContents: lit,
					}
				}
			} else {
				return Token{
					pos:           l.pos,
					tokenType:     ILLEGAL,
					tokenContents: string(r),
				}
			}
		}
	}
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			} else {
				panic("Why did we get here?")
			}
		}

		l.pos.column++

		if unicode.IsLetter(r) || unicode.IsNumber(r) || l.isIdentSymbol(r) || unicode.IsSpace(r) {
			lit = lit + string(r)
		} else if r == ',' {
			nextChar, err := l.reader.Peek(1)
			if err != nil {
				if err == io.EOF {
					// at the end of the identifier
					return lit
				} else {
					panic("Why did we get here?")
				}
			}

			if nextChar[0] == ' ' || unicode.IsLetter(rune(nextChar[0])) || unicode.IsNumber(rune(nextChar[0])) {
				lit = lit + string(r)
			} else {
				// scanned something not in the identifier
				return lit
			}
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) isIdentSymbol(r rune) bool {
	if r == '.' || r == '-' || r == '_' || r == '$' || r == '+' || r == '!' || r == '(' || r == ')' {
		return true
	} else {
		return false
	}
}
