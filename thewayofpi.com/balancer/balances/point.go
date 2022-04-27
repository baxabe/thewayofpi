package balances

import (
	"errors"
	"fmt"
	"math"
)

type (
	Weight int
	Delta  int
)

type (
	point struct {
		weight   Weight
		location []Delta
	}
)

func newPoint() point {
	return point{Weight(rnd.Intn(math.MaxUint8) + 1), nil}
}

func emptyPoint() point {
	return point{}
}

func (p point) insert(dim, deep byte) (*protobal, error) {
	if deep < 1 {
		return nil, errors.New("balances:Point.Insert - limit deep reached")
	}
	bal := newProtobal(dim, p, newPoint())
	return bal, nil
}

func (p point) force() Weight {
	return p.weight
}

func (p *point) place() []Delta {
	return p.location
}

func (p *point) setPlace(l []Delta) {
	p.location = l
}

func (p point) canonize() []token {
	return []token{token{POINT, p}}
}

func (p point) structure() string {
	return fmt.Sprintf(".")
}

func (p point) ponder() {}

func (p point) Eval(w []Weight) []result {
	return []result{UNKNOWN}
}

func (p point) solve() []Weight {
	return []Weight{p.weight}
}

func (p point) String() string {
	return fmt.Sprintf("<%v@%v>", p.weight, p.location)
}
