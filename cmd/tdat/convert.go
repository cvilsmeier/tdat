package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/cvilsmeier/tdat"
	"io"
)

func convertToJSON(r io.Reader, w io.Writer, indent string) error {
	model, err := tdat.ParseFromReader(r)
	if err != nil {
		return err
	}
	err = tdat.ValidateModel(model)
	if err != nil {
		return err
	}
	jsModel := map[string]interface{}{}
	for _, table := range model.Tables {
		jsTable := []map[string]interface{}{}
		for _, row := range table.Rows {
			jsRow := map[string]interface{}{}
			for columnIndex, column := range table.Columns {
				value := row.Values[columnIndex]
				if value.Null {
					jsRow[column.Name] = nil
				} else {
					switch value.Type {
					case tdat.IntValue:
						jsRow[column.Name] = value.AsInt
					case tdat.FloatValue:
						jsRow[column.Name] = value.AsFloat
					case tdat.BoolValue:
						jsRow[column.Name] = value.AsBool
					case tdat.StringValue:
						jsRow[column.Name] = value.AsString
					case tdat.TimeValue:
						jsRow[column.Name] = value.AsTime
					default:
						panic("invalid value type")
					}
				}
			}
			jsTable = append(jsTable, jsRow)
		}
		jsModel[table.Name] = jsTable
	}
	// to json
	enc := json.NewEncoder(w)
	if indent != "" {
		enc.SetIndent("", indent)
	}
	err = enc.Encode(jsModel)
	return err
}

func convertToCSV(r io.Reader, w io.Writer) error {
	model, err := tdat.ParseFromReader(r)
	if err != nil {
		return err
	}
	err = tdat.ValidateModel(model)
	if err != nil {
		return err
	}
	csvWriter := csv.NewWriter(w)
	csvWriter.Comma = ';'
	for _, table := range model.Tables {
		// write table name
		{
			err := csvWriter.Write([]string{table.Name})
			if err != nil {
				return err
			}
		}
		if len(table.Columns) > 0 {
			// write columns
			{
				record := []string{}
				for _, column := range table.Columns {
					record = append(record, column.Name)
				}
				err := csvWriter.Write(record)
				if err != nil {
					return err
				}
			}
			// write rows
			{
				for _, row := range table.Rows {
					record := []string{}
					for columnIndex := range table.Columns {
						value := row.Values[columnIndex]
						cell := ""
						if !value.Null {
							switch value.Type {
							case tdat.IntValue:
								cell = fmt.Sprintf("%d", value.AsInt)
							case tdat.FloatValue:
								cell = fmt.Sprintf("%f", value.AsFloat)
							case tdat.BoolValue:
								cell = fmt.Sprintf("%t", value.AsBool)
							case tdat.StringValue:
								cell = value.AsString
							case tdat.TimeValue:
								cell = value.AsTime.Format("2006-01-02 15:04:05")
							default:
								panic("invalid value type")
							}
						}
						record = append(record, cell)
					}
					err := csvWriter.Write(record)
					if err != nil {
						return err
					}
				}
			}
		}
		// write empty line
		{
			err := csvWriter.Write([]string{})
			if err != nil {
				return err
			}
		}
	}
	csvWriter.Flush()
	return nil
}
