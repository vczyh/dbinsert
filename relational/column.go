package relational

import (
	"fmt"
	"github.com/vczyh/dbinsert/generator"
)

type FieldType uint8

const (
	FieldTypeInt8 = iota
	FieldTypeInt16
	FieldTypeInt24
	FieldTypeInt32
	FieldTypeInt64
	FieldTypeChar
	FieldTypeVarChar
	FieldTypeText
)

func (ft FieldType) String() string {
	switch ft {
	case FieldTypeInt8:
		return "INT8"
	case FieldTypeChar:
		return "CHAR"
	default:
		return fmt.Sprintf("unknown field type: %d", ft)
	}
}

type Field struct {
	Name     string
	TypeName string

	autoIncrement bool
	fieldType     FieldType
	len           int
	precision     int
	scale         int
	generator     generator.Generator
}

func NewField(dialect Dialect, name, typeName string, autoIncrement bool) (*Field, error) {
	field := &Field{
		Name:     name,
		TypeName: typeName,
	}
	filedDesc, err := DialectParsers[dialect].FieldDesc(field.TypeName)
	if err != nil {
		return nil, err
	}
	field.fieldType = filedDesc.FieldType
	field.len = filedDesc.Len
	field.autoIncrement = filedDesc.IsAutoIncreased

	switch dialect {
	case DialectMysql:
		field.autoIncrement = autoIncrement
	}

	if err := field.defaultGenerator(); err != nil {
		return nil, err
	}
	return field, nil
}

func MustNewField(dialect Dialect, name, typeName string) *Field {
	field, err := NewField(dialect, name, typeName, false)
	if err != nil {
		panic(err)
	}
	return field
}

func (f *Field) GenerateData() interface{} {
	return f.generator.Generate()
}

func (f *Field) defaultGenerator() error {
	switch f.fieldType {
	case FieldTypeInt8, FieldTypeInt16, FieldTypeInt24, FieldTypeInt32, FieldTypeInt64:
		f.generator = generator.NewInt8()
	case FieldTypeChar, FieldTypeVarChar:
		g, err := generator.NewString(f.len)
		if err != nil {
			return err
		}
		f.generator = g
	default:
		return fmt.Errorf("unsupported field type: %s", f.fieldType.String())
	}

	return nil
}

func (f *Field) AutoIncrement() bool {
	return f.autoIncrement
}
