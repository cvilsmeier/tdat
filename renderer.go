package tdat

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// RenderToString is like RenderToWriter but renders to a string.
func RenderToString(model *Model, colWidth int) (string, error) {
	buffer := &bytes.Buffer{}
	err := RenderToWriter(model, colWidth, buffer)
	if err != nil {
		return "", err
	}
	txt := string(buffer.Bytes())
	return txt, nil
}

// RenderToFile is like RenderToWriter but renders into a file.
// If the file exists, it is overwritten.
func RenderToFile(model *Model, colWidth int, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return RenderToWriter(model, colWidth, file)
}

// RenderToWriter renders a model to a io.Writer.
// The renderer pads columns with spaces, so that each column
// has at least colWidth characters.
// If colWidth <= 0, no padding is applied.
func RenderToWriter(model *Model, colWidth int, w io.Writer) error {
	r := renderer{w: w, colWidth: colWidth}
	r.renderModel(model)
	return r.err
}

// ------------------------------------------------------------

type renderer struct {
	w        io.Writer
	colWidth int
	err      error
}

func (r renderer) renderModel(model *Model) {
	for _, table := range model.Tables {
		r.renderTable(table)
	}
}

func (r renderer) renderTable(table *Table) {
	if r.err != nil {
		return
	}
	// table name
	r.printf("%s\n", table.Name)
	// columns
	r.renderColumns(table.Columns)
	// rows
	for _, row := range table.Rows {
		if r.err != nil {
			return
		}
		r.renderRow(row)
	}
	r.printf("\n")
}

func (r renderer) renderColumns(columns []*Column) {
	colCount := len(columns)
	for colIndex, col := range columns {
		cell := fmt.Sprintf("%s:%c", col.Name, col.Type)
		if r.colWidth <= 0 || colIndex >= colCount-1 {
			r.printf("|%s", cell)
		} else {
			r.printf("|%-*s", r.colWidth, cell)
		}
	}
	if colCount > 0 {
		r.printf("\n")
	}
}

func (r renderer) renderRow(row *Row) {
	valCount := len(row.Values)
	for valIndex, val := range row.Values {
		cell := ""
		if !val.Null {
			switch val.Type {
			case IntValue:
				cell = fmt.Sprintf("%d", val.AsInt)
			case FloatValue:
				cell = fmt.Sprintf("%f", val.AsFloat)
			case BoolValue:
				cell = fmt.Sprintf("%t", val.AsBool)
			case StringValue:
				cell = fmt.Sprintf("%q", val.AsString)
			case TimeValue:
				cell = val.AsTime.UTC().Format("2006-01-02T15:04:05.999")
			default:
				panic("wrong value type")
			}
		}
		if r.colWidth <= 0 || valIndex >= valCount-1 {
			r.printf("|%s", cell)
		} else {
			r.printf("|%-*s", r.colWidth, cell)
		}
	}
	if valCount > 0 {
		r.printf("\n")
	}
}

func (r renderer) printf(format string, args ...interface{}) {
	if r.err != nil {
		return
	}
	s := fmt.Sprintf(format, args...)
	_, err := r.w.Write([]byte(s))
	if err != nil {
		r.err = err
	}
}
