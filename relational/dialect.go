package relational

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	DialectParsers = map[Dialect]DialectParser{
		DialectMysql:    &MysqlDialectParser{},
		DialectPostgres: &PostgresDialectParser{},
	}
)

type DialectParser interface {
	FieldDesc(typeName string) (*FieldDataType, error)
}

type FieldDataType struct {
	FieldType       FieldType
	Len             int
	IsAutoIncreased bool
}

type MysqlDialectParser struct{}

func (p *MysqlDialectParser) FieldDesc(typeName string) (*FieldDataType, error) {
	tn := strings.ToUpper(typeName)
	switch tn {
	case "TINYINT":
		return &FieldDataType{FieldTypeInt8, 0, false}, nil
	case "SMALLINT":
		return &FieldDataType{FieldTypeInt16, 0, false}, nil
	case "MEDIUMINT":
		return &FieldDataType{FieldTypeInt24, 0, false}, nil
	case "INT", "INTEGER":
		return &FieldDataType{FieldTypeInt32, 0, false}, nil
	case "BIGINT":
		return &FieldDataType{FieldTypeInt64, 0, false}, nil
	case "TEXT":
		return &FieldDataType{FieldTypeText, 0, false}, nil
	}

	if strings.HasPrefix(tn, "CHAR") {
		length, err := typeLength(tn)
		if err != nil {
			return nil, err
		}
		return &FieldDataType{FieldTypeChar, length, false}, nil

	} else if strings.HasPrefix(tn, "VARCHAR") {
		length, err := typeLength(tn)
		if err != nil {
			return nil, err
		}
		return &FieldDataType{FieldTypeVarChar, length, false}, nil
	}

	return nil, fmt.Errorf("unsupported type: %s", typeName)
}

type PostgresDialectParser struct{}

func (p *PostgresDialectParser) FieldDesc(typeName string) (*FieldDataType, error) {
	tn := strings.ToLower(typeName)
	switch tn {
	case "smallint":
		return &FieldDataType{FieldTypeInt16, 0, false}, nil
	case "integer":
		return &FieldDataType{FieldTypeInt32, 0, false}, nil
	case "bigint":
		return &FieldDataType{FieldTypeInt64, 0, false}, nil
	case "smallserial":
		return &FieldDataType{FieldTypeInt16, 0, true}, nil
	case "serial":
		return &FieldDataType{FieldTypeInt32, 0, true}, nil
	case "bigserial":
		return &FieldDataType{FieldTypeInt64, 0, true}, nil
	case "text":
		return &FieldDataType{FieldTypeText, 0, false}, nil
	}

	if strings.HasPrefix(tn, "character varying") || strings.HasPrefix(tn, "varchar") {
		length, err := typeLength(tn)
		if err != nil {
			return nil, err
		}
		return &FieldDataType{FieldTypeVarChar, length, false}, nil

	} else if strings.HasPrefix(tn, "character") || strings.HasPrefix(tn, "char") {
		length, err := typeLength(tn)
		if err != nil {
			return nil, err
		}
		return &FieldDataType{FieldTypeChar, length, false}, nil
	}

	return nil, fmt.Errorf("unsupported type: %s", typeName)
}

func typeLength(typeName string) (int, error) {
	pattern := regexp.MustCompile(`\((\d+)\)`)
	lens := pattern.FindAllStringSubmatch(typeName, -1)
	if len(lens) != 1 {
		return 0, fmt.Errorf("unsupported type: %s", typeName)
	}
	return strconv.Atoi(lens[0][1])
}
