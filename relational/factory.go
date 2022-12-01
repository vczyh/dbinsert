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

type PostgresConfig struct {
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
	cm, err := NewMySQLConnectionManager(
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
		DialectMysql,
		cm,
		bi,
		tables,
		WithManagerAutoCreateDatabase(c.CreateDatabase),
		WithManagerAutoCreateTable(c.CreateTable),
		WithManagerTimeout(c.Timeout))
}

func CreateManagerForPostgres(c *PostgresConfig) (*Manager, error) {
	cm, err := NewPostgresConnectionManager(
		c.Host,
		c.Username,
		c.Password,
		WithPostgresConnectionPort(c.Port))
	if err != nil {
		return nil, err
	}

	bi, err := NewBulkInsert()
	if err != nil {
		return nil, err
	}

	tables, err := NewTablesFromSchemaFile(
		DialectPostgres,
		c.SchemaFile,
		WithSchemaDatabasesRepeat(c.DatabaseRepeat),
		WithSchemaTableSize(c.TableSize))
	if err != nil {
		return nil, err
	}

	return NewManager(
		DialectPostgres,
		cm,
		bi,
		tables,
		WithManagerAutoCreateDatabase(c.CreateDatabase),
		WithManagerAutoCreateTable(c.CreateTable),
		WithManagerTimeout(c.Timeout))
}
