package relation

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
)

var (
	//go:embed schema
	schemaF embed.FS

	schemaFiles = map[Dialect]map[string]string{
		DialectMysql:    {"sysbench": "schema/sysbench_mysql.json"},
		DialectPostgres: {"sysbench": "schema/sysbench_postgres.json"},
	}
)

type Schema struct {
	initDbRepeat  int
	initTableSize int

	dialect Dialect
	dbs     map[string][]*SchemaTable
	tables  []*SchemaTable
}

func ParseSchemaFromFile(dialect Dialect, file string, opts ...SchemaOption) (*Schema, error) {
	tables, err := readSchemaTablesFromFileOrTemplate(dialect, file)
	if err != nil {
		return nil, err
	}
	return NewSchema(dialect, tables, opts...)
}

func NewSchema(dialect Dialect, tables []*SchemaTable, opts ...SchemaOption) (*Schema, error) {
	s := &Schema{
		dialect:       dialect,
		initDbRepeat:  0,
		initTableSize: 0,
	}
	for _, opt := range opts {
		if err := opt.apply(s); err != nil {
			return nil, err
		}
	}

	s.dbs = make(map[string][]*SchemaTable)
	for _, table := range tables {
		if s.initTableSize != 0 {
			table.Size = s.initTableSize
		}
		s.dbs[table.Database] = append(s.dbs[table.Database], table)
	}

	if s.initDbRepeat > 0 {
		newDbs := make(map[string][]*SchemaTable)
		for db, tables := range s.dbs {
			//newDbs[db] = tables
			for i := 0; i < s.initDbRepeat+1; i++ {
				newDb := fmt.Sprintf("%s_%d", db, i+1)
				var newTables []*SchemaTable
				for _, table := range tables {
					newTable := table.Clone()
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

func WithSchemaDatabasesRepeat(dbRepeat int) SchemaOption {
	return schemaOptionFun(func(s *Schema) error {
		s.initDbRepeat = dbRepeat
		return nil
	})
}

func WithSchemaTableSize(tableSize int) SchemaOption {
	return schemaOptionFun(func(s *Schema) error {
		s.initTableSize = tableSize
		return nil
	})
}

type SchemaOption interface {
	apply(schema *Schema) error
}

type schemaOptionFun func(*Schema) error

func (f schemaOptionFun) apply(schema *Schema) error {
	return f(schema)
}

type SchemaTable struct {
	Database             string         `json:"database"`
	Table                string         `json:"table"`
	Size                 int            `json:"size"`
	Fields               []*SchemaField `json:"fields"`
	PrimaryKeyFieldNames []string       `json:"primaryKeyFieldNames"`
}

func (st *SchemaTable) Clone() *SchemaTable {
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

func readSchemaTablesFromFileOrTemplate(dialect Dialect, name string) (schemas []*SchemaTable, err error) {
	var b []byte
	switch name {
	case "", "sysbench":
		file := schemaFiles[dialect]["sysbench"]
		b, err = schemaF.ReadFile(file)
		if err != nil {
			return nil, err
		}
	default:
		b, err = os.ReadFile(name)
		if err != nil {
			return nil, err
		}
	}
	return schemas, json.Unmarshal(b, &schemas)
}
