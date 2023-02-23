package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestInsert(t *testing.T) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"100.100.5.222:6379"},
		Username: "unicloud",
		Password: "Zggyy2019!",
	})

	//rdb := redis.NewClient(&redis.Options{
	//	Addr:                  "100.100.5.222:6379",
	//	Username:              "unicloud",
	//	Password:              "Zggyy2019!",
	//	DialTimeout:           2 * time.Second,
	//	WriteTimeout:          2 * time.Second,
	//	ReadTimeout:           2 * time.Second,
	//	ContextTimeoutEnabled: true,
	//})
	//defer rdb.Close()

	pipeline := rdb.Pipeline()
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("dbinsert_%d", i)
		statusCmd := pipeline.Set(context.Background(), key, key, 0)
		if err := statusCmd.Err(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pipeline.Exec(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
