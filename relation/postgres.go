package relation

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"net"
	"net/url"
	"strconv"
)

type PostgresConnectionManager struct {
	host     string
	port     int
	user     string
	password string
}

func NewPostgresConnectionManager(host, user, password string, opts ...PostgresConnectionManagerOption) (*PostgresConnectionManager, error) {
	cm := &PostgresConnectionManager{
		host:     host,
		port:     5432,
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

func (cm *PostgresConnectionManager) Create(ctx context.Context) (*sql.DB, error) {
	return cm.open("postgres")
}

func (cm *PostgresConnectionManager) CreateInDatabase(ctx context.Context, database string) (*sql.DB, error) {
	return cm.open(database)
}

func (cm *PostgresConnectionManager) open(database string) (*sql.DB, error) {
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

func WithPostgresConnectionPort(port int) PostgresConnectionManagerOption {
	return postgresConnectionManagerOptionFun(func(cm *PostgresConnectionManager) error {
		cm.port = port
		return nil
	})
}

type PostgresConnectionManagerOption interface {
	apply(manager *PostgresConnectionManager) error
}

type postgresConnectionManagerOptionFun func(manager *PostgresConnectionManager) error

func (f postgresConnectionManagerOptionFun) apply(manager *PostgresConnectionManager) error {
	return f(manager)
}
