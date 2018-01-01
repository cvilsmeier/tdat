package main

import (
	"bytes"
	"github.com/cvilsmeier/tdat/assert"
	"testing"
)

func TestConvertEmptyModelToJSON(t *testing.T) {
	txt := ""
	in := bytes.NewBufferString(txt)
	// to json
	out := &bytes.Buffer{}
	err := convertToJSON(in, out, " ")
	assert.Truef(t, err == nil, "err was %s", err)
	act := string(out.Bytes())
	exp := "{}\n"
	assert.EqStr(t, exp, act)
}

func TestConvertEmptyTablesToJSON(t *testing.T) {
	txt := "authors\n" +
		"\n" +
		"postings\n" +
		"\n"
	in := bytes.NewBufferString(txt)
	// to json
	out := &bytes.Buffer{}
	err := convertToJSON(in, out, " ")
	assert.Truef(t, err == nil, "err was %s", err)
	act := string(out.Bytes())
	exp := "{\n" +
		" \"authors\": [],\n" +
		" \"postings\": []\n" +
		"}\n"
	assert.EqStr(t, exp, act)
}

func TestConvertFullToJSON(t *testing.T) {
	txt := "authors\n" +
		"|id:i  |name:s            |registered:t         |rating:f\n" +
		"|1     |\"John Doe\"      |2017-12-12T10:00:00  |0.95\n" +
		"|2     |\"Mitch Kashmar\" |                     |\n" +
		"\n" +
		"postings\n" +
		"|id:i  |authorId:i            |date:t               |title:s\n" +
		"|1     |1                     |2017-12-12T10:00:00  |\"About something\"\n" +
		"|2     |1                     |2017-12-12T10:00:00  |\"About something else\"\n" +
		"\n"
	in := bytes.NewBufferString(txt)
	// to json
	out := &bytes.Buffer{}
	err := convertToJSON(in, out, " ")
	assert.Truef(t, err == nil, "err was %s", err)
	act := string(out.Bytes())
	exp := "{\n" +
		" \"authors\": [\n" +
		"  {\n" +
		"   \"id\": 1,\n" +
		"   \"name\": \"John Doe\",\n" +
		"   \"rating\": 0.95,\n" +
		"   \"registered\": \"2017-12-12T10:00:00Z\"\n" +
		"  },\n" +
		"  {\n" +
		"   \"id\": 2,\n" +
		"   \"name\": \"Mitch Kashmar\",\n" +
		"   \"rating\": null,\n" +
		"   \"registered\": null\n" +
		"  }\n" +
		" ],\n" +
		" \"postings\": [\n" +
		"  {\n" +
		"   \"authorId\": 1,\n" +
		"   \"date\": \"2017-12-12T10:00:00Z\",\n" +
		"   \"id\": 1,\n" +
		"   \"title\": \"About something\"\n" +
		"  },\n" +
		"  {\n" +
		"   \"authorId\": 1,\n" +
		"   \"date\": \"2017-12-12T10:00:00Z\",\n" +
		"   \"id\": 2,\n" +
		"   \"title\": \"About something else\"\n" +
		"  }\n" +
		" ]\n" +
		"}\n"
	assert.EqStr(t, exp, act)
}

func TestConvertEmptyModelToCSV(t *testing.T) {
	txt := ""
	in := bytes.NewBufferString(txt)
	// to csv
	out := &bytes.Buffer{}
	err := convertToCSV(in, out)
	assert.Truef(t, err == nil, "err was %s", err)
	act := string(out.Bytes())
	exp := ""
	assert.EqStr(t, exp, act)
}

func TestConvertEmptyTableToCSV(t *testing.T) {
	txt := "\n\n\n\nauthors\n\n\n\n\n\n\n\n"
	in := bytes.NewBufferString(txt)
	// to csv
	out := &bytes.Buffer{}
	err := convertToCSV(in, out)
	assert.Truef(t, err == nil, "err was %s", err)
	act := string(out.Bytes())
	exp := "authors\n\n"
	assert.EqStr(t, exp, act)
}

func TestConvertFullModelToCSV(t *testing.T) {
	txt := "authors\n" +
		"|id:i  |name:s                     |registered:t         |rating:f\n" +
		"|1     |\"John \\\"J.D.\\\" Doe\"  |2017-12-12T10:00:00  |0.95\n" +
		"|2     |\"Mitch Kashmar\"          |                     |\n" +
		"\n" +
		"postings\n" +
		"|id:i  |authorId:i            |date:t               |title:s\n" +
		"|1     |1                     |2017-12-12T10:00:00  |\"About something\"\n" +
		"|2     |1                     |2017-12-12T10:00:00  |\"About something else\"\n" +
		"|3     |2                     |2017-12-12T10:00:00  |\"儿童游戏\\\";;;\\\"\\r\\nfoo\"\n" +
		"\n"
	in := bytes.NewBufferString(txt)
	// to csv
	out := &bytes.Buffer{}
	err := convertToCSV(in, out)
	assert.Truef(t, err == nil, "err was %s", err)
	act := string(out.Bytes())
	exp := "" +
		"authors\n" +
		"id;name;registered;rating\n" +
		"1;\"John \"\"J.D.\"\" Doe\";2017-12-12 10:00:00;0.950000\n" +
		"2;Mitch Kashmar;;\n" +
		"\n" +
		"postings\n" +
		"id;authorId;date;title\n" +
		"1;1;2017-12-12 10:00:00;About something\n" +
		"2;1;2017-12-12 10:00:00;About something else\n" +
		"3;2;2017-12-12 10:00:00;\"儿童游戏\"\";;;\"\"\r\nfoo\"\n" +
		"\n"
	assert.EqStr(t, exp, act)
}
