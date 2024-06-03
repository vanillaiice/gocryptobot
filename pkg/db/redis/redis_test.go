package redis_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/vanillaiice/gocryptobot/pkg/db"
	"github.com/vanillaiice/gocryptobot/pkg/db/redis"
)

var DB *redis.DB

var receipts = []*db.Receipt{
	{
		Symbol:    "BTCUSDT",
		TxType:    "BUY",
		Quantity:  1.0,
		Profit:    1.0,
		Price:     1.0,
		Timestamp: 1,
	},
	{
		Symbol:    "BTCUSDT",
		TxType:    "SELL",
		Quantity:  1.0,
		Profit:    1.0,
		Price:     1.0,
		Timestamp: 2,
	},
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(err)
	}

	if err = pool.Client.Ping(); err != nil {
		log.Fatal(err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7.2.5-alpine",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatal(err)
	}

	addr := net.JoinHostPort("localhost", resource.GetPort("6379/tcp"))
	dbURL := fmt.Sprintf("redis://%s", addr)

	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		DB, err = redis.Open(dbURL, "BTCUSDT", context.Background())
		return err
	}); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestInsertReceipt(t *testing.T) {
	err := DB.InsertReceipt(receipts[1])
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetLastReceipts(t *testing.T) {
	err := DB.InsertReceipt(receipts[0])
	if err != nil {
		t.Fatal(err)
	}

	n := 2

	r, err := DB.GetLastReceipts(n)
	if err != nil {
		t.Fatal(err)
	}

	if len(receipts) == 0 {
		t.Fatalf("no receipts")
	}

	if len(receipts) != n {
		t.Fatalf("got %d, want %d", len(receipts), n)
	}

	for i, receipt := range receipts {
		if receipt == r[i] {
			t.Errorf("got %v, want %v", r[i], receipt)
		}
	}
}
