package mysql

import (
	"github.com/vczyh/dbinsert/relational"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	cnf := &Config{
		Host:           "100.100.1.194",
		Port:           3306,
		Username:       "cloudos",
		Password:       "Zggyy2019!",
		CreateTable:    true,
		CreateDatabase: true,
		Tables:         nil,
		Timeout:        5 * time.Second,
	}

	cnf.Tables = append(cnf.Tables, relational.DefaultTable)

	m, err := NewManager(cnf)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Start(); err != nil {
		t.Fatal(err)
	}
}
