package random

import (
	"math/rand"
	"time"
)

func New() *rand.Rand {
	t := time.Date(1969, time.January, 1, 0, 0, 0, 0, time.UTC)
	d := time.Since(t)
	s := rand.NewSource(d.Nanoseconds())
	return rand.New(s)
}

