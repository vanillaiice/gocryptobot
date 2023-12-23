package simple

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Receipt struct {
	TxType    Tx
	Quantity  float32
	Price     float32
	Profit    float32
	Timestamp int64
}

func DBOpen(symbol string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	db, err = sql.Open("sqlite3", fmt.Sprintf("db/tx_receipts/%s.db", symbol))
	if err != nil {
		return db, err
	}

	stmt := `CREATE TABLE IF NOT EXISTS receipts(id INTEGER PRIMARY KEY, type TEXT, quantity REAL, profit REAL, cumProfit REAL, price REAL, timestamp INTEGER)`
	_, err = db.Exec(stmt)

	if err != nil {
		return db, err
	}

	return db, nil
}

func DBGetLastProfit(db *sql.DB) (float32, error) {
	var lastProfit float32
	err := db.QueryRow(`SELECT cumProfit FROM receipts ORDER BY ID DESC LIMIT 1`).Scan(&lastProfit)
	if err == sql.ErrNoRows {
		return float32(0), nil
	} else if err != nil {
		return lastProfit, err
	}
	return lastProfit, nil
}

func DBInsertReceipt(receipt *Receipt, db *sql.DB) error {
	lastProfit, err := DBGetLastProfit(db)
	if err != nil {
		return err
	}
	stmt := fmt.Sprintf(
		"INSERT INTO receipts(type, quantity, profit, cumProfit, price, timestamp) VALUES(%q, %.5f, %.5f, %.5f, %.5f, %d)",
		receipt.TxType,
		receipt.Quantity,
		receipt.Profit,
		lastProfit+receipt.Profit,
		receipt.Price,
		receipt.Timestamp,
	)
	_, err = db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
