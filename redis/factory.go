package redis

import "time"

type Config struct {
	User       string
	Password   string
	Timeout    time.Duration
	KeyCount   int
	ValueLen   int
	EnableTLS  bool
	CaCert     string
	SkipVerify bool
	Cert       string
	Key        string

	// Standalone or Master-slave
	Host string
	Port int

	// Cluster
	Addresses   []string
	ClusterMode bool
}

func CreateManager(cnf *Config) (*Manager, error) {
	return NewManager(
		cnf.Host,
		cnf.User,
		cnf.Password,
		WithManagerOptionPort(cnf.Port),
		WithManagerOptionKeyCount(cnf.KeyCount),
		WithManagerOptionTimeout(cnf.Timeout),
		WithManagerOptionValueLen(cnf.ValueLen),
		WithManagerOptionIsCluster(cnf.ClusterMode),
		WithManagerOptionAddresses(cnf.Addresses),
		WithManagerOptionEnableTLS(cnf.EnableTLS),
		WithManagerOptionCaCert(cnf.CaCert),
		WithManagerOptionSkipVerify(cnf.SkipVerify),
		WithManagerOptionCert(cnf.Cert),
		WithManagerOptionKey(cnf.Key),
	)
}
