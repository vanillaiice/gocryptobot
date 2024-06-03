package binance

// Client is a client for the Binance API.
type Client struct {
	BaseURL   string
	SecretKey string
	ApiKey    string
}

// OrderSide is an enum for the order side.
type OrderSide string

// OrderSideBuy represents the buy side of an order.
const OrderSideBuy OrderSide = "BUY"

// OrderSideSell represents the sell side of an order.
const OrderSideSell OrderSide = "SELL"

// OrderType is an enum for the order type.
type OrderType string

// OrderTypeMarket represents a market order.
const OrderTypeMarket OrderType = "MARKET"

// OrderTypeLimit represents a limit order.
const OrderTypeLimit OrderType = "LIMIT"

// NewClient creates a new client for the Binance API.
func NewClient(BaseURL, SecretKey, ApiKey string) *Client {
	return &Client{
		BaseURL:   BaseURL,
		SecretKey: SecretKey,
		ApiKey:    ApiKey,
	}
}

// Time returns the server time.
func (c *Client) Time() (int64, error) {
	return time(c.BaseURL)
}

// AccountInfo returns the account information.
func (c *Client) AccountInfo() (*ResponseAccount, error) {
	return info(c)
}

// AccountInfoPretty returns the account information in a pretty format.
func (c *Client) AccountInfoPretty() (string, error) {
	return infoPretty(c)
}

// ExchangeInfo returns the exchange info.
func (c *Client) ExchangeInfo() (*ResponseExchangeInfo, error) {
	return exchangeInfo(c.BaseURL)
}

// StepSize returns the step size for the symbol(s).
func (c *Client) StepSize(symbols ...string) (map[string]string, error) {
	return stepSize(c.BaseURL, symbols...)
}

// Order places an order.
func (c *Client) Order(opts map[string]string) (*ResponseOrder, error) {
	return place(c, opts)
}
