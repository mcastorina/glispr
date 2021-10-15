package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Token int

func (t Token) String() string {
	switch t {
	case tokError:
		return "<error>"
	case tokAtom:
		return "<atom>"
	case tokMinus:
		return "-"
	case tokNumber:
		return "<number>"
	case tokString:
		return "<string>"
	case tokComment:
		return "<comment>"
	case tokLParen:
		return "("
	case tokRParen:
		return ")"
	case tokQuote:
		return "'"
	case tokEOF:
		return "<EOF>"
	}
	return fmt.Sprintf("<unknown token %d>", t)
}

const (
	// Token types
	tokError Token = iota
	tokAtom
	tokMinus
	tokNumber // non-negative
	tokString
	tokComment
	tokLParen
	tokRParen
	tokQuote
	tokEOF

	eof rune = 0
)

// Lexer reads inputs and produces a stream of Tokens and their matching string literals
type Lexer struct {
	r *bufio.Reader
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		r: bufio.NewReader(r),
	}
}

// read reads a rune from the input buffer
func (l *Lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread puts the last read rune back into the buffer
func (l *Lexer) unread() {
	if err := l.r.UnreadRune(); err != nil {
		// if we reach here, there is a problem in our logic
		panic(err)
	}
}

// Next determines the next Token and string literal. It returns tokEOF,
// "" when the input buffer is empty.
func (l *Lexer) Next() (Token, string) {
	ch := l.read()
	switch ch {
	case eof:
		return tokEOF, ""
	case '(':
		return tokLParen, string(ch)
	case ')':
		return tokRParen, string(ch)
	case '\'':
		return tokQuote, string(ch)
	case '-':
		return tokMinus, string(ch)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.unread()
		return l.nextNumber()
	case ';':
		l.unread()
		return l.nextComment()
	case '"':
		l.unread()
		return l.nextString()
	case ' ', '\t', '\n':
		for ch == ' ' || ch == '\t' || ch == '\n' {
			ch = l.read()
		}
		l.unread()
		return l.Next()
	default:
		l.unread()
		return l.nextAtom()
	}
}

// nextNumber continuously reads input characters in [0-9_] and returns
// the full string
func (l *Lexer) nextNumber() (Token, string) {
	var buf bytes.Buffer

	for ch := l.read(); ch != eof; ch = l.read() {
		if ch >= '0' && ch <= '9' || ch == '_' {
			_, _ = buf.WriteRune(ch)
		} else {
			l.unread()
			break
		}
	}
	return tokNumber, buf.String()
}

// nextComment continuously reads input characters until a newline and returns
// the full string
func (l *Lexer) nextComment() (Token, string) {
	var buf bytes.Buffer

	for ch := l.read(); ch != eof; ch = l.read() {
		if ch == '\n' {
			break
		}
		_, _ = buf.WriteRune(ch)
	}
	return tokComment, buf.String()
}

// nextString continuously reads input characters until a double quote
// and returns the full string. This method accounts for backslashes in
// the string and will not add the backslash to the output.
func (l *Lexer) nextString() (Token, string) {
	var buf bytes.Buffer

	_, _ = buf.WriteRune(l.read()) // write opening double quote
	for ch := l.read(); ch != eof; ch = l.read() {
		// handle backslash escape
		if ch == '\\' {
			if ch = l.read(); ch == eof {
				break
			}
			_, _ = buf.WriteRune(ch)
			continue
		} else if ch == '"' {
			_, _ = buf.WriteRune(ch)
			break
		}
		_, _ = buf.WriteRune(ch)
	}

	return tokString, buf.String()
}

// nextAtom continuously reads input characters until a valid ending
// character and returns the full string
func (l *Lexer) nextAtom() (Token, string) {
	var buf bytes.Buffer

	for ch := l.read(); ch != eof; ch = l.read() {
		if ch == '(' || ch == '"' || ch == ')' || ch == ';' || ch == ' ' || ch == '\t' || ch == '\n' {
			l.unread()
			break
		}
		_, _ = buf.WriteRune(ch)
	}

	return tokAtom, buf.String()
}
