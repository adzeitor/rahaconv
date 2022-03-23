package money

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adzeitor/rahaconv/internal/bigdecimal"
)

func TestMoney_Convert(t *testing.T) {
	t.Run("simple conversion", func(t *testing.T) {
		// arrange
		inBtc, _ := NewFromString("0.00268053", "BTC")

		// act
		inDollars := inBtc.Convert("USD", bigdecimal.MustFromString("42969.7"))

		// assert
		want, _ := NewFromString("115.181569941", "USD")
		assert.Equal(t, want, inDollars)
	})
}
