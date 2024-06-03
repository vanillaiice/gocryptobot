package db

// Receipt stores a receipt's data.
type Receipt struct {
	Symbol    string
	TxType    string
	Quantity  float32
	Price     float32
	Profit    float32
	Timestamp int64
}

// DB is a database interface.
type DB interface {
	GetLastReceipts(limit int) ([]*Receipt, error)
	InsertReceipt(receipt *Receipt) error
	Close() error
}
