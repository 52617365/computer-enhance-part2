package main

import (
  "io"
  "bufio"
  "unicode"
)

type Position struct {
  line int
  column int
}

type Lexer struct {
  pos Position
  reader *bufio.Reader
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}



func NewLexer(reader io.Reader) *Lexer {
  return &Lexer{
    pos: Position{line: 1, column: 0},
    reader: bufio.NewReader(reader),
  }
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}
	
	l.pos.column--
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
			}
		}
			
    l.pos.column++
		if unicode.IsLetter(r) || unicode.IsNumber(r) || isIdentSymbol(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) Lex() (Position, TokenType, string) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			// at this point there  isn't much we can do, and the compiler
			// should just return the raw error to the user
			panic(err)
		}
    l.pos.column++

    switch r {
        case '\n':
          l.resetPosition()
        case '{':
            return l.pos, CURLYOPEN, "{"
        case '}':
            return l.pos, CURLYCLOSE, "}"
        case '[':
            return l.pos, SQUAREOPEN, "["
        case ']':
            return l.pos, SQUARECLOSE, "]"
        case ':':
            return l.pos, COLON, ":"
        case ',':
            return l.pos, COMMA, ","
        case '"':
            return l.pos, QUOTATION, "\""
        default:
            if unicode.IsSpace(r) {
                continue // nothing to do here, just move on
            } else if unicode.IsDigit(r) || unicode.IsLetter(r) || isIdentSymbol(r) {
                // backup and let lexIdent rescan the beginning of the ident
                startPos := l.pos
                l.backup()
                lit := l.lexIdent() // TODO: make lexIdent and handle numbers inside of it because we don't want to separate them.
                return startPos, IDENT, lit
            } else {
                return l.pos, ILLEGAL, string(r)
            }
        }
    }
}

func isIdentSymbol(r rune) bool {
  if r == '.' || r == '-' {
    return true
  } else {
    return false
  }
}
