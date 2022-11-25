package postgres

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"strconv"
)

type ConnectionManager struct {
	host     string
	port     int
	user     string
	password string
}

func NewConnectionManager(host, user, password string, opts ...ConnectionManagerOption) (*ConnectionManager, error) {
	cm := &ConnectionManager{
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

func (cm *ConnectionManager) CreateInDatabase(ctx context.Context, database string) (*sql.DB, error) {
	query := make(url.Values)
	query.Set("connect_timeout", "5")
	query.Set("sslmode", "disable")
	//if c.sslRootCert != "" {
	//	query.Set("sslmode", "verify-ca")
	//	query.Set("sslrootcert", c.sslRootCert)
	//}

	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cm.user, cm.password),
		Host:     net.JoinHostPort(cm.host, strconv.Itoa(cm.port)),
		Path:     database,
		RawQuery: query.Encode(),
	}
	return sql.Open("postgres", dsn.String())
}

func WithConnectionPort(port int) ConnectionManagerOption {
	return connectionManagerOptionFun(func(cm *ConnectionManager) error {
		cm.port = port
		return nil
	})
}

type ConnectionManagerOption interface {
	apply(manager *ConnectionManager) error
}

type connectionManagerOptionFun func(manager *ConnectionManager) error

func (f connectionManagerOptionFun) apply(manager *ConnectionManager) error {
	return f(manager)
}
