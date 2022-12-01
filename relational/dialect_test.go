package relational

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMysqlDialectParser_FieldDesc(t *testing.T) {
	parser := DialectParsers[DialectMysql]

	cases := []struct {
		dataType  string
		fieldType FieldType
		len       int
	}{
		{"TINYINT", FieldTypeInt8, 0},
		{"SMALLINT", FieldTypeInt16, 0},
		{"MEDIUMINT", FieldTypeInt24, 0},
		{"INT", FieldTypeInt32, 0},
		{"BIGINT", FieldTypeInt64, 0},
		{"BIGINT", FieldTypeInt64, 0},
		{"BIGINT", FieldTypeInt64, 0},
		{"CHAR(60)", FieldTypeChar, 60},
		{"VARCHAR(60)", FieldTypeVarChar, 60},
		{"TEXT", FieldTypeText, 0},
	}

	for _, c := range cases {
		t.Run(c.dataType, func(t *testing.T) {
			desc, err := parser.FieldDesc(c.dataType)
			if err != nil {
				t.Fatal(err)
			}
			assert.NotNil(t, desc)
			assert.Equal(t, c.fieldType, desc.FieldType)
			assert.Equal(t, c.len, desc.Len)
			assert.Equal(t, false, desc.IsAutoIncreased)
		})
	}
}

func TestPostgresDialectParser_FieldDesc(t *testing.T) {
	parser := DialectParsers[DialectPostgres]

	cases := []struct {
		dataType                string
		expectedFieldType       FieldType
		expectedLen             int
		expectedIsAutoIncreased bool
	}{
		{"smallint", FieldTypeInt16, 0, false},
		{"integer", FieldTypeInt32, 0, false},
		{"bigint", FieldTypeInt64, 0, false},
		{"smallserial", FieldTypeInt16, 0, true},
		{"serial", FieldTypeInt32, 0, true},
		{"bigserial", FieldTypeInt64, 0, true},
		{"character varying(60)", FieldTypeVarChar, 60, false},
		{"varchar(60)", FieldTypeVarChar, 60, false},
		{"character(60)", FieldTypeChar, 60, false},
		{"char(60)", FieldTypeChar, 60, false},
		{"text", FieldTypeText, 0, false},
	}

	for _, c := range cases {
		t.Run(c.dataType, func(t *testing.T) {
			desc, err := parser.FieldDesc(c.dataType)
			if err != nil {
				t.Fatal(err)
			}
			assert.NotNil(t, desc)
			assert.Equal(t, c.expectedFieldType, desc.FieldType)
			assert.Equal(t, c.expectedLen, desc.Len)
			assert.Equal(t, c.expectedIsAutoIncreased, desc.IsAutoIncreased)
		})
	}
}
