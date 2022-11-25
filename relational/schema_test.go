package relational

import (
	"encoding/json"
	"testing"
)

func TestGenerateSchema(t *testing.T) {
	tables := []*SchemaTable{
		{
			Database: "db1",
			Table:    "tbl1",
			Size:     1000,
			Fields: []*SchemaField{
				{Name: "id1", TypeName: "INT"},
				{Name: "val1", TypeName: "char(110)"},
			},
		},
		{
			Database: "db2",
			Table:    "tbl2",
			Size:     1000,
			Fields: []*SchemaField{
				{Name: "id2", TypeName: "INT"},
				{Name: "val2", TypeName: "char(110)"},
			},
		},
	}

	schema, err := NewSchema(DialectMysql, tables, 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	for db, _ := range schema.Databases() {
		t.Log(db)
	}

	b, err := json.MarshalIndent(tables, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
