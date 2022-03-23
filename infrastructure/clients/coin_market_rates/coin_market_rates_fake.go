package coin_market_rates

import (
	"github.com/adzeitor/rahaconv/domain/money"
	"github.com/adzeitor/rahaconv/internal/bigdecimal"
)

type CoinMarketRatesFake struct {
	Rates map[[2]string]*bigdecimal.BigDecimal
	Err   error
}

func NewCoinMarketRatesFake() *CoinMarketRatesFake {
	return &CoinMarketRatesFake{
		Rates: make(map[[2]string]*bigdecimal.BigDecimal),
	}
}

func (fake *CoinMarketRatesFake) AddExchange(
	from money.Currency, to money.Currency, rate string,
) {
	fake.Rates[[2]string{from, to}] = bigdecimal.MustFromString(rate)
}

func (fake *CoinMarketRatesFake) FetchExchangeRate(
	from money.Currency, to money.Currency,
) (*bigdecimal.BigDecimal, error) {
	rate := fake.Rates[[2]string{from, to}]
	return rate, fake.Err
}
