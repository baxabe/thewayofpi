package balances

import (
	"fmt"
)

type (
	lexem rune
)

const (
	OPEN    lexem = '\u0028' // Open tree '('     dec: 40
	CLOSE   lexem = '\u0029' // Close tree ')'    dec: 41
	POINT   lexem = '\u002E' // Weight '.'        dec: 46
	NOTERMS lexem = '\u0053' // No terminal S 'S' dec: 83
	NOTERMA lexem = '\u0041' // No terminal A 'A' dec: 65
	NOTERMB lexem = '\u0042' // No terminal B 'B' dec: 66
	NOTERMC lexem = '\u0043' // No terminal C 'C' dec: 67
	EPSILON lexem = '\u03F5' // Epsilon ε 'ε'	dec: 1013
	EOF     lexem = '\u0024' // End of input '$'	dec: 36
)

type (
	token struct {
		id    lexem
		value point
	}
)

func newToken(id lexem) token {
	return token{id, emptyPoint()}
}

func (t token) idx() int {
	switch t.id {
	case POINT:
		return 0
	case OPEN:
		return 1
	case CLOSE:
		return 2
	case EOF:
		return 3
	}
	return -1
}

func (t token) idy() int {
	switch t.id {
	case NOTERMS:
		return 0
	case NOTERMA:
		return 1
	case NOTERMB:
		return 2
	case NOTERMC:
		return 3
	}
	return -1
}

func (t token) isTerminal() bool {
	switch t.id {
	case POINT, OPEN, CLOSE, EPSILON:
		return true
	}
	return false
}

func (t token) isEOF() bool {
	return t.id == EOF
}

func (t token) isEpsilon() bool {
	return t.id == EPSILON
}

func (t token) equals(tt token) bool {
	return t.id == tt.id
}

func (t token) String() string {
	var result string
	switch t.id {
	case POINT:
		result += fmt.Sprintf("%v", t.value)
	case OPEN:
		result += fmt.Sprintf("%v", t.value)
		fallthrough
	case CLOSE:
		result += fmt.Sprintf("%c", t.id)
	}
	return result
}
