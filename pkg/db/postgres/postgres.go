package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/vanillaiice/gocryptobot/pkg/db"
)

type DB struct {
	conn *pgx.Conn
	ctx  context.Context
}

var tradingPairSymbol string

func Open(url, symbol string, ctx context.Context) (db *DB, err error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return
	}

	stmt :=
		`CREATE TABLE IF NOT EXISTS receipts(
			id SERIAL PRIMARY KEY,
			symbol TEXT NOT NULL,
			type TEXT,
			quantity REAL,
			profit REAL,
			price REAL,
			timestamp INTEGER
		);`

	if _, err = conn.Exec(ctx, stmt); err != nil {
		return
	}

	tradingPairSymbol = symbol

	return &DB{conn: conn, ctx: ctx}, nil
}

func (d *DB) Close() error {
	return d.conn.Close(d.ctx)
}

func (d *DB) GetLastReceipts(limit int) ([]*db.Receipt, error) {
	receipts := []*db.Receipt{}

	stmt := `
		SELECT symbol, type, quantity, profit, price, timestamp
		FROM receipts
		WHERE symbol = $1
		ORDER BY timestamp DESC
		LIMIT $2;`

	rows, err := d.conn.Query(d.ctx, stmt, tradingPairSymbol, limit)
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
		INSERT INTO receipts (symbol, type, quantity, profit, price, timestamp)
		VALUES (@symbol, @type, @quantity, @profit, @price, @timestamp);`

	args := pgx.NamedArgs{
		"symbol":    receipt.Symbol,
		"type":      receipt.TxType,
		"quantity":  receipt.Quantity,
		"profit":    receipt.Profit,
		"price":     receipt.Price,
		"timestamp": receipt.Timestamp,
	}

	_, err = d.conn.Exec(d.ctx, stmt, args)

	return
}
