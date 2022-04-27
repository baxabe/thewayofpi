package balances

import (
	"container/list"
	"errors"
	"fmt"
)

import "github.com/hishboy/gocommons/lang"

type (
	result rune
)

const (
	UNKNOWN result = '\u003F' // '?'
	EQUAL   result = '\u003D' // '='
	LESS    result = '\u003C' // '<'
	GREAT   result = '\u003E' // '>'
	START   result = '\u0028' // '('
	END     result = '\u0029' // ')'
)

type (
	parsTree struct { // Parser tree
		mother   *parsTree
		output   evaluation
		children *list.List
		carrier  []int
		visited  int
		hook     point
	}

	parsHelp struct { // Parser helper
		current *parsTree
		root    *parsTree
	}
)

func eval(balance balance, weights []Weight) (evaluation, error) {
	if result, err := parse(balance, weights); err == nil {
		return result.root.children.Front().Value.(*parsTree).flat(), nil
	} else {
		return nil, errors.New(err.Error())
	}
}

func parse(balance balance, weights []Weight) (*parsHelp, error) {
	result := newParsHelp()
	rules := make([][]token, 8)
	rules[0] = []token{newToken(OPEN), newToken(NOTERMA), newToken(CLOSE)}
	rules[1] = []token{newToken(POINT), newToken(NOTERMB)}
	rules[2] = []token{newToken(NOTERMS), newToken(NOTERMB)}
	rules[3] = []token{newToken(POINT), newToken(NOTERMC)}
	rules[4] = []token{newToken(NOTERMS), newToken(NOTERMC)}
	rules[5] = []token{newToken(POINT), newToken(NOTERMC)}
	rules[6] = []token{newToken(NOTERMS), newToken(NOTERMC)}
	rules[7] = []token{newToken(EPSILON)}
	ttable := make([][]int, 4)
	ttable[0] = []int{-1, 0, -1, 0}
	ttable[1] = []int{1, 2, -1, -1}
	ttable[2] = []int{3, 4, -1, -1}
	ttable[3] = []int{5, 6, 7, -1}

	input := balance.canonical
	stack := lang.NewStack()
	stack.Push(newToken(EOF))
	stack.Push(newToken(NOTERMS))
	ip := 0
	item := stack.Peek()
	if item == nil {
		return nil, errors.New("balance:parse - stack is empty")
	}
	X, ok := item.(token)
	if !ok {
		return nil, errors.New("balance:parse - assertion error")
	}
	for !X.isEOF() {
		lookahead := input[ip]
		// fmt.Println(input[ip].Value)
		if X.equals(lookahead) {
			_ = stack.Pop()
			// Match terminal
			if err := matchTerminal(lookahead, &weights, &result); err != nil {
				return nil, errors.New(err.Error())
			}
			ip++
		} else if X.isTerminal() {
			// Error
			return nil, errors.New("balance:parse - X.IsTerminal")
		} else if x, y := lookahead.idx(), X.idy(); ttable[y][x] < 0 {
			// Error
			return nil, errors.New(fmt.Sprintf("balance:parse - ttable[y][x] < 0 - [y = %v - ip = %v - x = %v]", y, ip, x))
		} else {
			rule := rules[ttable[y][x]]
			// Output the production 'rule'
			_ = stack.Pop()
			if !rule[0].isEpsilon() {
				for i := len(rule) - 1; i >= 0; i-- {
					stack.Push(rule[i])
				}
			}
		}
		item = stack.Peek()
		if item == nil {
			return nil, errors.New("balance:parse - stack is empty")
		}
		if X, ok = item.(token); !ok {
			return nil, errors.New("balance:parse - assertion error")
		}
	}
	return &result, nil
}

func matchTerminal(token token, weights *[]Weight, helper *parsHelp) error {
	if !token.isTerminal() {
		return errors.New("balance:matchTerminal - token is not a terminal")
	}
	if helper == nil {
		return errors.New("balance:matchTerminal - helper is nil")
	}
	switch token.id {
	case OPEN:
		if err := matchOpen(point{0, token.value.location}, helper); err != nil {
			return err
		}
	case POINT:
		if *weights == nil || len(*weights) == 0 {
			return errors.New("balance:matchTerminal - not enough forces")
		}
		if err := matchPoint(point{(*weights)[0], token.value.location}, helper); err != nil {
			return err
		}
		*weights = (*weights)[1:]
	case CLOSE:
		if err := matchClose(helper); err != nil {
			return err
		}
	}
	return nil
}

func matchOpen(point point, helper *parsHelp) error {
	// fmt.Println(fmt.Sprintf("balance:matchOpen - |%v|(|", point))
	if err := helper.addChild(point); err != nil {
		return errors.New(fmt.Sprintf("balance:matchOpen - addChild got error: %v", err))
	}
	return nil
}

func matchPoint(point point, helper *parsHelp) error {
	// fmt.Println(fmt.Sprintf("balance:matchPoint - |%v|", point))
	if err := helper.addPoint(point); err != nil {
		return errors.New(fmt.Sprintf("balance:matchPoint - addPoint got error: %v", err))
	}
	return nil
}

func matchClose(helper *parsHelp) error {
	// fmt.Println("balance:matchClose - |)|")
	hook := helper.current.hook
	if err := helper.up(); err != nil {
		return errors.New(fmt.Sprintf("balance:matchClose - up got error: %v", err))
	}
	if err := helper.addPoint(hook); err != nil {
		return errors.New(fmt.Sprintf("balance:matchClose - addPoint got error: %v", err))
	}
	return nil
}

func newParsTree() *parsTree {
	children := list.New()
	return &parsTree{mother: nil, output: evaluation{}, children: children, carrier: []int{}, visited: 0, hook: emptyPoint()}
}

func (t parsTree) flat() evaluation {
	var result evaluation
	result = append(result, START)
	for _, v := range t.output {
		result = append(result, v)
	}
	for child := t.children.Front(); child != nil; child = child.Next() {
		result = append(result, child.Value.(*parsTree).flat()...)
	}
	result = append(result, END)
	return result
}

func newParsHelp() parsHelp {
	root := newParsTree()
	current := root
	return parsHelp{current, root}
}

func (h *parsHelp) up() error {
	if h == nil {
		return errors.New("balance:up - h is nil")
	}
	if h.current == nil {
		return errors.New("balance:up - h.current is nil")
	}
	if h.current.carrier == nil {
		return errors.New("balance:up - h.current.carrier is nil")
	}
	for _, v := range h.current.carrier {
		switch {
		case v < 0:
			h.current.output = append(h.current.output, LESS)
		case v > 0:
			h.current.output = append(h.current.output, GREAT)
		default:
			h.current.output = append(h.current.output, EQUAL)
		}
	}
	h.current = h.current.mother
	return nil
}

func (h *parsHelp) addChild(point point) error {
	if h == nil {
		return errors.New("balance:addChild - h is nil")
	}
	if h.current == nil {
		return errors.New("balance:addChild - h.current is nil")
	}
	child := newParsTree()
	child.mother = h.current
	h.current.children.PushBack(child)
	h.current = child
	h.current.hook = point
	return nil
}

func (h *parsHelp) addPoint(p point) error {
	if h == nil {
		return errors.New("balance:addPoint - h is nil")
	}
	if h.current == nil {
		return errors.New("balance:addPoint - h.current is nil")
	}
	if len(h.current.carrier) == 0 {
		h.current.carrier = make([]int, len(p.location))
	}
	if len(h.current.carrier) != len(p.location) {
		return errors.New("balance:addPoint - size mismatch")
	}
	h.current.hook.weight += p.weight
	for i, val := range p.location {
		h.current.carrier[i] += int(p.weight) * int(val)
	}
	return nil
}

func (r result) String() string {
	return fmt.Sprintf("%c", rune(r))
}

// Evaluation related
//
// Grammar:
// r0: S → (A)
// r1: A → pB
// r2: A → SB
// r3: B → pC
// r4: B → SC
// r5: C → pC
// r6: C → SC
// r7: C → ε
//
// Examples;
// (pp)
// (p(pp))
// (((ppp)(pp))pp(p(ppp)p))
//
// Sets:
// First('p') = {'p'}
// First('(') = {'('}
// First(')') = {')'}
// First(S) = {'('}		Follow(S) = {'p', '(', ')', 'ε', '$'}
// First(A) = {'p', '('}	Follow(A) = {')'}
// First(B) = {'p', '('}	Follow(B) = {')'}
// First(C) = {'p', '(', 'ε'}	Follow(C) = {')', 'ε'}
//
// Table:
//		| 'p' | '(' | ')' | '$' |
//	     ----------------------------
//	      S | err | r0  | err | Ok! |
//	     ----------------------------
//	      A | r1  | r2  | err | err |
//	     ----------------------------
//	      B | r3  | r4  | err | err |
//	     ----------------------------
//	      C | r5  | r6  | r7  | err |
//	     ----------------------------
//
// Aho, Lam, Sethi, Ullman: Compilers. Principles, Techniques, & Tools (Second Edition).
// 				Section 4.4.4 Nonrecursive Predictive Parsing. Algorithm 4.34.
// 				Page: 227/8. Fig: 4.20 Predictive parsing algorithm.
// input buffer: w$
// set ip to point to the first symbol of w;
// set X to the top stack symbol;
// while ( X != $ ) { // stack is not empty
// 	if ( X is a ) pop the stack and advance ip;
// 	else if ( X is a terminal ) error();
// 	else if ( M [ X , a] is an error entry ) error();
// 	else if ( M[X,a] = X → Y1Y2...Yk ) {
// 		output the production X → Y1Y2...Yk;
// 		pop the stack;
// 		push Yk,Yk-1,...,Y1 onto the stack, with Y1 on top;
// 	}
//	set X to the top stack symbol;
// }
//

//func (balance *Balance) Load(ps []datastore.Property) error {
//	return nil
//}

//func (balance *Balance) Save() ([]datastore.Property, error) {
//	return nil, nil
//}
