package bbolt

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/vanillaiice/gocryptobot/pkg/db"
	"go.etcd.io/bbolt"
)

type DB struct {
	conn *bbolt.DB
}

var tradingPairSymbol string

func Open(url, symbol string) (db *DB, err error) {
	conn, err := bbolt.Open(url, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return
	}

	err = conn.Update(func(tx *bbolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(symbol))
		return err
	})
	if err != nil {
		return
	}

	tradingPairSymbol = symbol

	return &DB{conn: conn}, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) GetLastReceipts(limit int) ([]*db.Receipt, error) {
	receipts := []*db.Receipt{}

	err := d.conn.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(tradingPairSymbol))

		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var receipt db.Receipt
			err := json.Unmarshal(v, &receipt)
			if err != nil {
				return err
			}

			receipts = append(receipts, &receipt)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return receipts, err
}

func (d *DB) InsertReceipt(receipt *db.Receipt) (err error) {
	return d.conn.Update(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(tradingPairSymbol))

		id, err := b.NextSequence()
		if err != nil {
			return
		}

		receiptJSON, err := json.Marshal(receipt)
		if err != nil {
			return
		}

		return b.Put([]byte(strconv.FormatInt(int64(id), 2)), receiptJSON)
	})
}
