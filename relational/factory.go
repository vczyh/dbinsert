package relational

import (
	"time"
)

type MysqlConfig struct {
	SchemaFile     string
	Host           string
	Port           int
	Username       string
	Password       string
	CreateTable    bool
	CreateDatabase bool
	Timeout        time.Duration
	DatabaseRepeat int
	TableSize      int
}

func CreateManagerForMysql(c *MysqlConfig) (*Manager, error) {
	cm, err := NewConnectionManager(
		c.Host,
		c.Username,
		c.Password,
		WithConnectionPort(c.Port))
	if err != nil {
		return nil, err
	}

	bi, err := NewBulkInsert()
	if err != nil {
		return nil, err
	}

	tables, err := NewTablesFromSchemaFile(
		DialectMysql,
		c.SchemaFile,
		WithSchemaDatabasesRepeat(c.DatabaseRepeat),
		WithSchemaTableSize(c.TableSize))
	if err != nil {
		return nil, err
	}

	return NewManager(
		cm,
		bi,
		tables,
		WithManagerAutoCreateDatabase(c.CreateDatabase),
		WithManagerAutoCreateTable(c.CreateTable),
		WithManagerTimeout(c.Timeout))
}
