// count
package main

import (
	"fmt"
	"math"
)

const (
	maxLevel = 60
)
var (
	ntot [maxLevel]int
	nodos [maxLevel] int
	njug [maxLevel]int
	ntotxjug [maxLevel]int
)

func exp(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func main() {
	ntot[0] = 1
	ntot[1] = 3
	nodos[0] = 1
	nodos[1] = 2
	njug[0] = 1
	njug[1] = 1
	ntotxjug[0] = 1
	ntotxjug[1] = 3
	fmt.Println("Arbol")
	fmt.Printf("%3s\t%16s\t%16s\t%16s\t%16s\n", "Dim", "Nod", "Tot", "Jug", "NxJ")
	for n := 0; n < maxLevel; n++ {
		if n > 1 {
			m := n % 5
			switch m {
				case 2, 3, 4:
					ntot[n] = exp(2, n) + 2*ntot[n-1]
					nodos[n] = 1 + 2*nodos[n-1]
					njug[n] = 2*njug[n-1]
				case 1:
					ntot[n] = exp(2, n) + ntot[n-1]
					nodos[n] = 1 + nodos[n-1]
					njug[n] = njug[n-1]
				case 0:
					ntot[n] = exp(2, n) + 3*ntot[n-1]
					nodos[n] = 1 + 3*nodos[n-1]
					njug[n] = 3*njug[n-1]
			}
			ntotxjug[n] = ntot[n]/njug[n]
		}
		fmt.Printf("%3d\t%16d\t%16d\t%16d\t%16d\n", n, nodos[n], ntot[n], njug[n], ntotxjug[n])
	}
}
