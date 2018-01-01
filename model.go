package tdat

import (
	"time"
)

// A Model contains zero or more tables.
type Model struct {
	// The tables of the model
	Tables []*Table
}

// A Table contains zero or more columns and zero or more rows.
type Table struct {
	// Name is the name of the table
	Name string
	// The columns of the table
	Columns []*Column
	// A slice of data rows
	Rows []*Row
}

// A Column has a name and a type.
type Column struct {
	Name string
	Type ValueType
}

// Row contains zero or more values.
type Row struct {
	// Each value is either a int, a float, a bool, etc.
	Values []*Value
}

// ValueType represents the type of a column or value.
type ValueType byte

const (
	
	// IntValue represents a int64 value. Its code is 'i'.
	IntValue    ValueType = 'i'

	// FloatValue represents a float64 value. Its code is 'f'.
	FloatValue            = 'f'

	// BoolValue represents a bool value. Its code is 'b'.
	BoolValue             = 'b'

	// StringValue represents a string value. Its code is 's'.
	StringValue           = 's'

	// TimeValue represents a time.Time value. Its code is 't'.
	TimeValue             = 't'
)

// IsValid returns true if t is a valid ValueType.
func (t ValueType) IsValid() bool {
	switch t {
	case 'i', 'f', 'b', 's', 't':
		return true
	}
	return false
}

// Value represents a value in a table row.
type Value struct {

	// The type of the Value (IntValue, FloatValue, etc.).
	Type ValueType

	// Null is true if this Value is null (undefined, nil, nothing).
	Null bool

	// Holds the value for IntValue.
	AsInt int64

	// Holds the value for FloatValue.
	AsFloat float64

	// Holds the value for BoolValue.
	AsBool bool

	// Holds the value for StringValue.
	AsString string

	// Holds the value for TimeValue.
	AsTime time.Time
}
