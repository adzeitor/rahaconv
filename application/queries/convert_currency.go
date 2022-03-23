package queries

import (
	"github.com/adzeitor/rahaconv/domain/money"
	"github.com/adzeitor/rahaconv/internal/bigdecimal"
)

type ConvertCurrencyQuery struct {
	Amount       *bigdecimal.BigDecimal
	FromCurrency money.Currency
	ToCurrency   money.Currency
}

type ConvertCurrencyResult struct {
	Amount money.Money
}

type ConvertCurrencyQueryHandler struct {
	Rates CoinMarketRates
}

func NewConvertCurrencyQueryHandler(rates CoinMarketRates) *ConvertCurrencyQueryHandler {
	return &ConvertCurrencyQueryHandler{
		Rates: rates,
	}
}

func (handler *ConvertCurrencyQueryHandler) Handle(
	query ConvertCurrencyQuery,
) (ConvertCurrencyResult, error) {
	rate, err := handler.Rates.FetchExchangeRate(query.FromCurrency, query.ToCurrency)
	if err != nil {
		return ConvertCurrencyResult{}, err
	}

	original := money.New(query.Amount, query.FromCurrency)
	converted := original.Convert(query.ToCurrency, rate)
	return ConvertCurrencyResult{Amount: converted}, nil
}
