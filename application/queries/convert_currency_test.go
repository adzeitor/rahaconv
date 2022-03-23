package queries

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adzeitor/rahaconv/domain/money"
	"github.com/adzeitor/rahaconv/infrastructure/clients/coin_market_rates"
	"github.com/adzeitor/rahaconv/internal/bigdecimal"
)

func TestConvertCurrencyQueryHandler_Handle(t *testing.T) {
	t.Run("success conversion", func(t *testing.T) {
		// arrange
		fakeRates := coin_market_rates.NewCoinMarketRatesFake()
		fakeRates.AddExchange("USD", "BTC", "0.0000235924")
		handler := NewConvertCurrencyQueryHandler(fakeRates)

		// act
		query := ConvertCurrencyQuery{
			Amount:       bigdecimal.MustFromString("1012.23"),
			FromCurrency: "USD",
			ToCurrency:   "BTC",
		}
		result, err := handler.Handle(query)
		require.NoError(t, err)

		// assert
		want, _ := money.NewFromString("0.023880935052", "BTC")
		assert.Equal(t, want, result.Amount)
	})

	t.Run("returns error if coin market return error", func(t *testing.T) {
		// arrange
		fakeRates := coin_market_rates.NewCoinMarketRatesFake()
		fakeRates.Err = errors.New("unavailable")
		handler := NewConvertCurrencyQueryHandler(fakeRates)

		// act
		query := ConvertCurrencyQuery{
			Amount:       bigdecimal.MustFromString("1012.23"),
			FromCurrency: "USD",
			ToCurrency:   "BTC",
		}
		_, err := handler.Handle(query)

		// assert
		assert.Error(t, err)
	})
}
