package sqlite

import (
	"context"
	"database/sql"

	"github.com/vanillaiice/gocryptobot/pkg/db"
	_ "modernc.org/sqlite"
)

// DB stores a database connection.
type DB struct {
	conn *sql.DB
	ctx  context.Context
}

// tradingPairSymbol is the symbol of the trading pair.
var tradingPairSymbol string

// Open opens a database connection.
func Open(url, symbol string, ctx context.Context) (db *DB, err error) {
	conn, err := sql.Open("sqlite", url)
	if err != nil {
		return
	}

	stmt := `
		CREATE TABLE IF NOT EXISTS receipts(
			id INTEGER PRIMARY KEY,
			symbol TEXT NOT NULL,
			type TEXT NOT NULL,
			quantity REAL NOT NULL,
			profit REAL NOT NULL,
			price REAL NOT NULL,
			timestamp INTEGER NOT NULL
		);`

	if _, err = conn.ExecContext(ctx, stmt); err != nil {
		return
	}

	tradingPairSymbol = symbol

	return &DB{conn: conn, ctx: ctx}, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) GetLastReceipts(limit int) ([]*db.Receipt, error) {
	receipts := []*db.Receipt{}

	stmt := `
		SELECT symbol, type, quantity, profit, price, timestamp
		FROM receipts
		WHERE symbol = ?
		ORDER BY timestamp
		DESC LIMIT ?`

	rows, err := d.conn.QueryContext(d.ctx, stmt, tradingPairSymbol, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		receipt := &db.Receipt{}

		err = rows.Scan(
			&receipt.Symbol,
			&receipt.TxType,
			&receipt.Quantity,
			&receipt.Profit,
			&receipt.Price,
			&receipt.Timestamp,
		)

		if err != nil {
			return nil, err
		}

		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

func (d *DB) InsertReceipt(receipt *db.Receipt) (err error) {
	stmt := `
		INSERT INTO receipts(symbol, type, quantity, profit, price, timestamp)
		VALUES(?, ?, ?, ?, ?, ?)`

	_, err = d.conn.ExecContext(
		d.ctx,
		stmt,
		receipt.Symbol,
		receipt.TxType,
		receipt.Quantity,
		receipt.Profit,
		receipt.Price,
		receipt.Timestamp,
	)

	return
}
