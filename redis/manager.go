package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/vczyh/dbinsert/generator"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	bulkBuffer = 512 * 1024
)

type Manager struct {
	isCluster bool

	user       string
	password   string
	timeout    time.Duration
	keyCount   int
	valueLen   int
	enableTLS  bool
	caCert     string
	skipVerify bool
	cert       string
	key        string

	// Standalone or Master-slave
	host string
	port int

	// Cluster
	addresses []string

	rdb        *redis.Client
	rdbCluster *redis.ClusterClient
	tlsConfig  *tls.Config
}

func NewManager(host, user, password string, opts ...ManagerOption) (*Manager, error) {
	m := &Manager{
		host:     host,
		user:     user,
		password: password,
	}

	for _, opt := range opts {
		if err := opt.apply(m); err != nil {
			return nil, err
		}
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

	if err := m.insert(ctx); err != nil {
		fmt.Println("insert error: ", err)
		return err
	}

	return nil
}

func (m *Manager) close() {
	if m.isCluster {
		m.rdbCluster.Close()
	} else {
		m.rdb.Close()
	}
}

func (m *Manager) prepare(ctx context.Context) error {
	if m.enableTLS {
		caCertBytes, err := os.ReadFile(m.caCert)
		if err != nil {
			return err
		}
		rootCertPool := x509.NewCertPool()
		if ok := rootCertPool.AppendCertsFromPEM(caCertBytes); !ok {
			return fmt.Errorf("fail to append PEM")
		}

		m.tlsConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			RootCAs:            rootCertPool,
			InsecureSkipVerify: m.skipVerify,
			//ServerName: "*",
		}

		if m.cert != "" && m.key != "" {
			certKeyPair, err := tls.LoadX509KeyPair(m.cert, m.key)
			if err != nil {
				return err
			}
			m.tlsConfig.Certificates = []tls.Certificate{certKeyPair}
		}
	}

	if m.isCluster {
		m.rdbCluster = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:                 m.addresses,
			Username:              m.user,
			Password:              m.password,
			DialTimeout:           2 * time.Second,
			WriteTimeout:          2 * time.Second,
			ReadTimeout:           2 * time.Second,
			ContextTimeoutEnabled: true,
			TLSConfig:             m.tlsConfig,
		})
	} else {
		m.rdb = redis.NewClient(&redis.Options{
			Addr:                  net.JoinHostPort(m.host, strconv.Itoa(m.port)),
			Username:              m.user,
			Password:              m.password,
			DialTimeout:           2 * time.Second,
			WriteTimeout:          2 * time.Second,
			ReadTimeout:           2 * time.Second,
			ContextTimeoutEnabled: true,
			TLSConfig:             m.tlsConfig,
		})
	}

	_, err := m.cmd().Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("ping failed")
	}

	return nil
}

func (m *Manager) insert(ctx context.Context) error {
	stringGenerator, err := generator.NewString(m.valueLen)
	if err != nil {
		return err
	}
	pipe := m.cmd().Pipeline()

	curKeyCnt := 0
	bufKeyCnt := 0
	bufSize := 0
	for {
		key := fmt.Sprintf("dbinsert_%d", curKeyCnt+bufKeyCnt)
		value := stringGenerator.Generate()
		switch v := value.(type) {
		case string:
			bufSize += len(v)
		}

		statusCmd := pipe.Set(ctx, key, value, 0)
		if err := statusCmd.Err(); err != nil {
			return err
		}
		bufKeyCnt++

		isOver := m.keyCount != 0 && curKeyCnt+bufKeyCnt == m.keyCount
		if isOver || bufSize >= bulkBuffer {
			if _, err := pipe.Exec(ctx); err != nil {
				return err
			}
			curKeyCnt += bufKeyCnt
			bufKeyCnt = 0
			bufSize = 0
			if isOver {
				break
			}
		}
	}

	return nil
}

func (m *Manager) cmd() redis.Cmdable {
	if m.isCluster {
		return m.rdbCluster
	}
	return m.rdb
}

func WithManagerOptionPort(port int) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.port = port
		return nil
	})
}

func WithManagerOptionKeyCount(keyCount int) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.keyCount = keyCount
		return nil
	})
}

func WithManagerOptionTimeout(timeout time.Duration) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.timeout = timeout
		return nil
	})
}

func WithManagerOptionValueLen(valueLen int) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.valueLen = valueLen
		return nil
	})
}

func WithManagerOptionIsCluster(isCluster bool) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.isCluster = isCluster
		return nil
	})
}

func WithManagerOptionAddresses(addresses []string) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.addresses = addresses
		return nil
	})
}

func WithManagerOptionEnableTLS(enableTLS bool) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.enableTLS = enableTLS
		return nil
	})
}

func WithManagerOptionCaCert(caCert string) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.caCert = caCert
		return nil
	})
}

func WithManagerOptionSkipVerify(skipVerify bool) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.skipVerify = skipVerify
		return nil
	})
}

func WithManagerOptionCert(cert string) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.cert = cert
		return nil
	})
}

func WithManagerOptionKey(key string) ManagerOption {
	return managerOptionFun(func(m *Manager) error {
		m.key = key
		return nil
	})
}

type ManagerOption interface {
	apply(manager *Manager) error
}

type managerOptionFun func(manager *Manager) error

func (f managerOptionFun) apply(manager *Manager) error {
	return f(manager)
}
