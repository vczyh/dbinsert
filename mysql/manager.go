package mysql

//import (
//	"context"
//	"database/sql"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/vczyh/dbinsert/relational"
//	"sync"
//)
//
//type Manager struct {
//	cnf       *Config
//	db        *sql.DB
//	databases map[string][]*relational.Table
//	dbs       sync.Map
//	err       error
//}
//
//func NewManager(cnf *Config) (*Manager, error) {
//	m := new(Manager)
//	m.cnf = cnf
//
//	db, err := m.open("information_schema")
//	if err != nil {
//		return nil, err
//	}
//	m.db = db
//
//	m.databases = make(map[string][]*relational.Table)
//	for _, table := range m.cnf.Tables {
//		m.databases[table.Database] = append(m.databases[table.Database], table)
//	}
//
//	return m, nil
//}
//
//func (m *Manager) Start() error {
//	ctx, cancelFunc := context.WithTimeout(context.Background(), m.cnf.Timeout)
//	defer cancelFunc()
//
//	if err := m.prepare(ctx); err != nil {
//		return err
//	}
//
//	var wg sync.WaitGroup
//	for _, table := range m.cnf.Tables {
//		wg.Add(1)
//		go func(table *relational.Table) {
//			defer wg.Done()
//			if err := m.batchInsert(ctx, table); err != nil {
//				fmt.Println("insert error: ", err)
//				if m.err == nil {
//					m.err = err
//				}
//				cancelFunc()
//			}
//		}(table)
//	}
//	wg.Wait()
//
//	m.close()
//	return m.err
//}
//
//func (m *Manager) prepare(ctx context.Context) error {
//	if err := m.db.PingContext(ctx); err != nil {
//		return err
//	}
//
//	for database := range m.databases {
//		db, err := m.open(database)
//		if err != nil {
//			return err
//		}
//		m.dbs.Store(database, db)
//	}
//
//	if m.cnf.CreateDatabase {
//		if err := m.createDatabases(ctx); err != nil {
//			return err
//		}
//	}
//
//	if m.cnf.CreateTable {
//		if err := m.createTable(ctx); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func (m *Manager) createDatabases(ctx context.Context) error {
//	for database := range m.databases {
//		_, err := m.db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+database)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (m *Manager) createTable(ctx context.Context) error {
//	for database, tables := range m.databases {
//		db, err := m.GetDb(database)
//		if err != nil {
//			return err
//		}
//		for _, table := range tables {
//			// TODO
//			fmt.Println(table.DDL())
//			_, err := db.ExecContext(ctx, table.DDL())
//			if err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//func (m *Manager) batchInsert(ctx context.Context, table *relational.Table) error {
//	db, err := m.GetDb(table.Database)
//	if err != nil {
//		return err
//	}
//
//	batchInsertStmt, err := NewBatchInsertStmt(table.Name, table.NotAutoIncrementColumnNames())
//	if err != nil {
//		return err
//	}
//
//	for i := 0; i < table.Size; i++ {
//		for _, field := range table.Fields {
//			if !field.AutoIncrement() {
//				batchInsertStmt.Set(field.Name, field.GenerateData())
//			}
//		}
//		batchInsertStmt.AddBatch()
//		if (i+1)%2500 == 0 {
//			if err = batchInsertStmt.ExecuteBatch(ctx, db); err != nil {
//				return err
//			}
//			batchInsertStmt.CleanBatch()
//		}
//	}
//
//	if batchInsertStmt.HaveBatch() {
//		if err = batchInsertStmt.ExecuteBatch(ctx, db); err != nil {
//			return err
//		}
//		batchInsertStmt.CleanBatch()
//	}
//	return nil
//}
//
//func (m *Manager) close() {
//	m.db.Close()
//	m.dbs.Range(func(database, value interface{}) bool {
//		db := value.(*sql.DB)
//		db.Close()
//		m.dbs.Delete(db)
//		return true
//	})
//}
//
//func (m *Manager) GetDb(database string) (*sql.DB, error) {
//	value, ok := m.dbs.Load(database)
//	if !ok {
//		return nil, fmt.Errorf("db not found for %s", database)
//	}
//	return value.(*sql.DB), nil
//}
//
//func (m *Manager) open(database string) (*sql.DB, error) {
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
//		m.cnf.Username,
//		m.cnf.Password,
//		m.cnf.Host,
//		m.cnf.Port,
//		database)
//	db, err := sql.Open("mysql", dsn)
//	if err != nil {
//		return nil, err
//	}
//	// See "Important settings" section.
//	//db.SetConnMaxLifetime(time.Minute * 3)
//	db.SetMaxOpenConns(1)
//	//db.SetMaxIdleConns(10)
//	return db, nil
//}
