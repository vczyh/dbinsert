package mysql

import (
	"github.com/vczyh/dbinsert/relational"
	"time"
)

type Config struct {
	Host           string
	Port           int
	Username       string
	Password       string
	CreateTable    bool
	CreateDatabase bool
	Tables         []*relational.Table
	Timeout        time.Duration
}
