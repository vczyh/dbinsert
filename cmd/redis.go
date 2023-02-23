package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/vczyh/dbinsert/redis"
	"time"
)

// represents the redis command
var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Quick insert tool for redis.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return StartRedis()
	},
}

var (
	redisCnf = new(redis.Config)
)

func init() {
	rootCmd.AddCommand(redisCmd)

	redisCmd.Flags().StringVar(&redisCnf.User, "user", "default", "redis username")
	redisCmd.Flags().StringVar(&redisCnf.Password, "password", "", "redis password")
	redisCmd.Flags().DurationVar(&redisCnf.Timeout, "timeout", 10*time.Hour, "timeout")
	redisCmd.Flags().IntVar(&redisCnf.KeyCount, "key-count", 0, "key count")
	redisCmd.Flags().IntVar(&redisCnf.ValueLen, "value-len", 50, "value string length")

	// Standalone or Master-slave
	redisCmd.Flags().StringVar(&redisCnf.Host, "host", "127.0.0.1", "redis host")
	redisCmd.Flags().IntVar(&redisCnf.Port, "port", 6379, "redis port")

	// Cluster
	redisCmd.Flags().BoolVar(&redisCnf.ClusterMode, "cluster", false, "enable cluster mode")
	redisCmd.Flags().StringArrayVar(&redisCnf.Addresses, "addrs", []string{"127.0.0.1:6379"}, "cluster addresses")
}

func StartRedis() error {
	manager, err := redis.CreateManager(redisCnf)
	if err != nil {
		return err
	}
	return manager.Start(context.Background())
}
