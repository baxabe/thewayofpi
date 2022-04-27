package reac

import (
	"errors"
)

import (
	"thewayofpi.com/buffer/chemical/mol"
)

func buildInputPhase(inputs int) phase {
	result := make(phase, inputs)
	for i, v := range result {
		v.mol = *mol.NewEmptyMolecule()
		v.minLen = mol.MinMoleculeElements
		v.left = plug{idx: input, pin: i}
		v.right = plug{idx: void, pin: void}
	}
	return result
}

func fillInputPhase(params []mol.Molecule, inPhase phase) error {
	if len(params) != len(inPhase) {
		return errors.New("reac:fillInputPhase - inputs length don't match")
	}
	for i, m := range params {
		inPhase[i].mol = m
	}
	return nil
}

func completeOutputPhase(ph phase) {
	for _, v := range ph {
		v.right = plug{idx: output, pin: 0}
	}
}
