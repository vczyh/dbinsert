package relational

import (
	"fmt"
	"strings"
)

type Dialect uint8

const (
	DialectMysql = iota
)

var DefaultTable = &Table{
	Database: "zyh_db9",
	Name:     "tbl",
	Size:     1_000_000,
	Fields: []*Field{
		MustNewField("id", "INT"),
		MustNewField("val", "CHAR(80)"),
	},
}

type Table struct {
	Dialect              Dialect
	Database             string
	Name                 string
	Fields               []*Field
	Size                 int
	PrimaryKeyFieldNames []string
}

func NewTables(schema *Schema) (tables []*Table, err error) {
	schemaTables := schema.Tables()

	for _, schemaTable := range schemaTables {
		table := &Table{
			Dialect:              schema.Dialect(),
			Database:             schemaTable.Database,
			Name:                 schemaTable.Table,
			Fields:               nil,
			Size:                 schemaTable.Size,
			PrimaryKeyFieldNames: schemaTable.PrimaryKeyFieldNames,
		}
		for _, schemaField := range schemaTable.Fields {
			field, err := NewField(schemaField.Name, schemaField.TypeName, schemaField.AutoIncrement)
			if err != nil {
				return nil, err
			}
			table.Fields = append(table.Fields, field)
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (t *Table) ColumnNames() []string {
	var names []string
	for _, column := range t.Fields {
		names = append(names, column.Name)
	}
	return names
}

func (t *Table) NotAutoIncrementColumnNames() []string {
	var names []string
	for _, column := range t.Fields {
		if !column.autoIncrement {
			names = append(names, column.Name)
		}
	}
	return names
}

func (t *Table) DDL() string {
	sb := new(strings.Builder)
	fmt.Fprintf(sb, "CREATE TABLE IF NOT EXISTS %s\n", t.Name)
	fmt.Fprintln(sb, "(")
	for i, column := range t.Fields {
		fmt.Fprintf(sb, "\t%s %s", column.Name, column.TypeName)
		if column.autoIncrement {
			sb.WriteString(" " + "AUTO_INCREMENT")
		}
		if i != len(t.Fields)-1 {
			fmt.Fprint(sb, ",")
		} else {
			if len(t.PrimaryKeyFieldNames) != 0 {
				fmt.Fprintln(sb, ",")
				fmt.Fprintf(sb, "\tPRIMARY KEY (%s)", strings.Join(t.PrimaryKeyFieldNames, ","))
			}
		}
		fmt.Fprintln(sb)
	}
	fmt.Fprintln(sb, ")")
	return sb.String()
}

//func (t *Table) DML() string {
//	sb := new(strings.Builder)
//
//	var names, placeholders []string
//	for _, column := range t.Columns {
//		names = append(names, column.Name)
//		placeholders = append(placeholders, "?")
//	}
//
//	fmt.Fprintf(sb, "INSERT INTO %s", t.Name)
//	fmt.Fprintf(sb, " ( %s )", strings.Join(names, ", "))
//	fmt.Fprintf(sb, " VALUES ( %s )", strings.Join(placeholders, ", "))
//
//	return sb.String()
//}
