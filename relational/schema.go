package relational

import (
	"encoding/json"
	"fmt"
	"os"
)

type Schema struct {
	dialect Dialect
	dbs     map[string][]*SchemaTable
	tables  []*SchemaTable
}

func ParseSchemaFromFile(dialect Dialect, file string, dbRepeat int) (*Schema, error) {
	tables, err := readSchemaTablesFromFile(file)
	if err != nil {
		return nil, err
	}
	return NewSchema(dialect, tables, dbRepeat)
}

func NewSchema(dialect Dialect, tables []*SchemaTable, dbRepeat int) (*Schema, error) {
	s := &Schema{
		dialect: dialect,
	}

	s.dbs = make(map[string][]*SchemaTable)
	for _, table := range tables {
		s.dbs[table.Database] = append(s.dbs[table.Database], table)
	}

	if dbRepeat > 0 {
		newDbs := make(map[string][]*SchemaTable)
		for db, tables := range s.dbs {
			newDbs[db] = tables
			for i := 0; i < dbRepeat; i++ {
				newDb := fmt.Sprintf("%s_%d", db, i+1)
				var newTables []*SchemaTable
				for _, table := range tables {
					newTable := table.CLone()
					newTable.Database = newDb
					newTables = append(newTables, newTable)
				}
				newDbs[newDb] = newTables
			}
		}
		s.dbs = newDbs
	}

	return s, nil
}

func (s *Schema) Dialect() Dialect {
	return s.dialect
}

func (s *Schema) Tables() []*SchemaTable {
	if s.tables == nil {
		for _, tables := range s.dbs {
			s.tables = append(s.tables, tables...)
		}
	}
	return s.tables
}

func (s *Schema) Databases() map[string][]*SchemaTable {
	return s.dbs
}

type SchemaTable struct {
	Database             string         `json:"database"`
	Table                string         `json:"table"`
	Size                 int            `json:"size"`
	Fields               []*SchemaField `json:"fields"`
	PrimaryKeyFieldNames []string       `json:"primaryKeyFieldNames"`
}

func (st *SchemaTable) CLone() *SchemaTable {
	c := &SchemaTable{
		Database:             st.Database,
		Table:                st.Table,
		Size:                 st.Size,
		Fields:               nil,
		PrimaryKeyFieldNames: nil,
	}
	for _, field := range st.Fields {
		c.Fields = append(c.Fields, field.Clone())
	}
	for _, name := range st.PrimaryKeyFieldNames {
		c.PrimaryKeyFieldNames = append(c.PrimaryKeyFieldNames, name)
	}
	return c
}

type SchemaField struct {
	Name          string `json:"name"`
	TypeName      string `json:"type"`
	AutoIncrement bool   `json:"autoIncrement"`
}

func (sf *SchemaField) Clone() *SchemaField {
	return &SchemaField{
		Name:          sf.Name,
		TypeName:      sf.TypeName,
		AutoIncrement: sf.AutoIncrement,
	}
}

func readSchemaTablesFromFile(file string) (schemas []*SchemaTable, err error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return schemas, json.Unmarshal(b, &schemas)
}
