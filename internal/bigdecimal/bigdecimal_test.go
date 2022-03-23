package bigdecimal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		value   string
		want    string
		wantErr bool
	}{
		{
			value: "0",
			want:  "0",
		},
		{
			value: "00000",
			want:  "0",
		},
		{
			value: "00000.000",
			want:  "0",
		},
		{
			value: "0.000000001",
			want:  "0.000000001",
		},
		{
			value: "12345",
			want:  "12345",
		},
		{
			value: "-12345",
			want:  "-12345",
		},
		{
			value: "12345.00",
			want:  "12345",
		},
		{
			value: "00012345000.000000",
			want:  "12345000",
		},
		{
			value: "0012300.00123",
			want:  "12300.00123",
		},
		{
			value: "0012300.0099900",
			want:  "12300.00999",
		},
		{
			value: "10000.11",
			want:  "10000.11",
		},
		{
			value: "1238483214.482348384",
			want:  "1238483214.482348384",
		},
		{
			value:   "1238483214.482348384.5",
			wantErr: true,
		},
		{
			value:   "",
			wantErr: true,
		},
		{
			value:   "foobar",
			wantErr: true,
		},
		{
			value:   "abc.def",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got, err := FromString(tt.value)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.want, got.String())
			}
		})
	}
}

func TestBigDecimal_Mul(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		want string
	}{
		{
			a:    "0.1",
			b:    "0.2",
			want: "0.02",
		},
		{
			a:    "0.00268605",
			b:    "42894.48",
			want: "115.216718004",
		},
		{
			a:    "-563.2344",
			b:    "42894.48",
			want: "-24159646.706112",
		},
		{
			a:    "-0.0000001",
			b:    "-99999.99",
			want: "0.009999999",
		},
	}
	for _, tt := range tests {
		t.Run(tt.a+"*"+tt.b, func(t *testing.T) {
			got := MustFromString(tt.a).Mul(MustFromString(tt.b))
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestBigDecimal_UnmarshalText(t *testing.T) {
	type structWithDecimal struct {
		Amount *BigDecimal
	}

	t.Run("unmarshal to struct from json", func(t *testing.T) {
		data := `
			{
				"Amount": 12345.12
			}
		`
		var got structWithDecimal
		err := json.Unmarshal([]byte(data), &got)
		assert.NoError(t, err)

		assert.Equal(t, "12345.12", got.Amount.String())
	})

	t.Run("return error on non decimal", func(t *testing.T) {
		data := `
			{
				"Amount": ""
			}
		`
		var got structWithDecimal
		err := json.Unmarshal([]byte(data), &got)
		assert.Error(t, err)
	})
}
