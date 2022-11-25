package relational

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type Manager struct {
	cm                 ConnectionManager
	bi                 BulkInserter
	tables             []*Table
	timeout            time.Duration
	autoCreateDatabase bool
	autoCreateTable    bool

	db        *sql.DB
	databases map[string][]*Table
	dbs       sync.Map
	err       error
}

func NewManager(cm ConnectionManager, bi BulkInserter, tables []*Table, opts ...ManagerOption) (*Manager, error) {
	m := &Manager{
		cm:                 cm,
		bi:                 bi,
		tables:             tables,
		timeout:            10000 * time.Hour,
		autoCreateDatabase: false,
		autoCreateTable:    false,
	}

	for _, opt := range opts {
		if err := opt.apply(m); err != nil {
			return nil, err
		}
	}

	m.databases = make(map[string][]*Table)
	for _, table := range m.tables {
		m.databases[table.Database] = append(m.databases[table.Database], table)
	}

	return m, nil
}

func (m *Manager) Start(ctx context.Context) error {
	ctx, cancelFunc := context.WithTimeout(ctx, m.timeout)
	defer cancelFunc()

	if err := m.prepare(ctx); err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, table := range m.tables {
		wg.Add(1)
		go func(table *Table) {
			defer wg.Done()
			if err := m.batchInsert(ctx, table); err != nil {
				// TODO
				fmt.Println("insert error: ", err)
				if m.err == nil {
					m.err = err
				}
				cancelFunc()
			}
		}(table)
	}
	wg.Wait()

	m.close()
	return m.err
}

func (m *Manager) prepare(ctx context.Context) error {
	db, err := m.cm.Create(ctx)
	if err != nil {
		return err
	}
	m.db = db

	if err := m.db.PingContext(ctx); err != nil {
		return err
	}

	for database := range m.databases {
		db, err := m.open(ctx, database)
		if err != nil {
			return err
		}
		m.dbs.Store(database, db)
	}

	if m.autoCreateDatabase {
		if err := m.createDatabases(ctx); err != nil {
			return err
		}
	}

	if m.autoCreateTable {
		if err := m.createTable(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) createDatabases(ctx context.Context) error {
	for database := range m.databases {
		_, err := m.db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+database)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) createTable(ctx context.Context) error {
	for database, tables := range m.databases {
		db, err := m.GetDb(database)
		if err != nil {
			return err
		}
		for _, table := range tables {
			// TODO
			fmt.Println(table.DDL())
			_, err := db.ExecContext(ctx, table.DDL())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) batchInsert(ctx context.Context, table *Table) error {
	db, err := m.GetDb(table.Database)
	if err != nil {
		return err
	}

	var rows []map[string]interface{}
	for i := 0; i < table.Size; i++ {
		row := make(map[string]interface{})
		for _, field := range table.Fields {
			if !field.AutoIncrement() {
				row[field.Name] = field.GenerateData()
			}
		}
		rows = append(rows, row)
		// TODO
		if (i+1)%2500 == 0 {
			if err := m.bi.Insert(ctx, db, table, rows); err != nil {
				return err
			}
			rows = nil
		}
	}
	if len(rows) > 0 {
		return m.bi.Insert(ctx, db, table, rows)
	}

	return nil
}

func (m *Manager) close() {
	m.db.Close()
	m.dbs.Range(func(database, value interface{}) bool {
		db := value.(*sql.DB)
		db.Close()
		m.dbs.Delete(db)
		return true
	})
}

func (m *Manager) GetDb(database string) (*sql.DB, error) {
	value, ok := m.dbs.Load(database)
	if !ok {
		return nil, fmt.Errorf("db not found for %s", database)
	}
	return value.(*sql.DB), nil
}

func (m *Manager) open(ctx context.Context, database string) (*sql.DB, error) {
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
	//	m.cnf.Username,
	//	m.cnf.Password,
	//	m.cnf.Host,
	//	m.cnf.Port,
	//	database)
	//db, err := sql.Open("mysql", dsn)
	//if err != nil {
	//	return nil, err
	//}
	// See "Important settings" section.
	//db.SetConnMaxLifetime(time.Minute * 3)
	//db.SetMaxOpenConns(1)
	//db.SetMaxIdleConns(10)
	//return db, nil

	return m.cm.CreateInDatabase(ctx, database)
}

func WithManagerTimeout(timeout time.Duration) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.timeout = timeout
		return nil
	})
}

func WithManagerAutoCreateDatabase(autoCreateDatabase bool) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.autoCreateDatabase = autoCreateDatabase
		return nil
	})
}

func WithManagerAutoCreateTable(autoCreateTable bool) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.autoCreateTable = autoCreateTable
		return nil
	})
}

type ManagerOption interface {
	apply(manager *Manager) error
}

type managerOptionFun func(*Manager) error

func (f managerOptionFun) apply(manager *Manager) error {
	return f(manager)
}
