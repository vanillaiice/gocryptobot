package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/joho/godotenv"
	"github.com/vanillaiice/gocryptobot/strategy/simple"
)

const perm = 0644

func main() {
	var configPath, strategy string
	flag.StringVar(&strategy, "s", "simple", "trading strategy to use (simple)")
	flag.Parse()
	configPath = flag.Arg(0)
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		var createFileAns string
		var secretKey, apiKey string
		fmt.Print("no .env file found, create one ? (y/N): ")
		fmt.Scanln(&createFileAns)
		if strings.ToLower(createFileAns) == "no" || strings.ToLower(createFileAns) == "n" {
			fmt.Println("not creating .env file")
			os.Exit(0)
		}
		fmt.Print("please enter your binance secret key: ")
		fmt.Scanln(&secretKey)
		fmt.Print("please enter your binance api key: ")
		fmt.Scanln(&apiKey)
		err = os.WriteFile(".env", []byte(fmt.Sprintf("SECRET_KEY = %s\nAPI_KEY = %s", secretKey, apiKey)), perm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	var keysMap map[string]string
	keys := []string{"SECRET_KEY", "API_KEY"}
	keysMap, err := godotenv.Read()
	if err != nil {
		log.Fatalln(err)
	}
	for _, k := range keys {
		if keysMap[k] == "" {
			log.Fatalf("%s not found in .env file\n", k)
		}
	}

	strategies := []string{"simple"}
	if slices.Contains(strategies, strategy) == false {
		log.Fatalf("strategy %s not valid\n", strategy)
		os.Exit(1)
	}

	switch strategy {
	case "simple":
		simple.Start(configPath, keysMap)
	default:
		simple.Start(configPath, keysMap)
	}
}
