package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func FormatValue(value *big.Int) string {
	if value != nil {
		m := new(big.Int)
		b, _ := new(big.Int).DivMod(value, big.NewInt(params.Ether), m)

		m = m.Div(m, big.NewInt(params.GWei))

		return b.String() + "." + m.String()
	}
	return "0.00"
}
