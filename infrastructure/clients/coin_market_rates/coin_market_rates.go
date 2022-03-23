package coin_market_rates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/adzeitor/rahaconv/domain/money"
	"github.com/adzeitor/rahaconv/internal/bigdecimal"
)

const (
	defaultURI     = "https://pro-api.coinmarketcap.com"
	defaultTimeout = 10 * time.Second
)

type CoinMarketRates struct {
	Client http.Client
	ApiKey string
	Host   string
}

func NewCoinMarketRates(apiKey string) *CoinMarketRates {
	client := &CoinMarketRates{
		Client: http.Client{
			Timeout: defaultTimeout,
		},
		Host:   defaultURI,
		ApiKey: apiKey,
	}
	return client
}

type coinMarketQuote struct {
	Symbol string
	Price  *bigdecimal.BigDecimal
}

type coinMarketCurrency struct {
	Symbol string
	Quotes map[string]coinMarketQuote `json:"quote"`
}

type status struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type coinMarketQuotesResponse struct {
	Status status
	Data   map[string]coinMarketCurrency
}

// FIXME: improve readability, long method.
func (client *CoinMarketRates) FetchExchangeRate(
	from money.Currency, to money.Currency,
) (*bigdecimal.BigDecimal, error) {
	uri, err := url.Parse(client.Host + "/v1/cryptocurrency/quotes/latest")
	if err != nil {
		return nil, err
	}
	// FIXME: In cryptocurrency calls you would then send, for example id=1027, instead of symbol=ETH
	// It's strongly recommended that any production code utilize these IDs for cryptocurrencies,
	// exchanges, and markets to future-proof your code.
	// https://coinmarketcap.com/api/v1/#section/Best-Practices
	query := url.Values{
		"symbol":  []string{from},
		"convert": []string{to},
	}
	uri.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-CMC_PRO_API_KEY", client.ApiKey)

	resp, err := client.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response coinMarketQuotesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if response.Status.ErrorCode != 0 {
		err = fmt.Errorf("coinmarketcap returns error: %s", response.Status.ErrorMessage)
		return nil, err
	}

	currency, ok := response.Data[from]
	if !ok {
		return nil, fmt.Errorf("conversion from %s is not available", from)
	}
	quote, ok := currency.Quotes[to]
	if !ok {
		return nil, fmt.Errorf("conversion from %s to %s is not available", from, to)
	}
	return quote.Price, nil
}
