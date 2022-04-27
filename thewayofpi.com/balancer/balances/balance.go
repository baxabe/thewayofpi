package balances

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

import (
	"thewayofpi.com/buffer/random"
)

type (
	balance struct {
		canonical	[]token
		structure	string
		solution	[]Weight
	}
	evaluation	[]result
)

var (
	rnd *rand.Rand
)

func (b balance) Canonical() []token {
	return b.canonical
}

func (b balance) Structure() string {
	return b.structure
}

func (b balance) Eval(w []Weight) (evaluation, error) {
	if res, err := eval(b, w); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (b balance) Solution() []Weight {
	return b.solution
}

func (b balance) JSON(w []Weight) ([]byte, error) {
	ev, err := eval(b, w)
	if err != nil {
		return nil, err
	}
	result, err := json.Marshal(ev)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b balance) String() string {
	var result string
	for _, t := range b.canonical {
			result += fmt.Sprintf("%v", t)
	}
	return result
}

// dimension: Number of axis
// leaves: Total number of weights to place
// levels: Max deep of generated tree
func New(dim, leaves, deep byte) (*balance, error) {
	if dim < 1 {
		return nil, errors.New("balances:New - dim too small")
	}
	if leaves < 2 {
		return nil, errors.New("balances:New - too few leaves")
	}
	if deep < 1 {
		return nil, errors.New("balances:New - deep too low")
	}
	rnd = random.New()
	proto := newProtobal(dim, newPoint(), newPoint())
	for i := 0; i < int(leaves-2); i++ {
		if bal, err := proto.insert(dim, deep); err == nil {
			proto = bal
		} else {
			return nil, err
		}
	}
	proto.ponder()
	canonical:= proto.canonize()
	structure:= proto.structure()
	solution:= proto.solve()
	return &balance{canonical, structure, solution}, nil
}

//func (balance *Balance) Load(ps []datastore.Property) error {
//	return nil
//}

//func (balance *Balance) Save() ([]datastore.Property, error) {
//	return nil, nil
//}
