package config

import (
	"fmt"
	"strings"
)

func insertNewField[V int | bool | string](prompt, key string, m map[string]V) {
	var val V
	fmt.Print(prompt)
	fmt.Scanln(&val)
	m[key] = val
}

func New() string {
	var fieldsStr = map[string]string{}
	var fieldsInt = map[string]int{}
	var fieldsBool = map[string]bool{}

	var promptStr = map[string]string{
		"enter base currency: ":                           "base",
		"enter quote currency: ":                          "quote",
		"enter trading balance: ":                         "tradingBalance",
		"enter percent change buy: ":                      "percentChangeBuy",
		"enter percent change sell: ":                     "percentChangeSell",
		"enter server base endpoint: ":                    "serverBaseEndpoint",
		"enter websocket server base endpoint: ":          "wsServerBaseEndpoint",
		"enter type of first transaction (buy or sell): ": "firstTx",
	}

	var promptInt = map[string]string{
		"enter decision interval in milliseconds: ": "decisionIntervalMs",
	}

	/*
		var promptBool = map[string]string{
			"trade in testnet (true or false): ": "testnet",
		}
	*/

	for k, v := range promptStr {
		insertNewField(
			k,
			v,
			fieldsStr,
		)
	}

	for k, v := range promptInt {
		insertNewField(
			k,
			v,
			fieldsInt,
		)
	}

	/*
		for k, v := range promptBool {
			insertNewField(
				k,
				v,
				fieldsBool,
			)
		}
	*/

	return makeJson(fieldsStr, fieldsBool, fieldsInt)
}

func makeJson(fieldsStr map[string]string, fieldsBool map[string]bool, fieldsInt map[string]int) string {
	jsonStr := "{\n"

	for k, v := range fieldsStr {
		jsonStr += fmt.Sprintf("\"%s\": \"%s\",\n", k, v)
	}
	for k, v := range fieldsBool {
		jsonStr += fmt.Sprintf("\"%s\": \"%t\",\n", k, v)
	}
	for k, v := range fieldsInt {
		jsonStr += fmt.Sprintf("\"%s\": %d,\n", k, v)
	}

	jsonStr = strings.TrimRight(jsonStr, ",\n")
	jsonStr += "\n}"

	return jsonStr
}
