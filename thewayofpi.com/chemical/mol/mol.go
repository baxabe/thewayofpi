package mol

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

import (
	"thewayofpi.com/buffer/logs"
	"thewayofpi.com/buffer/chemical"
	"thewayofpi.com/buffer/chemical/elmt"
)

const (
	MinMoleculeElements = 2
	ParityOne = elmt.Aα
)

type (
	component struct {
		count   uint
		element elmt.Symbol
	}
	components []component
	Molecule struct {
		len        int
		components components
		parity     uint
		energy     int
		polarity   int
	}
)

func NewMolecule(elements elmt.SymbolList) *Molecule {
	if len(elements) == 0 {
		return NewEmptyMolecule()
	}
	comp, num := group(elements)
	if len(comp) == 0 {
		return NewEmptyMolecule()
	}
	par := elements.Parity()
	enrgy := elements.Energy()
	polarity := elements.Polarity()
	return &Molecule{len: num, components: comp, parity: par, energy: enrgy, polarity: polarity}
}

func NewEmptyMolecule() *Molecule {
	return &Molecule{len: 0, components: []component{}, parity: 0, energy: 0, polarity: 0}
}

func (m *Molecule) Formula() string {
	var formula string
	for _, v := range m.components {
		short := elmt.Table(v.element).Short()
		if v.count > 1 {
			formula += fmt.Sprintf("%v(%v)", short, v.count)
		} else {
			formula += fmt.Sprintf("%v", short)
		}
	}
	return formula
}

func (m Molecule) Expand() elmt.SymbolList {
	var result elmt.SymbolList
	for _, v := range m.components {
		for i := 0; i < int(v.count); i++ {
			result = append(result, v.element)
		}
	}
	sort.Sort(result)
	return result
}

func (m Molecule) Len() int {
	return m.len
}

func (m Molecule) Parity() uint {
	return m.parity
}

func (m Molecule) Energy() int {
	return m.energy
}

func (m Molecule) Polarity() int {
	return m.polarity
}

func (m Molecule) String() string {
	return fmt.Sprintf("Formula : %v\nParity  : %v\nEnergy  : %v\nPolarity: %v\n", m.Formula(), m.Parity(), m.Energy(), m.Polarity())
}

func (m Molecule) IsNull() bool {
	return m.components == nil || len(m.components) == 0
}

func (m *Molecule) AddElement(e elmt.Symbol, count uint) {
	if count > 0 && elmt.Table(e).IsValid() {
		symbols := m.Expand()
		for i := 0; i < int(count); i++ {
			symbols = append(symbols, e)
		}
		*m = *NewMolecule(symbols)
	}
}

func (m *Molecule) SubElement(e elmt.Symbol, count uint) {
	if count > 0 && elmt.Table(e).IsValid() {
		for _, v := range m.components {
			if v.element == e {
				v.count -= uint(math.Min(float64(count), float64(v.count)))
				*m = *NewMolecule(m.Expand())
				return
			}
		}
	}
}

func (m Molecule) CountElement(e elmt.Symbol) uint {
	if elmt.Table(e).IsValid() {
		for _, v := range m.components {
			if v.element == e {
				return v.count
			}
		}
	}
	return 0
}

func (m *Molecule) Adjust() {
	if m.parity > 0 {
		count := m.CountElement(elmt.Symbol(ParityOne))
		if m.parity < count {
			m.SubElement(elmt.Symbol(ParityOne), m.parity)
		} else {
			m.AddElement(elmt.Symbol(ParityOne), m.parity)
		}
		*m = *NewMolecule(m.Expand())
		if m.parity != 0 {
			logs.Error(errors.New("mol:Adjust - new parity is != 0"))
		}
	}
}

func (m *Molecule) Decay(propId chem.PropId) (err error) {
	ids := saneInput(m)
	switch propId {
	case chem.PolarId:
		sort.Sort(elmt.PolarOrder(ids))
		break
	case chem.ForceId:
		sort.Sort(elmt.EnergyOrder(ids))
		break
	case chem.LinkId:
		sort.Sort(elmt.LinksOrder(ids))
		break
	default:
		return logs.Error(fmt.Errorf("mol:Decay - propId is default: %v", propId))
	}
	if len(ids) >= MinMoleculeElements {
		ids[0] = elmt.Table(ids[0]).Decay()
		ids[len(ids)-1] = elmt.Table(ids[len(ids)-1]).Decay()
	} else {
		return logs.Error(fmt.Errorf("mol:Decay - len(exp) < MinMoleculeElements: %v", ids))
	}
	*m = *NewMolecule(ids)
	return nil
}

func (m *Molecule) decay(ids elmt.SymbolList) {
	if len(ids) >= MinMoleculeElements {
		ids[0] = elmt.Table(ids[0]).Decay()
		ids[len(ids)-1] = elmt.Table(ids[len(ids)-1]).Decay()
	}
	*m = *NewMolecule(ids)
}

func (m *Molecule) PolarDecay() {
	exp := saneInput(m)
	sort.Sort(elmt.PolarOrder(exp))
	m.decay(exp)
}

func (m *Molecule) EnergyDecay() {
	exp := saneInput(m)
	sort.Sort(elmt.EnergyOrder(exp))
	m.decay(exp)
}

func (m *Molecule) LinkDecay() {
	exp := saneInput(m)
	sort.Sort(elmt.LinksOrder(exp))
	m.decay(exp)
}

// Stuff

func group(elements elmt.SymbolList) ([]component, int) {
	table := elements.Count()
	return build(table)
}

func build(table []int) ([]component, int) {
	var result []component
	var num int
	for i := len(table) - 1; i >= 0; i-- {
		if table[i] > 0 {
			result = append(result, component{count: uint(table[i]), element: elmt.Table(elmt.Symbol(i)).Symbol()})
			num += table[i]
		}
	}
	return result, num
}

func saneInput(m *Molecule) elmt.SymbolList {
	if m.IsNull() {
		logs.Error(errors.New("mol:saneInput - null input"))
		return elmt.SymbolList{}
	}
	return m.Expand()
}

func (c components) Expand() elmt.SymbolList {
	var result elmt.SymbolList
	for _, v := range c {
		for i := 0; i < int(v.count); i++ {
			result = append(result, v.element)
		}
	}
	sort.Sort(result)
	return result
}

