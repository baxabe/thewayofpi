package main

import (
	"time"
	"math/rand"
)

var (
	rnd			*rand.Rand
)

func initRnd() {
	t := time.Date(1969, time.January, 1, 0, 0, 0, 0, time.UTC)
	d := time.Since(t)
	s := rand.NewSource(d.Nanoseconds())
	rnd = rand.New(s)
}
