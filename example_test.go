package tdat_test

import (
	"fmt"
	"github.com/cvilsmeier/tdat"
	"log"
	"time"
)

func ExampleParseFromString() {
	input := `
products
|id:i  |name:s    |date:t
|1     |"bottle"  |2017-12-12T10:11:12.013
|2     |"book"    |2017-11-11T10:11:12.013
`
	model, err := tdat.ParseFromString(input)
	if err != nil {
		log.Fatal(err)
	}
	table := model.Tables[0]
	fmt.Printf("%s\n", table.Name)
	fmt.Printf("%s, %s, %s\n", table.Columns[0].Name, table.Columns[1].Name, table.Columns[2].Name)
	for _, row := range table.Rows {
		id := row.Values[0].AsInt
		name := row.Values[1].AsString
		date := row.Values[2].AsTime
		fmt.Printf("%d, %s, %s\n", id, name, date)
	}
	// Output:
	// products
	// id, name, date
	// 1, bottle, 2017-12-12 10:11:12.013 +0000 UTC
	// 2, book, 2017-11-11 10:11:12.013 +0000 UTC
}

func ExampleBuilder() {
	builder := tdat.NewBuilder()
	// add "products" table
	table := builder.AddTable("products")
	// add 3 columns
	table.AddIntColumn("id")
	table.AddStringColumn("name")
	table.AddTimeColumn("date")
	// add a row of values
	row := table.AddRow()
	row.AddIntValue(int64(1))
	row.AddStringValue("bottle")
	row.AddTimeValue(time.Date(2012, 1, 1, 10, 0, 0, 0, time.UTC))
	// add another row of values
	row = table.AddRow()
	row.AddValue(tdat.IntValue, int64(2))
	row.AddValue(tdat.StringValue, "unknown")
	row.AddTimeValue(nil)
	// build and validate model
	model, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}
	// render model to string
	txt, err := tdat.RenderToString(model, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(txt)
	// Output:
	// products
	// |id:i      |name:s    |date:t
	// |1         |"bottle"  |2012-01-01T10:00:00
	// |2         |"unknown" |
}

func ExampleRenderToString() {
	// build model
	builder := tdat.NewBuilder()
	table := builder.AddTable("products")
	table.AddIntColumn("id")
	table.AddStringColumn("name")
	row := table.AddRow()
	row.AddIntValue(int64(1))
	row.AddStringValue("bottle")
	model := builder.MustBuild()
	// render model to string
	txt, err := tdat.RenderToString(model, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(txt)
	// Output:
	// products
	// |id:i      |name:s
	// |1         |"bottle"
}
