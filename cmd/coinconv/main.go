package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/adzeitor/rahaconv/application/queries"
	"github.com/adzeitor/rahaconv/infrastructure/clients/coin_market_rates"
	"github.com/adzeitor/rahaconv/internal/bigdecimal"
)

func printHelp() {
	fmt.Printf("use: \n" +
		"COINMARKET_API_KEY='b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c' coinconv 117.6 USD BTC\n")
}

// FIXME: Move this code to infrastructure?
func convertCurrencyQueryController(
	handler *queries.ConvertCurrencyQueryHandler,
) error {
	amount, err := bigdecimal.FromString(os.Args[1])
	if err != nil {
		return errors.New("amount should be in decimal format like 111.5")
	}
	query := queries.ConvertCurrencyQuery{
		Amount:       amount,
		FromCurrency: os.Args[2],
		ToCurrency:   os.Args[3],
	}

	result, err := handler.Handle(query)
	if err != nil {
		return fmt.Errorf("conversion error: %w", err)
	}

	fmt.Println(result.Amount)
	return nil
}

func main() {
	if len(os.Args) < 3 {
		printHelp()
		os.Exit(1)
	}
	coinMarketApiKey := os.Getenv("COINMARKET_API_KEY")
	if coinMarketApiKey == "" {
		printHelp()
		os.Exit(2)
	}

	coinMaketRates := coin_market_rates.NewCoinMarketRates(coinMarketApiKey)
	coinMarketHost := os.Getenv("COINMARKET_API_HOST")
	if coinMarketHost != "" {
		coinMaketRates.Host = coinMarketHost
	}

	convertCurrencyQueryHandler := queries.NewConvertCurrencyQueryHandler(coinMaketRates)

	err := convertCurrencyQueryController(convertCurrencyQueryHandler)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}
