package reac

import (
	"errors"
	"fmt"
)

import (
	"thewayofpi.com/buffer/chemical"
	"thewayofpi.com/buffer/chemical/elmt"
	"thewayofpi.com/buffer/chemical/mol"
	"thewayofpi.com/buffer/logs"
	"thewayofpi.com/buffer/math"
)

const (
	input = iota
	output
	void
)

type (
	opId int
)

const (
	PipeId opId = iota
	SplitId
	JoinId
	OpLen
)

type (
	stageSeed   map[opId]uint
	stageSchema []opId
)

type (
	Reaction struct {
		meta
		matrix
	}
	matrix struct {
		stages []stage
		phases []phase
	}
	meta struct {
		inputs  [][2]int // [][Column, Row]
		outputs [][2]int // [][Column, Row]
		cursor  int
	}
	stage struct {
		adapters
		eDelta int
	}
	adapters []adapter
	phase    []link
	adapter  struct {
		opId   opId
		propId chem.PropId
		op     operation
		inp    []int // inp[i] -> phases[prev][i]
		out    []int // out[i] -> phases[next][i]
		eDelta int
	}
	operation func(m []*mol.Molecule) ([]*mol.Molecule, int, error)
	link      struct {
		left   plug
		right  plug
		minLen int
		mol    mol.Molecule
	}
	plug struct {
		idx int
		pin int
	}
	io struct {
		inp int
		out int
	}
	opMap map[opId]io
)

var opSpec = opMap{
	PipeId:  io{inp: 1, out: 1},
	SplitId: io{inp: 1, out: 2},
	JoinId:  io{inp: 2, out: 1},
}

func New(start, end, steps int) (*Reaction, error) {
	r := &Reaction{nil, nil}
	matrix, err := buildMatrix(start, end, steps)
	if err != nil {
		return nil, logs.Error(fmt.Errorf("reac:New - error building matrix: %v", err))
	}
	r.phases = matrix.phases
	r.stages = matrix.stages
	r.inputs = make([][2]int, start)
	for i := 0; i < start; i++ {
		r.inputs[i][0] = 0
		r.inputs[i][1] = i
	}
	r.outputs = make([][2]int, end)
	for i := 0; i < end; i++ {
		r.outputs[i][0] = len(r.phases) - 1
		r.outputs[i][1] = i
	}
	return r, nil
}

func (r Reaction) Meta() *meta {
	return &r.meta
}

func (r Reaction) Stages() []stage {
	return r.stages
}

func (r Reaction) Phases() []phase {
	return r.phases
}

func (r Reaction) Init(params []mol.Molecule) error {
	if len(r.phases) == 0 {
		return logs.Error(fmt.Errorf("reac:Init - Reaction without phases"))
	}
	for _, p := range r.phases {
		for _, l := range p {
			l.mol = *mol.NewEmptyMolecule()
		}
	}
	if err := fillInputPhase(params, r.phases[0]); err != nil {
		return logs.Error(fmt.Errorf("reac:Init - error filling input phase: %v", err))
	}
	r.cursor = 0
	return nil
}

func (r Reaction) Results() []mol.Molecule {
	var result []mol.Molecule
	for _, o := range r.outputs {
		result = append(result, r.phases[o[0]][o[1]].mol)
	}
	return result
}

// TODO: func (r Reaction) State() ???? {}

func (r Reaction) Do() error {
	for i := r.cursor; i < len(r.stages); i++ {
		if err := r.doStep(); err != nil {
			return logs.Error(fmt.Errorf("reac:Do - error in stage %v: %v", i, err))
		}
	}
	return nil
}

func (r Reaction) DoStep() error {
	return r.doStep()
}

func (r Reaction) doStep() error {
	var energy int
	var err error
	if r.cursor >= len(r.stages) {
		return logs.Error(fmt.Errorf("reac:doStep - cursor out of range: %v - max: %v", r.cursor, len(r.stages)-1))
	}
	for i, a := range r.stages[r.cursor].adapters {
		if energy, err = a.run(r.phases[r.cursor], r.phases[r.cursor+1]); err != nil {
			return logs.Error(fmt.Errorf("reac:doStep - stage: %v; adapter: %v - error: %v", r.cursor, i, err))
		}
		r.stages[r.cursor].eDelta += energy
	}
	r.cursor++
	return nil
}

func (a *adapter) run(prev, next phase) (int, error) {
	var inp, out []*mol.Molecule
	var energy int
	var err error
	for _, i := range a.inp {
		inp = append(inp, &prev[i].mol)
	}
	if out, energy, err = a.op(inp); err != nil {
		return 0, logs.Error(fmt.Errorf("reac:run - Error reported by operation: %v", err))
	}
	if len(a.out) != len(out) {
		return 0, logs.Error(fmt.Errorf("reac:run - Output sizes don't match: [expected: %v | got: %v]", len(a.out), len(out)))
	}
	for _, i := range a.out {
		next[i].mol = *out[i]
	}
	return energy, nil
}

func selectOp(op opId, prop chem.PropId) (operation, error) {
	var selected operation
	var err error
	switch op {
	case PipeId:
		selected = opPipe()
		break
	case SplitId:
		selected = opSplit(prop)
		break
	case JoinId:
		selected = opJoin(prop)
		break
	default:
		err = logs.Error(fmt.Errorf("reac:selectOp - Unknown operation selected: %v", op))
	}
	return selected, err
}

func opPipe() operation {
	return func(m []*mol.Molecule) ([]*mol.Molecule, int, error) {
		var sl elmt.SymbolList
		var result []*mol.Molecule
		var err error
		if l := len(m); l == 0 {
			err = logs.Error(errors.New("reac:opPipe - Empty input"))
		} else {
			for _, val := range m {
				sl = val.Expand()
				newMol := mol.NewMolecule(sl)
				result = append(result, newMol)
			}
		}
		return result, 0, err
	}
}

func opJoin(id chem.PropId) operation {
	var prop = id
	return func(m []*mol.Molecule) ([]*mol.Molecule, int, error) {
		const MinInputs, Outputs = 2, 1
		var sl elmt.SymbolList
		var eInp, eOut int
		var err error
		var result []*mol.Molecule
		if l := len(m); l < MinInputs {
			err = logs.Error(fmt.Errorf("reac:opJoin - Too few inputs: %v", l))
		} else {
			for _, im := range m {
				sl = append(sl, im.Expand()...)
				eInp += im.Energy()
			}
			newMol := mol.NewMolecule(sl)
			if err = newMol.Decay(prop); err != nil {
				err = logs.Error(fmt.Errorf("reac:opJoin - Decay returns error: %v", prop))
			} else {
				result = make([]*mol.Molecule, Outputs)
				result[0] = newMol
				eOut = newMol.Energy() - eInp
			}
		}
		return result, eOut, err
	}
}

func opSplit(id chem.PropId) operation {
	var prop = id
	return func(m []*mol.Molecule) ([]*mol.Molecule, int, error) {
		const Inputs, Outputs = 1, 2
		var sl elmt.SymbolList
		var eInp, eOut, idx int
		var err error
		var result []*mol.Molecule
		if l := len(m); l != Inputs {
			err = logs.Error(fmt.Errorf("reac:opSplit - Incorrect number of inputs: %v", l))
		} else {
			sl = m[0].Expand()
			if l = len(sl); l < mol.MinMoleculeElements*Outputs {
				err = logs.Error(fmt.Errorf("reac:opSplit - molecule too short: len = %v", l))
			} else {
				eInp = m[0].Energy()
				result = make([]*mol.Molecule, Outputs)
				proto := make([]elmt.SymbolList, Outputs)
				balance := balancer(Outputs, prop)
				for _, el := range sl {
					if idx, err = balance(el); err != nil {
						break
					}
					proto[idx] = append(proto[idx], el)
				}
				if err == nil {
					for i, sl := range proto {
						result[i] = mol.NewMolecule(sl)
						eOut += result[i].Energy()
					}
				}
			}
		}
		return result, eOut - eInp, err
	}
}

func balancer(n int, propId chem.PropId) func(elmt.Symbol) (int, error) {
	count := make([]int, n)
	prop := propId
	return func(symbol elmt.Symbol) (int, error) {
		var val, idx int
		var err error
		switch prop {
		case chem.PolarId:
			val = int(elmt.Table(symbol).Polarity())
			break
		case chem.ForceId:
			val = int(elmt.Table(symbol).Energy())
			break
		case chem.LinkId:
			val = int(elmt.Table(symbol).Links())
			break
		default:
			err = logs.Error(fmt.Errorf("reac:deliver - Unknown property selected: %v", prop))
		}
		if err == nil {
			_, idx, err = math.FindMin(count)
			if err == nil {
				count[idx] += val
			}
		}
		return idx, err
	}
}

func buildMatrix(start, end, steps int) (*matrix, error) {
	if start == 0 {
		return nil, errors.New("reac:buildMatrix - zero inputs")
	}
	chain, err := buildStepsChain(uint(start), uint(end), uint(steps))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("reac:buildMatrix - %v", err))
	}
	result := &matrix{stages: nil, phases: nil}
	result.phases = append(result.phases, buildInputPhase(start))
	seed := make(stageSeed)
	var st stage
	var ph phase
	for i := 1; i < len(chain); i++ {
		if seed, err = buildStageSeed(chain[i-1], chain[i]); err != nil {
			return nil, fmt.Errorf("reac:buildMatrix - %v", err)
		}
		if &st, ph, err = buildStage(result.phases[i-1], seed); err != nil {
			return nil, fmt.Errorf("reac:buildMatrix - %v", err)
		}
		result.phases = append(result.phases, ph)
		result.stages = append(result.stages, st)
	}
	completeOutputPhase(result.phases[len(result.phases)-1])
	if err = linkMatrix(result); err != nil {
		return nil, fmt.Errorf("reac:buildMatrix - %v", err)
	}
	if err = fillMinLen(result); err != nil {
		return nil, fmt.Errorf("reac:buildMatrix - %v", err)
	}
	return result, nil
}

func linkMatrix(m *matrix) error {
	for i, s := range m.stages {
		lIdx, rIdx := 0, 0
		for j, a := range s.adapters {
			for l := range a.inp {
				a.inp[l] = lIdx
				m.phases[i][lIdx].right.idx = j
				m.phases[i][lIdx].right.pin = l
				lIdx++
			}
			for r := range a.out {
				a.inp[r] = rIdx
				rIdx++
			}
		}
		if len(m.phases[i]) != lIdx || len(m.phases[i+1]) != rIdx {
			return fmt.Errorf("reac:linkMatrix - [phases[%v]: %v, lIdx = %v, phases[%v]: %v, rIdx = %v]",
				i, len(m.phases[i]), lIdx, i+1, len(m.phases[i+1]), rIdx)
		}
	}
	return nil
}

func fillMinLen(m *matrix) error {
	if err := fillMinLenToRight(m); err != nil {
		return err
	}
	if err := fillMinLenToLeft(m); err != nil {
		return err
	}
	return nil
}

func fillMinLenToRight(m *matrix) error {
	for i, s := range m.stages {
		for _, a := range s.adapters {
			switch a.opId {
			case PipeId:
				break
			case SplitId:
				break
			case JoinId:
				m.phases[i+1][a.out[0]].minLen = math.Max(m.phases[i+1][a.out[0]].minLen, len(a.inp)*mol.MinMoleculeElements)
				break
			default:
				return logs.Error(fmt.Errorf("reac:fillMinLenToRight - Unknown opId: %v", a.opId))
			}
		}
	}
	return nil
}

func fillMinLenToLeft(m *matrix) error {
	for i := len(m.stages) - 1; i >= 0; i-- {
		for _, a := range m.stages[i].adapters {
			switch a.opId {
			case PipeId:
				break
			case SplitId:
				m.phases[i][a.inp[0]].minLen = math.Max(m.phases[i][a.inp[0]].minLen, len(a.out)*mol.MinMoleculeElements)
				break
			case JoinId:
				break
			default:
				return logs.Error(fmt.Errorf("reac:fillMinLenToLeft - Unknown opId: %v", a.opId))
			}
		}
	}
	return nil
}

//func TestBuildStepsChain(in, st, out uint) ([]uint, error) {
//	return buildStepsChain(in, st, out)
//}
