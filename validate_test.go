package tdat

import (
	"github.com/cvilsmeier/tdat/assert"
	"testing"
)

func TestValidateModelNoValues(t *testing.T) {
	model := &Model{
		[]*Table{
			&Table{
				"products",
				[]*Column{
					&Column{"id", IntValue},
				},
				[]*Row{
					&Row{},
				},
			},
		},
	}
	err := ValidateModel(model)
	assert.True(t, err != nil)
	assert.EqStr(t, "table \"products\": row 1: expected 1 values but got 0", err.Error())
}

func TestValidateTooManyValues(t *testing.T) {
	model := &Model{
		[]*Table{
			&Table{
				"products",
				[]*Column{
					&Column{"id", IntValue},
					&Column{"name", StringValue},
				},
				[]*Row{
					&Row{
						[]*Value{
							&Value{Type: IntValue},
							&Value{Type: StringValue},
							&Value{Type: BoolValue},
						},
					},
				},
			},
		},
	}
	err := ValidateModel(model)
	assert.True(t, err != nil)
	assert.EqStr(t, "table \"products\": row 1: expected 2 values but got 3", err.Error())
}

func TestValidateWrongValueType(t *testing.T) {
	model := &Model{
		[]*Table{
			&Table{
				"products",
				[]*Column{
					&Column{"id", IntValue},
					&Column{"name", StringValue},
				},
				[]*Row{
					&Row{
						[]*Value{
							&Value{Type: IntValue},
							&Value{Type: BoolValue},
						},
					},
				},
			},
		},
	}
	err := ValidateModel(model)
	assert.True(t, err != nil)
	assert.EqStr(t, "table \"products\": row 1, value 2: expected value type 's' but was 'b'", err.Error())
}

func TestValidateOk(t *testing.T) {
	model := &Model{
		[]*Table{
			&Table{
				"products",
				[]*Column{
					&Column{"id", IntValue},
					&Column{"name", StringValue},
				},
				[]*Row{
					&Row{
						[]*Value{
							&Value{Type: IntValue},
							&Value{Type: StringValue},
						},
					},
				},
			},
		},
	}
	err := ValidateModel(model)
	assert.Truef(t, err == nil, "err was %s", err)
}
