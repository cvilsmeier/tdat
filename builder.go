package tdat

import (
	"time"
)

// Builder can be used to build models. It is easier to use Builder than to
// build tables 'by hand'.
type Builder struct {
	tableBuilders []*TableBuilder
}

// NewBuilder creates a new Builder.
// Initially, the Builder is empty (contains no tables).
func NewBuilder() *Builder {
	return &Builder{}
}

// AddTable adds a new table to the builder. It returns a new TableBuilder
// that can be used to add columns and rows to the table.
func (b *Builder) AddTable(name string) *TableBuilder {
	tb := newTableBuilder(name)
	b.tableBuilders = append(b.tableBuilders, tb)
	return tb
}

// Build builds and validates the model. If validation fails,
// a non-nil error is returned.
func (b *Builder) Build() (*Model, error) {
	tables := []*Table{}
	for _, tb := range b.tableBuilders {
		tables = append(tables, tb.build())
	}
	model := &Model{tables}
	err := ValidateModel(model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// MustBuild is like Build, except it panics if
// validation fails.
func (b *Builder) MustBuild() *Model {
	model, err := b.Build()
	if err != nil {
		panic(err)
	}
	return model
}

// ----------------------------------------------------

// TableBuilder is used to build Tables.
type TableBuilder struct {
	name        string
	columns     []*Column
	rowBuilders []*RowBuilder
}

func newTableBuilder(name string) *TableBuilder {
	return &TableBuilder{name: name}
}

// AddColumn adds a new column to the table.
func (b *TableBuilder) AddColumn(name string, columnType ValueType) {
	b.columns = append(b.columns, &Column{name, columnType})
}

// AddIntColumn adds a new IntValue column to the table.
func (b *TableBuilder) AddIntColumn(name string) { b.AddColumn(name, IntValue) }

// AddFloatColumn adds a new FloatValue column to the table.
func (b *TableBuilder) AddFloatColumn(name string) { b.AddColumn(name, FloatValue) }

// AddBoolColumn adds a new BoolValue column to the table.
func (b *TableBuilder) AddBoolColumn(name string) { b.AddColumn(name, BoolValue) }

// AddStringColumn adds a new StringValue column to the table.
func (b *TableBuilder) AddStringColumn(name string) { b.AddColumn(name, StringValue) }

// AddTimeColumn adds a new TimeValue column to the table.
func (b *TableBuilder) AddTimeColumn(name string) { b.AddColumn(name, TimeValue) }

// AddRow adds a new row to the table. It returns a RowBuilder that can be used
// to add values to the new Row.
func (b *TableBuilder) AddRow() *RowBuilder {
	rb := newRowBuilder()
	b.rowBuilders = append(b.rowBuilders, rb)
	return rb
}

func (b *TableBuilder) build() *Table {
	rows := []*Row{}
	for _, rb := range b.rowBuilders {
		rows = append(rows, rb.build())
	}
	return &Table{b.name, b.columns, rows}
}

// ----------------------------------------------------

// RowBuilder can be used to build Rows.
type RowBuilder struct {
	values []*Value
}

func newRowBuilder() *RowBuilder {
	return &RowBuilder{}
}

// AddValue adds a value of a specific valueType to the Row.
// If val is nil, a null Value is added to the row.
// Otherwise, the val argument must fit the valueType:
//
//    for IntValue:
//        val must be of type int64
//    for FloatValue:
//        val must be of type float64
//    for BoolValue:
//        val must be of type bool
//    for StringValue:
//        val must be of type string
//    for TimeValue:
//        val must be of type time.Time
//
// If the type of val does not fit the valueType properly, AddValue will panic.
func (b *RowBuilder) AddValue(valueType ValueType, val interface{}) {
	value := &Value{Type: valueType}
	if val == nil {
		value.Null = true
	} else {
		switch valueType {
		case IntValue:
			value.AsInt = val.(int64)
		case FloatValue:
			value.AsFloat = val.(float64)
		case BoolValue:
			value.AsBool = val.(bool)
		case StringValue:
			value.AsString = val.(string)
		case TimeValue:
			value.AsTime = val.(time.Time)
		default:
			panic("unknown valueType")
		}
	}
	b.values = append(b.values, value)
}

// AddIntValue adds a Value of type IntValue.
// The val parameter must be nil or of type int64.
func (b *RowBuilder) AddIntValue(val interface{}) { b.AddValue(IntValue, val) }

// AddFloatValue adds a Value of type FloatValue.
// The val parameter must be nil or of type float64.
func (b *RowBuilder) AddFloatValue(val interface{}) { b.AddValue(FloatValue, val) }

// AddBoolValue adds a Value of type BoolValue.
// The val parameter must be nil or of type bool.
func (b *RowBuilder) AddBoolValue(val interface{}) { b.AddValue(BoolValue, val) }

// AddStringValue adds a Value of type StringValue.
// The val parameter must be nil or of type string.
func (b *RowBuilder) AddStringValue(val interface{}) { b.AddValue(StringValue, val) }

// AddTimeValue adds a Value of type TimeValue.
// The val parameter must be nil or of type time.Time.
func (b *RowBuilder) AddTimeValue(val interface{}) { b.AddValue(TimeValue, val) }

func (b *RowBuilder) build() *Row {
	return &Row{b.values}
}
