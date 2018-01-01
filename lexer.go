package tdat

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type tokenType byte

const (
	textToken tokenType = iota + 1
	separatorToken
	newlineToken
	eofToken
)

func (tt tokenType) String() string {
	switch tt {
	case textToken:
		return "text"
	case separatorToken:
		return "separator"
	case newlineToken:
		return "newline"
	case eofToken:
		return "eof"
	}
	panic("unknown token type")
}

// A token is a lexeme.
type token struct {
	line  int
	pos   int
	ttype tokenType
	text  string
}

func (t *token) String() string {
	return fmt.Sprintf("%d:%d %s(%s)", t.line, t.pos, t.ttype, t.text)
}

// The lexer scans input for tokens. It iterates over the
// runes in the reader processes each rune one-by-one.
// Initially, the lexer points to the first rune.
// At the end of the reader (if the reader has no more runes to read),
// the lexer stops.
type lexer struct {
	reader io.RuneReader
	r      rune
	err    error
	line   int
	pos    int
}

func newLexer(reader io.RuneReader) *lexer {
	l := &lexer{
		reader,
		-1,
		nil,
		1,
		0,
	}
	l.read()
	return l
}

// Next parses and returns the next token or returns an error.
// If the lexer has stopped (after reading the last rune from the reader)
// next will always return a eof token.
// Since eof is considered normal, eof is never returned as error.
// The returned token must not be retained by the caller.
// The returned token is valid until the next invocation of next.
func (l *lexer) next() (*token, error) {
	// eat whitespace
	for {
		if l.err != nil {
			return nil, l.err
		}
		if l.r > ' ' || l.r == 0 || l.r == '\n' {
			break
		}
		l.read()
	}
	// scan next token
	switch l.r {
	case 0:
		return &token{l.line, l.pos, eofToken, ""}, nil
	case '|':
		line, pos := l.line, l.pos
		l.read()
		return &token{line, pos, separatorToken, ""}, nil
	case '\n':
		line, pos := l.line, l.pos
		l.read()
		return &token{line, pos, newlineToken, ""}, nil
	case '"':
		line, pos := l.line, l.pos
		text, err := l.readQuotedText()
		if err != nil {
			return nil, err
		}
		return &token{line, pos, textToken, text}, nil
	default:
		line, pos := l.line, l.pos
		text, err := l.readText()
		if err != nil {
			return nil, err
		}
		return &token{line, pos, textToken, text}, nil
	}
}

// readText will collect the next runes, up to (and not including) the next
// separator or newline or EOF, whichever comes first.
// It trims trailing whitespace.
func (l *lexer) readText() (string, error) {
	runes := make([]rune, 0, 40)
	for {
		if l.err != nil {
			return "", l.err
		}
		if l.r == 0 || l.r == '|' || l.r == '\n' {
			return strings.TrimSpace(string(runes)), nil
		}
		runes = append(runes, l.r)
		l.read()
	}
}

// readQuotedText will collect the next runes, up to the
// first unescaped double quote '"', which will close a quoted string.
// Escaping applies, unicode escaping also.
func (l *lexer) readQuotedText() (string, error) {
	runes := make([]rune, 0, 40)
	for {
		l.read()
		if l.err != nil {
			return "", l.err
		}
		switch l.r {
		case 0:
			return "", l.errorf("unterminated string")
		case '\\':
			r, err := l.readEscapeSequence()
			if err != nil {
				return "", err
			}
			runes = append(runes, r)
		case '"':
			l.read()
			return strings.TrimSpace(string(runes)), nil
		default:
			runes = append(runes, l.r)
		}
	}
}

func (l *lexer) readEscapeSequence() (rune, error) {
	l.read()
	switch {
	case l.err != nil:
		return 0, l.err
	case l.r == 0:
		return 0, l.errorf("unterminated escape sequence")
	case l.r == 'b':
		return '\b', nil
	case l.r == 't':
		return '\t', nil
	case l.r == 'n':
		return '\n', nil
	case l.r == 'f':
		return '\f', nil
	case l.r == 'r':
		return '\r', nil
	case l.r == 'u':
		return l.readUniodeEscapeSequence()
	case l.r == '"':
		return '"', nil
	case l.r == '\\':
		return '\\', nil
	}
	return 0, l.errorf("illegal escape sequence")
}

func (l *lexer) readUniodeEscapeSequence() (rune, error) {
	runes := []rune{'\\', 'u', 0, 0, 0, 0}
	for i := 2; i < 6; i++ {
		l.read()
		if l.err != nil {
			return 0, l.err
		}
		if l.r == 0 {
			return 0, l.errorf("unterminated escape sequence")
		}
		runes[i] = l.r
	}
	value, _, _, err := strconv.UnquoteChar(string(runes), 0)
	return value, err
}

// read advances the lexer by reading the next rune from the reader.
// If the reader has no more runes (eof), special rune 0 is set as
// the current rune and no error is set.
// read will also update line and pos.
// Rune codepoint U+0000 ocurring in the input is considered an error:
// Valid input must not contain '\0' characters.
func (l *lexer) read() {
	if l.r == 0 || l.err != nil {
		return
	}
	if l.r == '\n' {
		l.line++
		l.pos = 1
	} else {
		l.pos++
	}
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.r = 0
			l.err = nil
		} else {
			l.err = err
		}
		return
	}
	if r < 0x20 && (r != 0x09 && r != 0x0A && r != 0x0D) {
		err = l.errorf("invalid char 0x%x", r)
	}
	l.r = r
	l.err = err
}

func (l *lexer) errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("line %d, pos %d: %s", l.line, l.pos, msg)
}
