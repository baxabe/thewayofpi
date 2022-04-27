package main

import (
	"fmt"
)

import (
	"thewayofpi.com/buffer/balancer/balances"
	"thewayofpi.com/buffer/logs"
)

var dimen = [32]uint8{
	0, 1, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3,
}

var levels = [32]uint8{
	0, 1, 2, 2, 3, 3, 4, 4,
	2, 2, 2, 2, 3, 3, 3, 3,
	4, 4, 4, 4, 2, 2, 2, 2,
	2, 2, 2, 2, 3, 3, 3, 3,
}

var vals = [32]uint8{
	0, 2, 3, 4, 5, 6, 7, 8,
	9, 9, 9, 9, 10, 10, 10, 10,
	11, 11, 11, 11, 12, 12, 12, 12,
	13, 13, 13, 13, 14, 14, 14, 14,
}

func main() {
	for infodim := 1; infodim < 32; infodim++ {
		for count := 0; count < 10; count++ {
			try(infodim, count)
			fmt.Println()
		}
	}
}

func try(infodim, count int) {
	dim := dimen[infodim]
	leaves := vals[infodim]
	deep := levels[infodim]
	fmt.Println()
	fmt.Printf("Infodim: %v - Try: %v", infodim, count)
	fmt.Printf("Parameters - dim: %v, leaves: %v, deep: %v", dim, leaves, deep)
	b, err := balances.New(dim, leaves, deep)
	if err != nil {
		fmt.Println(err)
	}
	sol := b.Solution()
	fmt.Printf("Balance...........: %v", b)
	fmt.Printf("Solution..........: %v", sol)
	fmt.Printf("Structure.........: %v", b.Structure())
	fmt.Printf("Canonical.........: %v", b.Canonical())
	if ev, err := b.Eval(sol); err != nil {
		logs.Error(fmt.Errorf("testbal: Eval(sol) got error - %v", err))
	} else {
		fmt.Printf("Eval(sol).........: %v", ev)
	}
	if ev, err := b.Eval(reverse(sol)); err != nil {
		logs.Error(fmt.Errorf("testbal: Eval(reverse(sol)) got error - %v", err))
	} else {
		fmt.Printf("Eval(reverse(sol)): %v", ev)
	}
	if js, err := b.JSON(sol); err != nil {
		logs.Error(fmt.Errorf("testbal: JSON(sol) got error - %v", err))
	} else {
		fmt.Printf(fmt.Sprintf("JSON(sol).........: %s", js))
	}
}

func reverse(slice []balances.Weight) []balances.Weight {
	ls := len(slice)
	for i := 0; i < ls/2; i++ {
		slice[i], slice[ls-1-i] = slice[ls-1-i], slice[i]
	}
	return slice
}
