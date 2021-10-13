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

type Lexer struct {
	r *bufio.Reader
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		r: bufio.NewReader(r),
	}
}

func (l *Lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (l *Lexer) unread() {
	if err := l.r.UnreadRune(); err != nil {
		// if we reach here, there is a problem in our logic
		panic(err)
	}
}

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
