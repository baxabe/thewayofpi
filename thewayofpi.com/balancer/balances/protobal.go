package balances

import (
	"errors"
	"fmt"
	"math"
)

type hanger interface {
	insert(byte, byte) (*protobal, error)
	force() Weight
	place() []Delta
	setPlace([]Delta)
	canonize() []token
	structure() string
	ponder()
	solve() []Weight
}

type (
	protobal struct {
		dimension byte
		hang      point
		children  []hanger
	}
)

func (b protobal) insert(dim, deep byte) (*protobal, error) {
	var err error
	lchildren := len(b.children)
	if deep > 0 && rnd.Intn(100) > 40 {
		if err = b.dig(dim, deep); err == nil {
			return &b, nil
		}
	}
	if lchildren < int(math.Pow(2, float64(b.dimension))) {
		p := newPoint()
		b.children = append(b.children, &p)
		b.hang.weight += p.weight
		return &b, nil
	}
	if deep > 0 {
		if err = b.dig(dim, deep); err == nil {
			return &b, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("balances:Balance.Insert - Point cann't be allocated.\n\t%v", err))
}

func (b *protobal) dig(dim, deep byte) error {
	var d byte = dim
	if d < 1 {
		d = 1
	}
	lchildren := len(b.children)
	perm := rnd.Perm(lchildren)
	for i := 0; i < lchildren; i++ {
		if bal, err := b.children[perm[i]].insert(d, deep-1); err == nil {
			b.children[perm[i]] = bal
			b.hang.weight = 0
			for _, c := range b.children {
				b.hang.weight += c.force()
			}
			return nil
		}
	}
	return errors.New(fmt.Sprintf("balances:Balance.dig - failed [dim: %v, deep: %v]", dim, deep))
}

func (b protobal) force() Weight {
	return b.hang.weight
}

func (b *protobal) place() []Delta {
	return b.hang.location
}

func (b *protobal) setPlace(place []Delta) {
	b.hang.location = place
}

func (b protobal) canonize() []token {
	result := []token{token{OPEN, point{b.hang.weight, b.hang.location}}}
	for _, c := range b.children {
		result = append(result, c.canonize()...)
	}
	result = append(result, token{CLOSE, point{b.hang.weight, b.hang.location}})
	return result
}

func (b protobal) structure() string {
	result := fmt.Sprintf("(")
	if len(b.children) > 0 {
		for _, c := range b.children {
			result = result + c.structure()
		}
	}
	result = result + fmt.Sprintf(")")
	return result
}

func (b protobal) ponder() {
	for _, c := range b.children {
		c.ponder()
	}
	dim := int(b.dimension)
	lchildren := len(b.children)
	if lchildren < 2 {
		return
	}
	b.children[0].setPlace(make([]Delta, dim))
	for i := 1; i < lchildren; i++ {
		b.children[i].setPlace(make([]Delta, dim))
		for j := 0; j < dim; j++ {
			b.children[i].place()[j] = Delta(math.Pow(-1, float64(rnd.Intn(2)))) * Delta(b.children[0].force()* Weight(rnd.Intn(math.MaxUint8+1)))
		}
	}
	for i := 0; i < dim; i++ {
		k := Delta(0)
		for j := 1; j < lchildren; j++ {
			k += -1 * Delta(b.children[j].force()) * b.children[j].place()[i]
		}
		b.children[0].place()[i] = k / Delta(b.children[0].force())
	}
}

func (b protobal) solve() []Weight {
	var result []Weight
	for _, c := range b.children {
		result = append(result, c.solve()...)
	}
	return result
}

func (b protobal) String() string {
	result := fmt.Sprintf("{%v}%c", b.hang, OPEN)
	for _, c := range b.children {
		result = result + fmt.Sprintf("%v", c)
	}
	result = result + fmt.Sprintf("%c", CLOSE)
	return result
}

func newProtobal(dim byte, p0, p1 point) *protobal {
	hang := point{p0.weight + p1.weight, root(dim)}
	children := []hanger{&p0, &p1}
	var d byte = dim
	if d < 1 {
		d = 1
	}
	return &protobal{d, hang, children}
}

func root(dim byte) []Delta {
	result := make([]Delta, int(dim))
	for i := 0; i < int(dim); i++ {
		result[i] = 0
	}
	return result
}
