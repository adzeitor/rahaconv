package money

import (
	"github.com/adzeitor/rahaconv/internal/bigdecimal"
)

type Money struct {
	// FIXME: just use some library for decimal from github, but NIH syndrome and just for fun ;)
	Amount   *bigdecimal.BigDecimal
	Currency Currency
	// FIXME: add scale (max digits after decimal point)?
}

func New(amount *bigdecimal.BigDecimal, currency Currency) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
	}
}

func NewFromString(amount string, currency Currency) (Money, error) {
	n, err := bigdecimal.FromString(amount)
	if err != nil {
		return Money{}, nil
	}
	return New(n, currency), nil
}

func (n Money) Convert(newCurrency Currency, rate *bigdecimal.BigDecimal) Money {
	return New(n.Amount.Mul(rate), newCurrency)
}

func (n Money) String() string {
	// TODO: number of digits after decimal point depenends on currency so we
	// should treat it here or add format method.
	return n.Amount.String() + " " + n.Currency
}
