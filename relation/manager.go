package relation

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

const (
	bulkBuffer = 512 * 1024
)

type Manager struct {
	dialect            Dialect
	cm                 ConnectionManager
	bi                 BulkInserter
	tables             []*Table
	timeout            time.Duration
	autoCreateDatabase bool
	autoCreateTable    bool
	progress           *Progress

	db        *sql.DB
	databases map[string][]*Table
	dbs       sync.Map
	err       error
}

func NewManager(dialect Dialect, cm ConnectionManager, bi BulkInserter, tables []*Table, opts ...ManagerOption) (*Manager, error) {
	m := &Manager{
		dialect:            dialect,
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

	if m.progress == nil {
		progress, err := NewProgress(m.tables)
		if err != nil {
			return nil, err
		}
		m.progress = progress
	}

	m.databases = make(map[string][]*Table)
	for _, table := range m.tables {
		m.databases[table.Database] = append(m.databases[table.Database], table)
	}

	return m, nil
}

func (m *Manager) Start(ctx context.Context) error {
	defer m.close()
	ctx, cancelFunc := context.WithTimeout(ctx, m.timeout)
	defer cancelFunc()

	if err := m.prepare(ctx); err != nil {
		return err
	}

	go m.progress.Render()
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
	if m.err != nil {
		return m.err
	}
	for {
		if m.progress.Ended() {
			break
		}
	}
	return nil
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
		if err := m.createDatabase(ctx, database); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) createDatabase(ctx context.Context, name string) error {
	switch m.dialect {
	case DialectPostgres:
		return m.createDatabaseForPostgres(ctx, name)
	default:
		_, err := m.db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+name)
		return err
	}
}

func (m *Manager) createDatabaseForPostgres(ctx context.Context, name string) error {
	s := fmt.Sprintf("SELECT COUNT(*) FROM pg_database WHERE datname = '%s'", name)
	row := m.db.QueryRowContext(ctx, s)
	var cnt int
	if err := row.Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}

	_, err := m.db.ExecContext(ctx, "CREATE DATABASE "+name)
	return err
}

func (m *Manager) createTable(ctx context.Context) error {
	for database, tables := range m.databases {
		db, err := m.GetDb(database)
		if err != nil {
			return err
		}
		for _, table := range tables {
			switch m.dialect {
			case DialectMysql:
				_, err := db.ExecContext(ctx, table.DDL())
				if err != nil {
					return err
				}
			case DialectPostgres:
				_, err := db.ExecContext(ctx, table.PostgresDDL())
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported dialect: %s", m.dialect)
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
	rowsSize := 0
	for i := 0; i < table.Size; i++ {
		row := make(map[string]interface{})
		for _, field := range table.Fields {
			if !field.AutoIncrement() {
				value := field.GenerateData()
				row[field.Name] = value
				switch v := value.(type) {
				case int, uint:
					rowsSize += 8
				case int8, uint8:
					rowsSize += 1
				case int16, uint16:
					rowsSize += 2
				case int32, uint32:
					rowsSize += 4
				case int64, uint64:
					rowsSize += 8
				case []byte:
					rowsSize += len(v)
				case string:
					rowsSize += len(v)
				}
			}
		}
		rows = append(rows, row)
		//if (i+1)%2500 == 0 {
		if rowsSize > bulkBuffer {
			if err := m.bi.Insert(ctx, db, table, rows); err != nil {
				return err
			}
			m.progress.Increment(table, len(rows))
			rows = nil
			rowsSize = 0
		}
	}
	if len(rows) > 0 {
		if err := m.bi.Insert(ctx, db, table, rows); err != nil {
			return err
		}
		m.progress.Increment(table, len(rows))
	}

	return nil
}

func (m *Manager) close() {
	m.progress.Stop()
	m.db.Close()
	m.dbs.Range(func(database, value interface{}) bool {
		db := value.(*sql.DB)
		db.Close()
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
