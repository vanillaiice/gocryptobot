# Simple Go Crypto Bot

Simple Crypto Trading Bot using Binance Spot and Websocket API.
This Bot buys and sells crypto at user specified margins.
This is basically a port of
[V Crypto Bot](https://github.com/vanillaiice/vcryptobot) in Go.

# Installation

```sh
$ go install github.com/vanillaiice/gocryptobot@latest
```

# Usage

```sh
Simple cryptocurrency trading bot using the Binance API and Websocket.

Usage:
  gocryptobot [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  simple      Use the simple trading strategy

Flags:
  -c, --config string   config file (default is $HOME/.gocryptobot.yaml)
  -h, --help            help for gocryptobot
  -t, --toggle          Help message for toggle
  -v, --version         version for gocryptobot

Use "gocryptobot [command] --help" for more information about a command.
```

- Also, an .env file with your Binance Secret and API key should be present in the following format:

```
SECRET_KEY = "<YOUR SECRET KEY>"
API_KEY = "<YOUR API KEY>"
```
## Sample Config File

```yaml
base: "BTC"
quote: "USDT"
trading-balance: "0.5"
state-file: "state.json"
first-tx: "buy"
percent-change-sell: "0.50"
percent-change-buy: "0.25"
server-base-endpoint: "testnet.binance.vision"
ws-server-base-endpoint: "testnet.binance.vision"
save-receipt: true
db-backend: "sqlite"
db-url: "receipts.db"
log-level: "info"
```

# Contributing

## bugs, issues, feature requests, etc.

Please fork the project, make your changes and submit a pull request.

## trading strategies

Please fork the project, create your own strategy under the `strategies` folder, and then create a new
command using cobra-cli:

```sh
$ cobra-cli add root -p <STRATEGY_NAME>
```

# Additional tools

## sqlite

To view the data in the sqlite databases (transaction receipts), 
you can install [DB Browser for sqlite](https://sqlitebrowser.org/dl/).

# Disclaimer

- No warranty whatsoever, use at your own risk
- Trading crypto is very risky, *only invest in what you can afford to lose*

# Author

vanillaiice

# License

GPLv3
