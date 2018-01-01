package tdat

import (
	"encoding/json"
	"github.com/cvilsmeier/tdat/assert"
	"io/ioutil"
	"testing"
	"time"
)

func TestRenderToString(t *testing.T) {
	locBerlin, _ := time.LoadLocation("Europe/Berlin")
	model := &Model{
		[]*Table{
			{
				"persons",
				[]*Column{
					{"id", IntValue},
					{"size", FloatValue},
					{"flag", BoolValue},
					{"name", StringValue},
					{"birth", TimeValue},
				},
				[]*Row{
					{
						[]*Value{
							{Type: IntValue, AsInt: int64(1)},
							{Type: FloatValue, AsFloat: float64(1.83)},
							{Type: BoolValue, AsBool: true},
							{Type: StringValue, AsString: "Joe \u2602 Smith"},
							{Type: TimeValue, AsTime: time.Date(2001, 1, 2, 10, 11, 12, 13000000, locBerlin)},
						},
					},
					{
						[]*Value{
							{Type: IntValue, Null: true},
							{Type: FloatValue, Null: true},
							{Type: BoolValue, Null: true},
							{Type: StringValue, Null: true},
							{Type: TimeValue, Null: true},
						},
					},
				},
			},
		},
	}
	//
	colwidth := 0
	txt, err := RenderToString(model, colwidth)
	assert.Truef(t, err == nil, "err=%s", err)
	exp := "persons\n"
	exp += "|id:i|size:f|flag:b|name:s|birth:t\n"
	exp += "|1|1.830000|true|\"Joe \u2602 Smith\"|2001-01-02T09:11:12.013\n"
	exp += "|||||\n"
	exp += "\n"
	assert.EqStr(t, exp, txt)
	//
	colwidth = 10
	txt, err = RenderToString(model, colwidth)
	assert.Truef(t, err == nil, "error was %s", err)
	exp = "persons\n"
	exp += "|id:i      |size:f    |flag:b    |name:s    |birth:t\n"
	exp += "|1         |1.830000  |true      |\"Joe \u2602 Smith\"|2001-01-02T09:11:12.013\n"
	exp += "|          |          |          |          |\n"
	exp += "\n"
	assert.EqStr(t, exp, txt)
}

func BenchmarkRenderTdat(b *testing.B) {
	rowCount := 100 * 1000
	rows := make([]*Row, 0, rowCount)
	for i := 0; i < rowCount; i++ {
		row := &Row{
			[]*Value{
				{Type: IntValue, AsInt: int64(1)},
				{Type: FloatValue, AsFloat: float64(13000.12)},
				{Type: BoolValue, AsBool: true},
				{Type: StringValue, AsString: "joe"},
			},
		}
		rows = append(rows, row)
	}
	model := &Model{
		[]*Table{
			{
				"persons",
				[]*Column{
					{"id", IntValue},
					{"rate", FloatValue},
					{"flag", BoolValue},
					{"name", StringValue},
				},
				rows,
			},
		},
	}
	b.ResetTimer()
	err := RenderToWriter(model, 0, ioutil.Discard)
	if err != nil {
		panic(err)
	}
}

func BenchmarkRenderJson(b *testing.B) {
	rowCount := 100 * 1000
	persons := []map[string]interface{}{}
	for i := 0; i < rowCount; i++ {
		person := map[string]interface{}{
			"id":   1,
			"rate": 13000.12,
			"flag": true,
			"name": "joe",
		}
		persons = append(persons, person)
	}
	model := map[string]interface{}{
		"persons": persons,
	}
	b.ResetTimer()
	encod := json.NewEncoder(ioutil.Discard)
	err := encod.Encode(model)
	if err != nil {
		panic(err)
	}
}
