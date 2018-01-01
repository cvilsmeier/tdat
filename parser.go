package tdat

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// ParseFromString is like ParseFromRuneReader but reads input from a string.
func ParseFromString(input string) (*Model, error) {
	runeReader := strings.NewReader(input)
	return ParseFromRuneReader(runeReader)
}

// ParseFromFile is like ParseFromRuneReader but reads input from a file.
func ParseFromFile(name string) (*Model, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	runeReader := bufio.NewReader(file)
	return ParseFromRuneReader(runeReader)
}

// ParseFromReader is like ParseFromRuneReader but reads input from a reader.
func ParseFromReader(reader io.Reader) (*Model, error) {
	runeReader := bufio.NewReader(reader)
	return ParseFromRuneReader(runeReader)
}

// ParseFromRuneReader parses a model from an io.RuneReader.
// It returns any error that occurs while parsing the input.
func ParseFromRuneReader(reader io.RuneReader) (*Model, error) {
	p := newParser(newLexer(reader))
	return p.parse()
}

// ----------------------------------------------

type tokenError struct {
	tok *token
	err error
}

func (e tokenError) Error() string {
	return fmt.Sprintf("line %d, pos %d: %s", e.tok.line, e.tok.pos, e.err)
}

/*
	TOKEN TYPES
	-----------
	Text
	Separator
	NewLine
	EOF


	PARSER STATES
	-------------

	[Start]
		Text / add new table
			---> [AfterName]
		Separator / check if we have table / create new row
			---> [AfterDataSeparator]
		NewLine
			---> [Start]
		EOF
			---> [End]

	[AfterName]
		NewLine
			---> [AfterNameLine]
		EOF
			---> [End]

	[AfterNameLine]
		Text / add new table
			---> [AfterName]
		Separator
			---> [AfterHeaderSeparator]
		NewLine
			---> [AfterNameLine]
		EOF
			---> [End]

	[AfterHeaderSeparator]
		Text / add new column to table
			---> [AfterHeaderText]
		EOF
			--> [End]

	[AfterHeaderText]
		Separator
			---> [AfterHeaderSeparator]
		NewLine
			---> [Start]
		EOF
			--> [End]

	[AfterDataSeparator]
		Text / add non-null value to last row, check type with i-th column
			---> [AfterDataText]
		Separator / add null value to last row, defer type from i-th column
			---> [AfterDataSeparator]
		NewLine / add null value to last row, defer type from i-th column / check last_row_width == header_width
			---> [Start]
		EOF / check last_row_width == header_width
			--> [End]

	[AfterDataText]
		Separator / check last_row_width < header width
			---> [AfterDataSeparator]
		NewLine / check last_row_width == header_width
			---> [Start]
		EOF / check last_row_width == header_width
			---> [End]

	[End]

*/

type parserState int

func (ps parserState) String() string {
	switch ps {
	case startState:
		return "start"
	case afterNameState:
		return "afterName"
	case afterNameLineState:
		return "afterNameLine"
	case afterHeaderSeparatorState:
		return "afterHeaderSeparator"
	case afterHeaderTextState:
		return "afterHeaderText"
	case afterDataSeparatorState:
		return "afterDataSeparator"
	case afterDataTextState:
		return "afterDataText"
	case endState:
		return "end"
	}
	panic("unknown parserState")
}

const (
	startState parserState = iota + 1
	afterNameState
	afterNameLineState
	afterHeaderSeparatorState
	afterHeaderTextState
	afterDataSeparatorState
	afterDataTextState
	endState
)

type parser struct {
	lex    *lexer
	state  parserState
	tables []*Table
	table  *Table
}

func newParser(lex *lexer) *parser {
	return &parser{lex, startState, []*Table{}, nil}
}

func (p *parser) parse() (*Model, error) {
	for {
		tok, err := p.lex.next()
		if err != nil {
			return nil, err
		}
		//fmt.Printf("%-20s %s\n", p.state, tok)
		switch p.state {
		case startState:
			err = p.forStart(tok)
		case afterNameState:
			err = p.forAfterName(tok)
		case afterNameLineState:
			err = p.forAfterNameLine(tok)
		case afterHeaderSeparatorState:
			err = p.forAfterHeaderSeparator(tok)
		case afterHeaderTextState:
			err = p.forAfterHeaderText(tok)
		case afterDataSeparatorState:
			err = p.forAfterDataSeparator(tok)
		case afterDataTextState:
			err = p.forAfterDataText(tok)
		case endState:
			panic("cannot parse beyond end")
		default:
			panic("invalid parser state")
		}
		if err != nil {
			return nil, tokenError{tok, err}
		}
		if p.state == endState {
			return &Model{p.tables}, nil
		}
	}
}

func (p *parser) forStart(tok *token) error {
	/*
		[Start]
			Text / add new table
				---> [AfterName]
			Separator / check if we have table / create new row
				---> [AfterDataSeparator]
			NewLine
				---> [Start]
			EOF
				---> [End]
	*/
	switch tok.ttype {
	case textToken:
		columns := make([]*Column, 0, 10)
		rows := make([]*Row, 0, 1000)
		p.table = &Table{tok.text, columns, rows}
		p.tables = append(p.tables, p.table)
		p.state = afterNameState
		return nil
	case separatorToken:
		if p.table == nil {
			return fmt.Errorf("unexpected separator")
		}
		values := make([]*Value, 0, 20)
		row := &Row{values}
		p.table.Rows = append(p.table.Rows, row)
		p.state = afterDataSeparatorState
		return nil
	case newlineToken:
		p.state = startState
		return nil
	case eofToken:
		p.state = endState
		return nil
	default:
		panic("invalid token type")
	}
}

func (p *parser) forAfterName(tok *token) error {
	/*
		[AfterName]
			NewLine
				---> [AfterNameLine]
			EOF
				---> [End]
	*/
	switch tok.ttype {
	case textToken:
		return fmt.Errorf("unexpected text")
	case separatorToken:
		return fmt.Errorf("unexpected separator")
	case newlineToken:
		p.state = afterNameLineState
		return nil
	case eofToken:
		p.state = endState
		return nil
	default:
		panic("invalid token type")
	}
}

func (p *parser) forAfterNameLine(tok *token) error {
	/*
		[AfterNameLine]
			Text / add new table
				---> [AfterName]
			Separator
				---> [AfterHeaderSeparator]
			NewLine
				---> [AfterNameLine]
			EOF
				---> [End]
	*/
	switch tok.ttype {
	case textToken:
		p.table = &Table{tok.text, []*Column{}, []*Row{}}
		p.tables = append(p.tables, p.table)
		p.state = afterNameState
		return nil
	case separatorToken:
		p.state = afterHeaderSeparatorState
		return nil
	case newlineToken:
		p.state = afterNameLineState
		return nil
	case eofToken:
		p.state = endState
		return nil
	default:
		panic("invalid token type")
	}
}

func (p *parser) forAfterHeaderSeparator(tok *token) error {
	/*
		[AfterHeaderSeparator]
			Text / add new column to table
				---> [AfterHeaderText]
			EOF
				--> [End]
	*/
	switch tok.ttype {
	case textToken:
		column, err := p.parseColumn(tok.text)
		if err != nil {
			return err
		}
		p.table.Columns = append(p.table.Columns, column)
		p.state = afterHeaderTextState
		return nil
	case separatorToken:
		return fmt.Errorf("unexpected separator")
	case newlineToken:
		return fmt.Errorf("unexpected end of line")
	case eofToken:
		p.state = endState
		return nil
	default:
		panic("invalid token type")
	}
}

func (p *parser) forAfterHeaderText(tok *token) error {
	/*
		[AfterHeaderText]
			Separator
				---> [AfterHeaderSeparator]
			NewLine
				---> [Start]
			EOF
				--> [End]
	*/
	switch tok.ttype {
	case textToken:
		return fmt.Errorf("unexpected text")
	case separatorToken:
		p.state = afterHeaderSeparatorState
		return nil
	case newlineToken:
		p.state = startState
		return nil
	case eofToken:
		p.state = endState
		return nil
	default:
		panic("invalid token type")
	}
}

func (p *parser) forAfterDataSeparator(tok *token) error {
	/*
		[AfterDataSeparator]
			Text / add non-null value to last row, check type with i-th column
				---> [AfterDataText]
			Separator / add null value to last row, defer type from i-th column
				---> [AfterDataSeparator]
			NewLine / add null value to last row, defer type from i-th column / check last_row_width == header_width
				---> [Start]
			EOF / check last_row_width == header_width
				--> [End]
	*/
	row := p.table.Rows[len(p.table.Rows)-1]
	columns := p.table.Columns
	switch tok.ttype {
	case textToken:
		// find type of i-th column
		colIndex := len(row.Values)
		if colIndex > len(columns)-1 {
			return fmt.Errorf("too many data values")
		}
		colType := columns[colIndex].Type
		// parse value
		value, err := p.parseValue(colType, tok.text)
		if err != nil {
			return err
		}
		// append value to last row
		row.Values = append(row.Values, value)
		p.state = afterDataTextState
		return nil
	case separatorToken:
		// find type of i-th column
		colIndex := len(row.Values)
		if colIndex > len(columns)-1 {
			return fmt.Errorf("too many data values")
		}
		colType := columns[colIndex].Type
		// append null value to last row
		value := &Value{Type: colType, Null: true}
		row.Values = append(row.Values, value)
		p.state = afterDataSeparatorState
		return nil
	case newlineToken:
		// find type of i-th column
		colIndex := len(row.Values)
		if colIndex > len(columns)-1 {
			return fmt.Errorf("too many data values")
		}
		colType := columns[colIndex].Type
		// append null value to last row
		value := &Value{Type: colType, Null: true}
		row.Values = append(row.Values, value)
		// check row_width == header_width
		if len(row.Values) < len(columns) {
			return fmt.Errorf("too few data values")
		}
		if len(row.Values) > len(columns) {
			return fmt.Errorf("too many data values")
		}
		p.state = startState
		return nil
	case eofToken:
		// check row_width == header_width
		if len(row.Values) < len(columns) {
			return fmt.Errorf("too few data values")
		}
		if len(row.Values) > len(columns) {
			return fmt.Errorf("too many data values")
		}
		p.state = endState
		return nil
	default:
		panic("invalid token type")
	}
}

func (p *parser) forAfterDataText(tok *token) error {
	/*
		[AfterDataText]
			Separator / check last_row_width < header width
				---> [AfterDataSeparator]
			NewLine / check last_row_width == header_width
				---> [Start]
			EOF / check last_row_width == header_width
				---> [End]
	*/
	row := p.table.Rows[len(p.table.Rows)-1]
	columns := p.table.Columns
	rowWidth := len(row.Values)
	headerWidth := len(columns)
	switch tok.ttype {
	case textToken:
		return fmt.Errorf("unexpected text")
	case separatorToken:
		if rowWidth >= headerWidth {
			return fmt.Errorf("too many data values")
		}
		p.state = afterDataSeparatorState
		return nil
	case newlineToken:
		if rowWidth > headerWidth {
			return fmt.Errorf("too many data values")
		}
		if rowWidth < headerWidth {
			return fmt.Errorf("too few data values")
		}
		p.state = startState
		return nil
	case eofToken:
		if rowWidth > headerWidth {
			return fmt.Errorf("too many data values")
		}
		if rowWidth < headerWidth {
			return fmt.Errorf("too few data values")
		}
		p.state = endState
		return nil
	default:
		panic("invalid token type")
	}
}

func (p *parser) parseColumn(text string) (*Column, error) {
	n := len(text)
	if n < 3 {
		return nil, fmt.Errorf("invalid column definition")
	}
	if text[n-2] != ':' {
		return nil, fmt.Errorf("invalid column definition")
	}
	typeChar := text[n-1]
	name := text[:n-2]
	switch typeChar {
	case 'i', 'f', 'b', 's', 't':
		return &Column{name, ValueType(typeChar)}, nil
	default:
		return nil, fmt.Errorf("invalid column type")
	}
}

func (p *parser) parseValue(colType ValueType, text string) (*Value, error) {
	v := &Value{Type: colType, Null: false}
	switch colType {
	case IntValue:
		x, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse as int: %s", err)
		}
		v.AsInt = x
	case FloatValue:
		x, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse as float: %s", err)
		}
		v.AsFloat = x
	case BoolValue:
		x, err := strconv.ParseBool(text)
		if err != nil {
			return nil, fmt.Errorf("cannot parse as bool: %s", err)
		}
		v.AsBool = x
	case StringValue:
		v.AsString = text
	case TimeValue:
		x, err := time.Parse("2006-01-02T15:04:05.999", text)
		if err != nil {
			return nil, fmt.Errorf("cannot parse as time: %s", err)
		}
		v.AsTime = x
	default:
		panic("wrong column type")
	}
	return v, nil
}
