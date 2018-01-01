/*

The tdat tool provides utilities for handling TDAT files.

	$ tdat

*/
package main

import (
	"flag"
	"fmt"
	"github.com/cvilsmeier/tdat"
	"os"
	"time"
)

var cmdFlag = ""
var inFlag = "-"
var outFlag = "-"
var indentFlag = ""

func usage() {
	fmt.Fprintf(os.Stderr, "tdat - a tool for handling TDAT files\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "tdat -cmd validate [-in <filename>] [-out <filename>]\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "    Cmd validate parses and validates a tdat model. If the model\n")
	fmt.Fprintf(os.Stderr, "    is valid, tdat will print nothing and exit with code 0.\n")
	fmt.Fprintf(os.Stderr, "    If the model is not valid, tdat will print an error message to\n")
	fmt.Fprintf(os.Stderr, "    stderr and exit with code 1.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "tdat -cmd json [-in <filename>] [-out <filename>] [-indent <pattern>]\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "    Cmd json parses and validates a tdat model and convert it to\n")
	fmt.Fprintf(os.Stderr, "    JSON format. If indent is \"\" (the default), the JSON will be\n")
	fmt.Fprintf(os.Stderr, "    written as one line. If indent is not empty, the JSON will\n")
	fmt.Fprintf(os.Stderr, "    be multi-line, each line indented by the indent pattern.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "tdat -cmd csv [-in <filename>] [-out <filename>]\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "    Cmd csv parses and validates a tdat model and convert it to\n")
	fmt.Fprintf(os.Stderr, "    CSV format.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	flag.StringVar(&cmdFlag, "cmd", cmdFlag, "specifies the command to execute")
	flag.StringVar(&inFlag, "in", inFlag, "read from the specified file. '-' means stdin.")
	flag.StringVar(&outFlag, "out", outFlag, "write to the specified file. '-' means stdout.")
	flag.StringVar(&indentFlag, "indent", indentFlag, "indentation of json output")
	flag.Usage = usage
	flag.Parse()
	switch cmdFlag {
	case "sample":
		sample()
		os.Exit(0)
	case "validate":
		err := validate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "json", "csv":
		err := convert()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "", "help":
		usage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", cmdFlag)
		os.Exit(2)
	}
}

func sample() {
	builder := tdat.NewBuilder()
	table := builder.AddTable("products")
	table.AddIntColumn("id")
	table.AddFloatColumn("rating")
	table.AddBoolColumn("in_stock")
	table.AddStringColumn("name")
	table.AddTimeColumn("dateOfEntry")
	{
		row := table.AddRow()
		row.AddIntValue(int64(1))
		row.AddFloatValue(float64(112.13))
		row.AddBoolValue(true)
		row.AddStringValue("a book")
		row.AddTimeValue(time.Now().Add(-10000 * time.Hour))
	}
	model := builder.MustBuild()
	colwidth := 15
	txt, err := tdat.RenderToString(model, colwidth)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(txt)
}

func validate() error {
	r := os.Stdin
	if inFlag != "-" {
		f, err := os.Open(inFlag)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}
	model, err := tdat.ParseFromReader(r)
	if err != nil {
		return err
	}
	err = tdat.ValidateModel(model)
	if err != nil {
		return err
	}
	return nil
}

func convert() error {
	r := os.Stdin
	if inFlag != "-" {
		f, err := os.Open(inFlag)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}
	w := os.Stdout
	if outFlag != "-" {
		f, err := os.OpenFile(outFlag, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	// convert
	switch cmdFlag {
	case "json":
		err := convertToJSON(r, w, indentFlag)
		if err != nil {
			return err
		}
	case "csv":
		err := convertToCSV(r, w)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown format %q", cmdFlag)
	}
	return nil
}
