package coin_market_rates

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adzeitor/rahaconv/internal/bigdecimal"
)

func TestCoinMarketRates_CoinMarketRates(t *testing.T) {
	t.Run("success return exchange rate", func(t *testing.T) {
		response := `
			{
				"status": {
					"timestamp": "2022-03-23T22:42:09.471Z",
					"error_code": 0,
					"error_message": null,
					"elapsed": 0,
					"credit_count": 1,
					"notice": null
				},
				"data": {
					"BTC": {
						"id": 1326,
						"name": "5qvblo09mg9",
						"symbol": "BTC",
						"slug": "aqq90isinyr",
						"is_active": 6421,
						"is_fiat": null,
						"circulating_supply": 1780,
						"total_supply": 8473,
						"max_supply": 5978,
						"date_added": "2022-03-23T22:42:09.471Z",
						"num_market_pairs": 6269,
						"cmc_rank": 5100,
						"last_updated": "2022-03-23T22:42:09.471Z",
						"tags": [
							"7kjo7g8nwkk",
							"rp6ghuj4u3o",
							"kjwvmo0nf6l",
							"54uzi75e8v6",
							"si2wxpu697o",
							"4iw2qfzgupf",
							"9q65lbd1zeu",
							"5i1wcmu2uqw",
							"45qtsbc6318",
							"sdm847vatcf"
						],
						"platform": null,
						"self_reported_circulating_supply": null,
						"self_reported_market_cap": null,
						"quote": {
							"EUR": {
								"price": 0.47246190005249633,
								"volume_24h": 0.027082409623853554,
								"volume_change_24h": 0.6264242968772216,
								"percent_change_1h": 0.7397633920536006,
								"percent_change_24h": 0.38842564743327834,
								"percent_change_7d": 0.25178153742002185,
								"percent_change_30d": 0.1892082781795954,
								"market_cap": 0.5237368298850309,
								"market_cap_dominance": 9101,
								"fully_diluted_market_cap": 0.47683025383999356,
								"last_updated": "2022-03-23T22:42:09.471Z"
							}
						}
					}
				}
			}
		`

		var wantRequestURI string
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wantRequestURI = r.RequestURI
			w.Write([]byte(response))
		}))

		rates := NewCoinMarketRates("apikey")
		rates.Host = testServer.URL

		// act
		gotRate, err := rates.FetchExchangeRate("BTC", "EUR")
		require.NoError(t, err)

		// assert
		assert.Equal(t, "/v1/cryptocurrency/quotes/latest?convert=EUR&symbol=BTC", wantRequestURI)
		assert.Equal(t, bigdecimal.MustFromString("0.47246190005249633"), gotRate)
	})

	t.Run("returns error on error status", func(t *testing.T) {
		response := `
			{
				"status": {
					"timestamp": "2022-03-23T19:36:57.223Z",
					"error_code": 1001,
					"error_message": "This API Key is invalid.",
					"elapsed": 0,
					"credit_count": 0
				}
			}
		`
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(response))
		}))

		rates := NewCoinMarketRates("apikey")
		rates.Host = testServer.URL

		// act
		_, err := rates.FetchExchangeRate("BTC", "USD")

		// assert
		assert.EqualError(t, err, "coinmarketcap returns error: This API Key is invalid.")
	})
}
