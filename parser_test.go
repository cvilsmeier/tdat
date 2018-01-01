package tdat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cvilsmeier/tdat/assert"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	fis, err := ioutil.ReadDir("testdata")
	assert.Truef(t, err == nil, "err was %s", err)
	for _, fi := range fis {
		fname := fi.Name()
		if !strings.HasPrefix(fname, "parser_") {
			break
		}
		t.Run(fname, func(t *testing.T) {
			buf, err := ioutil.ReadFile("testdata/" + fname)
			assert.Truef(t, err == nil, "err was %s", err)
			all := string(buf)
			i := strings.Index(all, "---EOF---")
			assert.Truef(t, i >= 0, "i was %d", i)
			input := all[:i]
			exp := strings.TrimSpace(all[i+9:])
			// parse input
			model, err := ParseFromString(input)
			// stringify parse result
			act := ""
			if model != nil {
				for _, table := range model.Tables {
					act += stringifyTable(table)
				}
			}
			if err != nil {
				act += err.Error()
			}
			act = strings.TrimSpace(act)
			assert.EqStr(t, exp, act)
		})
	}
}

func stringifyTable(table *Table) string {
	str := fmt.Sprintf("table %q\n", table.Name)
	for _, col := range table.Columns {
		str += fmt.Sprintf("  col %q(%c)\n", col.Name, col.Type)
	}
	for rowIndex, row := range table.Rows {
		str += fmt.Sprintf("row %d\n", rowIndex+1)
		for _, val := range row.Values {
			if val.Null {
				str += fmt.Sprintf("  val null(%c)\n", val.Type)
			} else {
				valText := ""
				switch val.Type {
				case IntValue:
					valText = fmt.Sprintf("%d", val.AsInt)
				case FloatValue:
					valText = fmt.Sprintf("%f", val.AsFloat)
				case BoolValue:
					valText = fmt.Sprintf("%t", val.AsBool)
				case StringValue:
					valText = fmt.Sprintf("%s", val.AsString)
				case TimeValue:
					valText = fmt.Sprintf("%s", val.AsTime)
				default:
					panic("wrong type")
				}
				str += fmt.Sprintf("  val %s(%c)\n", valText, val.Type)
			}
		}
	}
	return str
}

func BenchmarkAllocateModel(b *testing.B) {
	b.Skip()
	rowCount := 200 * 1000
	rows := make([]*Row, 0)
	for i := 0; i < rowCount; i++ {
		id, _ := strconv.ParseInt("1", 10, 64)
		idValue := &Value{Type: IntValue, AsInt: id}
		rate, _ := strconv.ParseFloat("13000.12", 64)
		rateValue := &Value{Type: FloatValue, AsFloat: rate}
		flag, _ := strconv.ParseBool("true")
		flagValue := &Value{Type: BoolValue, AsBool: flag}
		nameValue := &Value{Type: StringValue, AsString: "joe"}
		row := &Row{[]*Value{idValue, rateValue, flagValue, nameValue}}
		rows = append(rows, row)
	}
}

func BenchmarkParseValue(b *testing.B) {
	b.Skip()
	p := &parser{}
	rowCount := 100 * 1000
	for i := 0; i < rowCount; i++ {
		p.parseValue(IntValue, "13")
		p.parseValue(FloatValue, "1300.13")
		p.parseValue(BoolValue, "true")
		p.parseValue(StringValue, "lorem ipsum")
		p.parseValue(TimeValue, "2017-12-12T10:00:00.113")
	}
}

func BenchmarkParseTdat(t *testing.B) {
	rowCount := 100 * 1000
	var input string
	{
		b := &bytes.Buffer{}
		b.WriteString("persons\n")
		b.WriteString("|id:i|rate:f|flag:b|name:s\n")
		for i := 0; i < rowCount; i++ {
			b.WriteString("|1|13000.00|true|\"joe\"\n")
		}
		input = string(b.Bytes())
	}
	t.ResetTimer()
	m, err := ParseFromString(input)
	if err != nil {
		panic(err)
	}
	if m == nil {
		panic("no model")
	}
}

func BenchmarkParseJson(t *testing.B) {
	rowCount := 100 * 1000
	var input []byte
	{
		b := &bytes.Buffer{}
		b.WriteString("{\"persons\":[")
		for i := 0; i < rowCount; i++ {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString("{\"id\":1,\"rate\":13000.00,\"flag\":true,\"name\":\"joe\"}")
		}
		b.WriteString("]}")
		input = b.Bytes()
	}
	t.ResetTimer()
	model := map[string]interface{}{}
	t.ResetTimer()
	err := json.Unmarshal(input, &model)
	if err != nil {
		panic(err)
	}
}
