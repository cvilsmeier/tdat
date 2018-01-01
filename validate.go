package tdat

import (
	"fmt"
	"strings"
)

// ValidateModel validates a model. If the model is invalid, it returns a
// non-nil error.
func ValidateModel(model *Model) error {
	tableNames := map[string]bool{}
	for _, table := range model.Tables {
		// validate table name
		err := ValidateName(table.Name)
		if err != nil {
			return fmt.Errorf("table %q: %s", table.Name, err)
		}
		if tableNames[table.Name] {
			return fmt.Errorf("duplicate table %q", table.Name)
		}
		tableNames[table.Name] = true
		// validate table
		err = ValidateTable(table)
		if err != nil {
			return fmt.Errorf("table %q: %s", table.Name, err)
		}
	}
	return nil
}

// ValidateTable validates a table. If the table is invalid, it returns a
// non-nil error.
func ValidateTable(table *Table) error {
	// validate columns
	columnNames := map[string]bool{}
	for _, column := range table.Columns {
		// validate name
		err := ValidateName(column.Name)
		if err != nil {
			return fmt.Errorf("column %q: %s", column.Name, err)
		}
		if columnNames[column.Name] {
			return fmt.Errorf("duplicate column %q", column.Name)
		}
		columnNames[column.Name] = true
		// validate type
		ct := column.Type
		if !column.Type.IsValid() {
			return fmt.Errorf("column %q has invalid type '%c'", column.Name, ct)
		}
	}
	colCount := len(table.Columns)
	// validate rows
	for rowIndex, row := range table.Rows {
		valCount := len(row.Values)
		if valCount != colCount {
			return fmt.Errorf("row %d: expected %d values but got %d", rowIndex+1, colCount, valCount)
		}
		for valueIndex, value := range row.Values {
			column := table.Columns[valueIndex]
			if value.Type != column.Type {
				return fmt.Errorf("row %d, value %d: expected value type '%c' but was '%c'", rowIndex+1, valueIndex+1, column.Type, value.Type)
			}
		}
	}
	return nil
}

// ValidateName validates a table or column name.
// If the name is not valid, it returns a non-nil error.
func ValidateName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("name is empty")
	}
	if name != strings.TrimSpace(name) {
		return fmt.Errorf("name contains whitespace")
	}
	for _, r := range name {
		if r <= ' ' {
			return fmt.Errorf("name contains invalid character '%c'", r)
		}
	}
	return nil
}
