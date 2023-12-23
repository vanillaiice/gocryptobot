package gobinance

type Client struct {
	BaseURL   string
	SecretKey string
	ApiKey    string
}

type Side string
type OrderType string

const SideTypeBuy Side = "BUY"
const SideTypeSell Side = "SELL"
const OrderTypeMarket OrderType = "MARKET"
const OrderTypeLimit OrderType = "LIMIT"

func NewClient(BaseURL, SecretKey, ApiKey string) *Client {
	return &Client{BaseURL, SecretKey, ApiKey}
}

func (c *Client) Time() (int64, error) {
	return time(c.BaseURL)
}

func (c *Client) AccountInfo() (ResponseAccount, error) {
	return info(c.BaseURL, c.SecretKey, c.ApiKey)
}

func (c *Client) AccountInfoPretty() (string, error) {
	return infoPretty(c.BaseURL, c.SecretKey, c.ApiKey)
}

func (c *Client) Order(side Side, orderType OrderType, symbol, quantity string) (ResponseOrder, error) {
	options := map[string]string{"side": string(side), "symbol": symbol, "quantity": quantity, "type": string(orderType)}
	return place(c.BaseURL, c.SecretKey, c.ApiKey, options)
}
