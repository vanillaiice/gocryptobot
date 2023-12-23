# Simple Go Crypto Bot

Simple Crypto Trading Bot using Binance Spot and Websocket API.
This Bot buys and sells crypto at user specified margins.
This is basically a rewrite of 
[V Crypto Bot](https://github.com/vanillaiice/vcryptobot) in Go.

# Installation

```go
$ go install github.com/vanillaiice/gocryptobot/cmd/gocryptobot@latest
# or
$ git clone github.com/vanillaiice/gocryptobot
$ cd gocryptobot
$ make
```

# Usage

```sh
# Example Usage
$ gocryptobot [path to config file]
```

- If the path to the config file is omitted, a prompt asking you if you
want to create one will be shown.

- Also, you can create a .env file with your Binance Secret and Api keys
in the following format:

```
SECRET_KEY = "<YOUR SECRET KEY>"
API_KEY = "<YOUR API KEY>"
```

- If the .env file is not present when running the program, a prompt asking you
if you want to create one will be shown.

# Config File

- ```base```, base currency of the trading pair.
> example: "BTC"

- ```quote```, quote currency of the trading pair.
> example: "USDT"

- ```tradingBalance```, initial trading balance of the bot.
> example for BTC: "0.01

- ```percentChangeSell```, percent change between the current price and last buy price, at which the bot will sell.
> example value (%): "5.0"

- ```percentChangeBuy```, percent change between the last sell price and current price, at which the bot will buy.
> example value (%): "5.0"

- ```firstTx```, type of the first transaction executed by the bot.
> accepted values: "buy" or "sell"

- ```decisionIntervalMs```, time in milliseconds at which the bot will decide to buy or sell.
> example value in ms: 1500

- ```serverBaseEndpoint```, base endpoint of the Binance API server.
> example base endpoint for testnet server: "testnet.binance.vision"

> example base endpoint for normal server: "api.binance.com"

- ```wsServerBaseEndpoint```, base endpoint of the Binance Websocket API server.
> example base endpoint for normal server: "api.binance.com"

- ```logDB```, if the bot should log the transaction receipts in a sqlite database.
> accepted values: true or false

## Sample Config File

```json
{
"base": "BTC",
"quote": "USDT",
"tradingBalance": "0.1",
"decisionIntervalMs": 1500,
"firstTx": "buy",
"percentChangeSell": "3.0",
"percentChangeBuy": "1.0",
"logDB": true,
"serverBaseEndpoint": "testnet.binance.vision",
"wsServerBaseEndpoint": "testnet.binance.vision"
}
```

# Dev Dependencies

- [sqlite](https://modules.vlang.io/db.sqlite.html)
- make (optional)

# Additional tools

To view the data in the sqlite databases (transaction receipts), 
you can install [DB Browser for sqlite](https://sqlitebrowser.org/dl/).

# Disclaimer

- No warranty whatsoever, use at your own risk
- Trading crypto is very risky, *only invest in what you can afford to lose*

# Author

vanillaiice

# License

GPLv3
