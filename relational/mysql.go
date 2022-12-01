package relational

import (
	"context"
	"database/sql"
	"fmt"
)

type MySQLConnectionManager struct {
	host     string
	port     int
	user     string
	password string
}

func NewMySQLConnectionManager(host, user, password string, opts ...MysqlConnectionManagerOption) (*MySQLConnectionManager, error) {
	cm := &MySQLConnectionManager{
		host:     host,
		port:     3306,
		user:     user,
		password: password,
	}
	for _, opt := range opts {
		if err := opt.apply(cm); err != nil {
			return nil, err
		}
	}

	return cm, nil
}

func (cm *MySQLConnectionManager) Create(ctx context.Context) (*sql.DB, error) {
	return cm.open("information_schema")
}

func (cm *MySQLConnectionManager) CreateInDatabase(ctx context.Context, database string) (*sql.DB, error) {
	return cm.open(database)
}

func (cm *MySQLConnectionManager) open(database string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cm.user,
		cm.password,
		cm.host,
		cm.port,
		database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	//db.SetMaxOpenConns(1)
	//db.SetMaxIdleConns(10)
	return db, nil
}

func WithConnectionPort(port int) MysqlConnectionManagerOption {
	return mysqlConnectionManagerOptionFun(func(cm *MySQLConnectionManager) error {
		cm.port = port
		return nil
	})
}

type MysqlConnectionManagerOption interface {
	apply(manager *MySQLConnectionManager) error
}

type mysqlConnectionManagerOptionFun func(manager *MySQLConnectionManager) error

func (f mysqlConnectionManagerOptionFun) apply(manager *MySQLConnectionManager) error {
	return f(manager)
}

type BulkInsert struct {
	//db              *sql.DB
	//table           *relational.Table
	//batchInsertStmt *relational.BatchInsertStmt
}

//func NewBulkInsert(db *sql.DB, table *relational.Table) (*BulkInsert, error) {
//	bi := &BulkInsert{
//		db:    db,
//		table: table,
//	}
//
//	batchInsertStmt, err := relational.NewBatchInsertStmt(bi.table.Name, bi.table.NotAutoIncrementColumnNames())
//	if err != nil {
//		return nil, err
//	}
//	bi.batchInsertStmt = batchInsertStmt
//	return bi, nil
//}

func NewBulkInsert() (*BulkInsert, error) {
	return &BulkInsert{}, nil
}

func (bi *BulkInsert) Insert(ctx context.Context, db *sql.DB, table *Table, rows []map[string]interface{}) error {
	batchInsertStmt, err := NewBatchInsertStmt(table.Name, table.NotAutoIncrementColumnNames())
	if err != nil {
		return err
	}

	for _, row := range rows {
		for _, field := range table.Fields {
			if !field.AutoIncrement() {
				batchInsertStmt.Set(field.Name, row[field.Name])
			}
		}
		batchInsertStmt.AddBatch()
	}

	if batchInsertStmt.HaveBatch() {
		if err := batchInsertStmt.ExecuteBatch(ctx, db); err != nil {
			return err
		}
		batchInsertStmt.CleanBatch()
	}

	return nil
}
