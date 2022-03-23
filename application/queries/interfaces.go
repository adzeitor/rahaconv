package queries

import (
	"github.com/adzeitor/rahaconv/domain/money"
	"github.com/adzeitor/rahaconv/internal/bigdecimal"
)

type CoinMarketRates interface {
	FetchExchangeRate(from money.Currency, to money.Currency) (*bigdecimal.BigDecimal, error)
}
