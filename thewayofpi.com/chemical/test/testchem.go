package main

import (
	"fmt"
	"sort"
)

import (
	"thewayofpi.com/buffer/chemical/elmt"
	"thewayofpi.com/buffer/chemical/mol"
	//"thewayofpi.com/buffer/chemical/reac"
	"thewayofpi.com/buffer/random"
	"errors"
	"math"
	// "thewayofpi.com/buffer/chemical/subs"
)

func main() {
	//calcParity(4, false)
	//cmpParity(4)
	//findMolecules(5)
	//testPolarSort()
	//testJoinReaction(4, 10)
	//testSplitPolarReaction(4, 10)
	//testSplitForceReaction(4, 10)
	//testSplitLinksReaction(4, 10)
	//calcReactionStatistics(6, 16, 1000)
	//testLog2Function(200)
	//voidTest()
	//testNormal(100)
	testRandPerm(5)
}

func calcParity(n int, detail bool) {
	test := elmt.NewTestIds(n)
	if detail {
		fmt.Println(test)
	}
	fmt.Println(test.Parity())
	molecule := mol.NewMolecule(test)
	fmt.Println(molecule)
}

func cmpParity(n int) {
	test := elmt.NewTestIds(n)
	molecule := mol.NewMolecule(test)
	fmt.Println(test.Parity(), molecule.Parity())
}

func findMolecule(minElem, maxElem int) (*mol.Molecule, error) {
	tries := 10000
	rnd := random.New()
	found := false
	for !found && tries > 0 {
		num := rnd.Intn(maxElem - minElem + 1) + minElem
		ids := elmt.NewTestIds(num)
		molec := mol.NewMolecule(ids)
		if molec.Parity() == 0 {
			return molec, nil
		}
		tries--
	}
	return nil, errors.New("Could not found molecule")
}

func testPolarSort(){
	list := make([]elmt.Symbol, elmt.MaxElements)
	for i, _ := range(list) {
		list[i] = elmt.Symbol(i)
	}
	sort.Sort(elmt.PolarOrder(list))
	for _, id := range(list) {
		fmt.Printf("Symbol: %v\tPolarity: %v\tForce: %v\n", id, elmt.Table(id).Polarity(), elmt.Table(id).Energy())
	}

}

func findMolecules(n int) {
	const MaxElem = 6
	const MinElem = 2
	rnd := random.New()
	found, tries := 0, 0
	for found < n {
		num := rnd.Intn(MaxElem - MinElem + 1) + MinElem
		ids := elmt.NewTestIds(num)
		molec := mol.NewMolecule(ids)
		if molec.Parity() == 0 {
			fmt.Println(molec.Formula())
			found++
		}
		tries++
	}
	fmt.Printf("%v molecules found in %v tries\n", found, tries)
}

//func testJoinReaction(minElem, maxElem int) {
//	fmt.Println("Testing join:")
//	m0, err := findMolecule(minElem, maxElem)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	} else {
//		fmt.Println("Molecule A:", m0.Formula())
//	}
//	m1, err := findMolecule(minElem, maxElem)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	} else {
//		fmt.Println("Molecule B:", m1.Formula())
//	}
//	if enrgy, err := m0.PolarJoin(m1); err != nil {
//		fmt.Println(err.Error())
//	} else {
//		fmt.Printf("Formula: %v - Energy: %v\n", m0.Formula(), enrgy)
//	}
//}
//
//func testSplitPolarReaction(minElem, maxElem int) {
//	fmt.Println("Testing split polar:")
//	m, err := findMolecule(minElem, maxElem)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	} else {
//		fmt.Println("Molecule:", m.Formula())
//	}
//	if j, k, enrgy, err := m.PolarSplit(); err != nil {
//		fmt.Println(err.Error())
//	} else {
//		fmt.Println("Chunk A:", j.Formula())
//		fmt.Println("Chunk B:", k.Formula())
//		fmt.Println("Energy :", enrgy)
//	}
//}
//
//func testSplitForceReaction(minElem, maxElem int) {
//	fmt.Println("Testing split force:")
//	m, err := findMolecule(minElem, maxElem)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	} else {
//		fmt.Println("Molecule:", m.Formula())
//	}
//	if j, k, enrgy, err := m.ForceSplit(); err != nil {
//		fmt.Println(err.Error())
//	} else {
//		fmt.Println("Chunk A:", j.Formula())
//		fmt.Println("Chunk B:", k.Formula())
//		fmt.Println("Energy :", enrgy)
//	}
//}
//
//func testSplitLinksReaction(minElem, maxElem int) {
//	fmt.Println("Testing split links:")
//	m, err := findMolecule(minElem, maxElem)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	} else {
//		fmt.Println("Molecule:", m.Formula())
//	}
//	if j, k, enrgy, err := m.LinkSplit(); err != nil {
//		fmt.Println(err.Error())
//	} else {
//		fmt.Println("Chunk A:", j.Formula())
//		fmt.Println("Chunk B:", k.Formula())
//		fmt.Println("Energy :", enrgy)
//	}
//}
//
//func calcReactionStatistics(minElem, maxElem, count int) error {
//	fmt.Println("Calculating reactions stats:")
//	var p, f, l, pf, pl, fl, pfl int
//	okp, okf := false, false
//	for i := 0; i < count; i++ {
//		m, err := findMolecule(minElem, maxElem)
//		if err != nil {
//			return err
//		}
//		if _, _, _, err := m.PolarSplit(); err == nil {
//			okp = true
//			p++
//		}
//		if _, _, _, err := m.ForceSplit(); err == nil {
//			okf = true
//			f++
//			if okp {
//				pf++
//			}
//		}
//		if _, _, _, err := m.LinkSplit(); err == nil {
//			l++
//			if okp {
//				pl++
//			}
//			if okf {
//				fl++
//			}
//			if okp && okf {
//				pfl++
//			}
//		}
//		okp, okf = false, false
//	}
//	fmt.Printf("Stats results (%v cases)\n", count)
//	fmt.Println("\tpolar: ", p)
//	fmt.Println("\tforce: ", f)
//	fmt.Println("\tlinks: ", l)
//	fmt.Println("\tp + l: ", pl)
//	fmt.Println("\tp + f: ", pf)
//	fmt.Println("\tf + l: ", fl)
//	fmt.Println("\tp+f+l: ", pfl)
//	return nil
//}

func testLog2Function(n int) {
	fmt.Println("Testing Log2 values for: ", n)
	for i := n; i >= 1; i = int(math.Log2(float64(i))) {
		fmt.Println(i)
	}
}

func voidTest() {
	fmt.Println(math.Ceil(float64(6)/float64(4)))
}

func testNormal(n int) {
	//var res [100]int
	rnd := random.New()
	for i := 0; i < n; i++ {
		//x := int(math.Floor(math.Abs(rnd.NormFloat64() * 10)))
		//for j := 0; j < 100; j++ {
		//	if x >= j && x < j + 1 {
		//		res[j]++
		//		break
		//	}
		//}
		fmt.Println(rnd.NormFloat64())
	}
	//fmt.Println(res)
}

func testRandPerm(n int) {
	rnd := random.New()
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := rnd.Intn(i + 1)
		fmt.Printf("i: %v\nj: %v\nm[i]: %v\nm[j]: %v\n", i, j, m[i], m[j])
		m[i] = m[j]
		m[j] = i
		fmt.Printf("m: %v\n\n", m)
	}
}