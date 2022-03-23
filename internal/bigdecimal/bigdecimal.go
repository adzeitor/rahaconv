package bigdecimal

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var zero = big.NewInt(0)
var ten = big.NewInt(10)

// BigDecimal is an arbitrary precision decimal numbers "suitable" (i think) for finance.
type BigDecimal struct {
	Value *big.Int
	// Number of digits after decimal point.
	Exponent int64
}

// FIXME: pointers is needed here? Maybe immutable variant, but big.Int with pointers...
func FromString(s string) (*BigDecimal, error) {
	number := &BigDecimal{
		Value:    &big.Int{},
		Exponent: 0,
	}
	intPart, fracPart, isCutted := strings.Cut(s, ".")
	if isCutted {
		number.Exponent = int64(len(fracPart))
	}

	_, ok := number.Value.SetString(intPart+fracPart, 10)
	if !ok {
		return nil, errors.New("not a number")
	}
	return number.Normalize(), nil
}

func MustFromString(s string) *BigDecimal {
	n, err := FromString(s)
	if err != nil {
		panic(err)
	}
	return n
}

// Normalize big decimal number.
// For example, 1234500 with exponent=2 equals to 12345 exponent=0.
func (n *BigDecimal) Normalize() *BigDecimal {
	normalized := (&big.Int{}).Set(n.Value)
	remainder := &big.Int{}
	for n.Exponent > 0 {
		normalized.QuoRem(normalized, ten, remainder)
		if remainder.Cmp(zero) != 0 {
			return n
		}
		n.Value.Set(normalized)
		n.Exponent--
	}
	return n
}

func (n BigDecimal) String() string {
	s := n.Value.String()
	if n.Exponent == 0 {
		return s
	}
	remainder := &big.Int{}
	power := &big.Int{}
	power.SetInt64(10).Exp(power, big.NewInt(n.Exponent), remainder)

	intPart := &big.Int{}
	intPart.Set(n.Value)
	intPart.QuoRem(intPart, power, remainder)

	// FIXME: OMG. Use leftpad.
	remainder.Abs(remainder)
	fracPart := fmt.Sprintf("%0*s", n.Exponent, remainder.String())
	return intPart.String() + "." + fracPart
}

// FIXME: implement UnmarshalText for string representation in json?
func (n *BigDecimal) UnmarshalJSON(data []byte) error {
	parsed, err := FromString(string(data))
	if err != nil {
		return err
	}
	n.Value = parsed.Value
	n.Exponent = parsed.Exponent
	return nil
}

func (n *BigDecimal) Mul(other *BigDecimal) *BigDecimal {
	result := &BigDecimal{
		Value: &big.Int{},
	}
	// x/(10**n) * y/(10**k) == x*y / (10**(n+k))
	result.Value = result.Value.Mul(n.Value, other.Value)
	result.Exponent = n.Exponent + other.Exponent
	return result.Normalize()
}
