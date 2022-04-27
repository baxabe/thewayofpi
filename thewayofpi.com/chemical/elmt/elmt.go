package elmt

import (
	"fmt"
)

import (
	"thewayofpi.com/buffer/logs"
	"thewayofpi.com/buffer/random"
	"sort"
	"errors"
)

const (
	MaxElements = len(eTable)
)

type (
	Symbol      byte
	SymbolList  []Symbol
	PolarOrder  []Symbol
	EnergyOrder []Symbol
	LinksOrder  []Symbol
)

type (
	Atom struct {
		symbol   Symbol
		name     tName
		short    tShort
		energy   tEnergy
		links    tLinks
		polarity tPolarity
		decay    Symbol
	}
)

type (
	tName      string
	tShort     string
	tLinks     byte
	tLinkList  []tLinks
	tEnergy    int
	tPolarity  int
)

func Table(idx Symbol) *Atom {
	if int(idx) >= MaxElements {
		logs.Error(fmt.Errorf("elmt:Table - Index out of range: %v", idx))
		return &eTable[ZΩ]
	}
	return &eTable[idx]
}

func (e *Atom) Symbol() Symbol {
	return e.symbol
}

func (e *Atom) Name() string {
	return string(e.name)
}

func (e *Atom) Short() string {
	return string(e.short)
}

func (e *Atom) Energy() int {
	return int(e.energy)
}

func (e *Atom) Links() int {
	return int(e.links)
}

func (e *Atom) Polarity() int {
	return int(e.polarity)
}

func (e *Atom) Decay() Symbol {
	return e.decay
}

func (e *Atom) IsNull() bool {
	return e.symbol == 0
}

func (e *Atom) IsValid() bool {
	return e.symbol > ZΩ && e.symbol < Symbol(MaxElements)
}

func (e *Atom) String() string {
	return fmt.Sprintf("Symbol  : %v\nName    : \"%v\"\nShort   : \"%v\"\nEnergy  : %v\nLinks   : %v\nPolarity: %v\nDecay: %v\n", e.Short(), e.Name(), e.Short(), e.Energy(), e.Links(), e.Polarity(), e.Short())
}

// SymbolList

func (l SymbolList) Len() int {
	return len(l)
}

func (l SymbolList) Less(i, j int) bool {
	return l[i] < l[j]
}

func (l SymbolList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l SymbolList) Count() []int {
	return count(l)
}

func (l SymbolList) Parity() uint {
	return calcParity(l)
}

func (l SymbolList) Energy() int {
	return int(calcEnergy(l))
}

func (l SymbolList) Polarity() int {
	return int(calcPolarity(l))
}

func (l SymbolList) String() string {
	var result string
	for _, id := range l {
		result += fmt.Sprintf("%v\n", Table(id))
	}
	return result
}

func (l SymbolList) links() tLinkList {
	return getLinks(l)
}

// tLinkList

func (l tLinkList) Len() int {
	return len(l)
}

func (l tLinkList) Less(i, j int) bool {
	return l[i] < l[j]
}

func (l tLinkList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// PolarOrder

func (p PolarOrder) Len() int {
	return len(p)
}

func (p PolarOrder) Less(i, j int) bool {
	return Table(p[i]).Polarity() < Table(p[j]).Polarity() ||
		(Table(p[i]).Polarity() == Table(p[j]).Polarity() &&
			((Table(p[i]).Polarity() < 0 && Table(p[i]).Energy() > Table(p[j]).Energy()) ||
				Table(p[i]).Polarity() >= 0 && Table(p[i]).Energy() < Table(p[j]).Energy()))
}

func (p PolarOrder) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// EnergyOrder

func (p EnergyOrder) Len() int {
	return len(p)
}

func (p EnergyOrder) Less(i, j int) bool {
	return Table(p[i]).Energy() < Table(p[j]).Energy()
}

func (p EnergyOrder) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// LinksOrder

func (p LinksOrder) Len() int {
	return len(p)
}

func (p LinksOrder) Less(i, j int) bool {
	return Table(p[i]).Links() < Table(p[j]).Links() ||
		(Table(p[i]).Links() == Table(p[j]).Links() && Table(p[i]).Energy() < Table(p[j]).Energy())
}

func (p LinksOrder) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Stuff

func getLinks(elements SymbolList) tLinkList {
	var links tLinkList
	for _, e := range elements {
		if int(e) < MaxElements {
			if Table(e).Links() > 0 {
				links = append(links, Table(e).links)
			}
		} else {
			logs.Warning(errors.New("elmt:getLinks - element index out of range"))
		}
	}
	sort.Sort(sort.Reverse(links))
	return links
}

func count(idList SymbolList) []int {
	result := make([]int, MaxElements)
	for _, e := range idList {
		if int(e) < MaxElements {
			result[e]++
		} else {
			logs.Warning(errors.New("elmt:count - element index out of range"))
		}
	}
	return result
}

func calcParity(idList SymbolList) uint {
	if len(idList) == 1 {
		return uint(Table(idList[0]).links)
	}
	links := idList.links()
	for len(links) > 1 {
		links[0]--
		links[len(links)-1]--
		links = reorder(links)
	}
	if len(links) != 0 {
		return uint(links[0])
	}
	return 0
}

func reorder(links tLinkList) tLinkList {
	if len(links) < 2 {
		return links
	}
	if links[len(links)-1] == 0 {
		links = links[:len(links)-1]
	}
	if links[0] == 0 {
		links = links[1:]
	}
	if len(links) < 2 {
		return links
	}
	for i := 0; i < len(links)-1; i++ {
		if links[i+1] <= links[0] {
			links.Swap(0, i)
			break
		}
	}
	return links
}

func calcEnergy(elements SymbolList) tEnergy {
	var enrgy tEnergy
	for _, v := range elements {
		enrgy += Table(v).energy
	}
	return enrgy
}

func calcPolarity(elements SymbolList) tPolarity {
	var polarity tPolarity
	for _, v := range elements {
		polarity += Table(v).polarity
	}
	return polarity
}

// Testing

func NewTestIds(n int) SymbolList {
	var result SymbolList
	rnd := random.New()
	for i := 0; i < n; i++ {
		result = append(result, Symbol(rnd.Intn(MaxElements-1)+1))
	}
	sort.Sort(sort.Reverse(result))
	return result

	//return []TId{1, 1, 1}
	//return []TId{6, 6, 2}
	//return []tId{2, 2, 1}
	//return []tId{6, 2}
	//return []tId{23, 15, 15, 1}
	//return []tId{5, 1, 1, 1}
	//return []tId{2, 2, 2}
	//return []tId{1, 1, 1, 1}
	//return []TId{1, 4, 7, 8, 11, 12, 17, 21, 22, 23, 26}
}
