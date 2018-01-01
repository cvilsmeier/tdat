package tdat

import (
	"bytes"
	"github.com/cvilsmeier/tdat/assert"
	"testing"
)

func TestLexer(t *testing.T) {
	{
		input := "\n"
		input += "persons\n"
		input += "|id:n   |name:s   \n"
		input += "|1   |\"joe\"\n"
		input += "|2|\"\"\n"
		input += "| | \n"
		input += "||\n"
		input += "  "
		lex := newLexer(bytes.NewBufferString(input))
		//
		assertNext(t, lex, "1:1 newline()")
		// persons
		assertNext(t, lex, "2:1 text(persons)")
		assertNext(t, lex, "2:8 newline()")
		// |id:n   |name:s   |carId:n
		assertNext(t, lex, "3:1 separator()")
		assertNext(t, lex, "3:2 text(id:n)")
		assertNext(t, lex, "3:9 separator()")
		assertNext(t, lex, "3:10 text(name:s)")
		assertNext(t, lex, "3:19 newline()")
		// |1   |"joe"
		assertNext(t, lex, "4:1 separator()")
		assertNext(t, lex, "4:2 text(1)")
		assertNext(t, lex, "4:6 separator()")
		assertNext(t, lex, "4:7 text(joe)")
		assertNext(t, lex, "4:12 newline()")
		// |2|""
		assertNext(t, lex, "5:1 separator()")
		assertNext(t, lex, "5:2 text(2)")
		assertNext(t, lex, "5:3 separator()")
		assertNext(t, lex, "5:4 text()")
		assertNext(t, lex, "5:6 newline()")
		// | |
		assertNext(t, lex, "6:1 separator()")
		assertNext(t, lex, "6:3 separator()")
		assertNext(t, lex, "6:5 newline()")
		// ||
		assertNext(t, lex, "7:1 separator()")
		assertNext(t, lex, "7:2 separator()")
		assertNext(t, lex, "7:3 newline()")
		// EOF
		assertNext(t, lex, "8:3 eof()")
		assertNext(t, lex, "8:3 eof()")
	}
	{
		input := "\r\n"
		input += "\t\t\r\n"
		input += "persons    \n"
		input += "  cars\n"
		input += "  \n"
		input += "\n"
		lex := newLexer(bytes.NewBufferString(input))
		assertNext(t, lex, "1:2 newline()")
		assertNext(t, lex, "2:4 newline()")
		assertNext(t, lex, "3:1 text(persons)")
		assertNext(t, lex, "3:12 newline()")
		assertNext(t, lex, "4:3 text(cars)")
		assertNext(t, lex, "4:7 newline()")
		assertNext(t, lex, "5:3 newline()")
		assertNext(t, lex, "6:1 newline()")
		assertNext(t, lex, "7:1 eof()")
		assertNext(t, lex, "7:1 eof()")
		assertNext(t, lex, "7:1 eof()")
	}
	{
		input := "| | \r\n"
		input += " | | \r\n"
		input += "\t|\t|\t\r\n"
		input += "||\r\n"
		input += "||\n"
		input += "|-|-\n"
		lex := newLexer(bytes.NewBufferString(input))
		// input := "| | \r\n"
		assertNext(t, lex, "1:1 separator()")
		assertNext(t, lex, "1:3 separator()")
		assertNext(t, lex, "1:6 newline()")
		// input += " | | \r\n"
		assertNext(t, lex, "2:2 separator()")
		assertNext(t, lex, "2:4 separator()")
		assertNext(t, lex, "2:7 newline()")
		// input += "\t|\t|\t\r\n"
		assertNext(t, lex, "3:2 separator()")
		assertNext(t, lex, "3:4 separator()")
		assertNext(t, lex, "3:7 newline()")
		// input += "||\r\n"
		assertNext(t, lex, "4:1 separator()")
		assertNext(t, lex, "4:2 separator()")
		assertNext(t, lex, "4:4 newline()")
		// input += "||\n"
		assertNext(t, lex, "5:1 separator()")
		assertNext(t, lex, "5:2 separator()")
		assertNext(t, lex, "5:3 newline()")
		// input += "|-|-\n"
		assertNext(t, lex, "6:1 separator()")
		assertNext(t, lex, "6:2 text(-)")
		assertNext(t, lex, "6:3 separator()")
		assertNext(t, lex, "6:4 text(-)")
		assertNext(t, lex, "6:5 newline()")
		// EOF
		assertNext(t, lex, "7:1 eof()")
		assertNext(t, lex, "7:1 eof()")
	}
}

func assertNext(t *testing.T, lex *lexer, exp string) {
	t.Helper()
	tok, err := lex.next()
	act := ""
	if tok != nil {
		act += tok.String()
	}
	if err != nil {
		act += err.Error()
	}
	assert.EqStr(t, exp, act)
}

func TestLexerReadText(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		exp   string
	}{
		{"text_2", "A", "A"},
		{"text_11", "abc  \n", "abc"},
		{"text_12", "2017-01-01  |", "2017-01-01"},
		{"text_14", "-|", "-"},
		{"text_32", "a\"b\"c", "a\"b\"c"},
		{"text_33", "a\"b\"|", "a\"b\""},
		{"text_34", "a\"b\"\r\n", "a\"b\""},
		{"text_53", "\"\"\n", ""},
		{"text_61", "\"blabla\"\n", "blabla"},
		{"text_62", "\"|\"\n", "|"},
		{"text_63", "\"\\\"\"\n", "\""},
		{"text_64", "\"\\u2602\"\n", "☂"},
		{"text_71", "\"\\u2602 ☂\"  |", "☂ ☂"},
		{"text_72", "\"\"ab\"\"  \n", ""},
		{"text_73", "\"hello\" \u0006", "hello"},
		{"err_11", "\"a", "line 1, pos 3: unterminated string"},
		{"err_12", "\"☂", "line 1, pos 3: unterminated string"},
		{"err_13", "\"\r", "line 1, pos 3: unterminated string"},
		{"err_14", "\"\\", "line 1, pos 3: unterminated escape sequence"},
		{"err_21", "\"\\u12", "line 1, pos 6: unterminated escape sequence"},
		{"err_22", "\"\\e", "line 1, pos 3: illegal escape sequence"},
		{"err_23", string([]byte{2}), "line 1, pos 1: invalid char 0x2"},
		{"err_24", string([]byte{'a', 1}), "line 1, pos 2: invalid char 0x1"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			lex := newLexer(bytes.NewBufferString(testCase.input))
			text, err := "", error(nil)
			if testCase.input[0] == '"' {
				text, err = lex.readQuotedText()
			} else {
				text, err = lex.readText()
			}
			act := text
			if err != nil {
				act += err.Error()
			}
			assert.EqStr(t, testCase.exp, act)
		})
	}
}

func BenchmarkLexerNext(b *testing.B) {
	var input []byte
	{
		b := &bytes.Buffer{}
		b.WriteString("persons\n")
		b.WriteString("|id:i|rate:f|flag:b|name:s\n")
		for i := 0; i < 200*1000; i++ {
			b.WriteString("|1|13000.00|true|\"joe\"\n")
		}
		input = b.Bytes()
	}
	lex := newLexer(bytes.NewBuffer(input))
	b.ResetTimer()
	for {
		tok, err := lex.next()
		if err != nil {
			panic(err)
		}
		if tok.ttype == eofToken {
			break
		}
	}
}
