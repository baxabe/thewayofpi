package math

// https://play.golang.org/p/hF9jklt5lp

import (
	"math/big"
)

func Arccot(x int64, unity *big.Int) *big.Int {
	bigx := big.NewInt(x)
	xsquared := big.NewInt(x*x)
	sum := big.NewInt(0)
	sum.Div(unity, bigx)
	xpower := big.NewInt(0)
	xpower.Set(sum)
	n := int64(3)
	zero := big.NewInt(0)
	sign := false

	term := big.NewInt(0)
	for {
		xpower.Div(xpower, xsquared)
		term.Div(xpower, big.NewInt(n))
		if term.Cmp(zero) == 0 {
			break
		}
		if sign {
			sum.Add(sum, term)
		} else {
			sum.Sub(sum, term)
		}
		sign = !sign
		n += 2
	}
	return sum
}
