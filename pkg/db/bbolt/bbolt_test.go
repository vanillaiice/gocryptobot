package bbolt_test

import (
	"log"
	"os"
	"testing"

	"github.com/vanillaiice/gocryptobot/pkg/db"
	"github.com/vanillaiice/gocryptobot/pkg/db/bbolt"
)

var DB *bbolt.DB

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
	var err error
	DB, err = bbolt.Open("test.db", "BTCUSDT")
	if err != nil {
		log.Fatal(err)
	}

	defer DB.Close()

	code := m.Run()

	os.Remove("test.db")

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
		t.Fatal("no receipts")
	}

	if len(receipts) != n {
		t.Fatalf("got %d, want %d", len(receipts), n)
	}

	for i, receipt := range receipts {
		if receipt == r[i] {
			t.Fatalf("got %v, want %v", r[i], receipt)
		}
	}
}
