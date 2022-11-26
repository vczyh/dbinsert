package relational

import (
	"fmt"
	"github.com/vczyh/dbinsert/generator"
	"regexp"
	"strconv"
	"strings"
)

type FieldType uint8

const (
	FieldTypeInt = iota
	FieldTypeChar
)

func (ft FieldType) String() string {
	switch ft {
	case FieldTypeInt:
		return "INT"
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

func NewField(name, typeName string, autoIncrement bool) (*Field, error) {
	tn := strings.ToUpper(typeName)

	field := new(Field)
	field.Name = name
	field.TypeName = typeName
	field.autoIncrement = autoIncrement

	switch {
	case tn == "INT":
		field.fieldType = FieldTypeInt

	case strings.HasPrefix(tn, "CHAR"):
		field.fieldType = FieldTypeChar
		pattern := regexp.MustCompile(`(\d+)`)
		lens := pattern.FindAllStringSubmatch(tn, -1)
		if len(lens) != 1 {
			return nil, fmt.Errorf("unsupported type: %s", tn)
		}
		len, err := strconv.Atoi(lens[0][1])
		if err != nil {
			return nil, err
		}
		field.len = len

	default:
		return nil, fmt.Errorf("unsupported type: %s", name)
	}

	if err := field.defaultGenerator(); err != nil {
		return nil, err
	}
	return field, nil
}

func MustNewField(name, typeName string) *Field {
	field, err := NewField(name, typeName, false)
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
	case FieldTypeInt:
		f.generator = generator.NewInt8()
	case FieldTypeChar:
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
