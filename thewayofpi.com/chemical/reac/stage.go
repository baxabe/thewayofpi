package reac

import (
	"fmt"
	"math"
)

import (
	"thewayofpi.com/buffer/chemical"
	"thewayofpi.com/buffer/random"
)

func buildStageSeed(inputs, outputs uint) (stageSeed, error) {
	if inputs < 1 || outputs < 1 {
		return nil, fmt.Errorf("reac:buildStageSpec - invalid constraints: inputs = %v; outputs = %v", inputs, outputs)
	}
	result := make(stageSeed)
	diff := uint(math.Abs(float64(inputs - outputs)))
	if inputs > outputs {
		result[JoinId] = diff
		inputs -= 2 * diff
		outputs -= diff
	} else if inputs < outputs {
		result[SplitId] = diff
		inputs -= diff
		outputs -= 2 * diff
	}
	// inputs == outputs
	rnd := random.New()
	remain := int(inputs)
	for i := remain; i > 0; i /= 2 {
		if int(rnd.Intn(10*remain)) < i {
			result[PipeId]++
			inputs--
		}
	}
	result[PipeId] += inputs % 3
	result[SplitId] += inputs / 3
	result[JoinId] += inputs / 3
	return result, nil
}

func buildStageSchema(seed stageSeed) (stageSchema, error) {
	if int(OpLen) < len(seed) {
		return nil, fmt.Errorf("reac:buildStageSchema - error in seed: OpLen = %v; len(seed) = %v", OpLen, len(seed))
	}
	var schema stageSchema
	for op, val := range seed {
		for i := 0; i < int(val); i++ {
			schema = append(schema, op)
		}
	}
	rnd := random.New()
	perm := rnd.Perm(len(schema))
	var result [len(schema)]opId
	for i := 0; i < len(schema); i++ {
		result[i] = schema[perm[i]]
	}
	return result[:], nil
}

func buildStage(inp phase, seed stageSeed) (*stage, phase, error) {
	var sch stageSchema
	var err error
	var back, fwd int
	if sch, err = buildStageSchema(seed); err != nil {
		return nil, nil, fmt.Errorf("reac:buildStage - error in buildStageSchema: %v", err)
	}
	rnd := random.New()
	adap := make(adapters, len(sch))
	for i, opId := range sch {
		if opId < 0 || opId >= OpLen {
			return nil, nil, fmt.Errorf("reac:buildStage - unknown opId: %v", opId)
		}
		prop := chem.PropId(rnd.Intn(int(chem.PropLen)))
		var op operation
		if op, err = selectOp(opId, prop); err != nil {
			return nil, nil, fmt.Errorf("reac:buildStage - error in selectOp: %v", err)
		}
		adap[i].opId = opId
		adap[i].propId = prop
		adap[i].op = op
		adap[i].inp = make([]int, opSpec[opId].inp)
		adap[i].out = make([]int, opSpec[opId].out)
		back += opSpec[opId].inp
		fwd += opSpec[opId].out
	}
	if back != len(inp) {
		return nil, nil, fmt.Errorf("reac:buildStage - back{%v} != len(inp){%v}", back, len(inp))
	}
	stage := &stage{adap, 0}
	next := make(phase, fwd)
	return stage, next, nil
}
