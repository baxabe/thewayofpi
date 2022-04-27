package math

// https://play.golang.org/p/hF9jklt5lp

import (
	"fmt"
	"math/big"
)

func PiDigits(d int64) string {
	digits := big.NewInt(d + 10)
	unity := big.NewInt(0)
	unity.Exp(big.NewInt(10), digits, nil)
	pi := big.NewInt(0)
	four := big.NewInt(4)
	pi.Mul(four, pi.Sub(pi.Mul(four, Arccot(5, unity)), Arccot(239, unity)))
	return fmt.Sprintln(pi)
}