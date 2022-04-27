package main

import (
	"fmt"
	"os"
	"strings"
)

import (
	piMath "thewayofpi.com/buffer/math"
)

const (
	MaxElements = 33
)

var (
	LatinChars = [...]string{"Z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}
	GreekChars = [...]string{"Omega", "Alpha", "Beta", "Gamma", "Delta", "Epsilon", "Zeta", "Eta", "Theta", "Iota",
		"Kappa", "Lamda", "Mu", "Nu", "Xi", "Omicron", "Pi"}
)

func main() {
	buildElementTable()
}

func buildSymbols() ([]string, []string) {
	const first rune = 'A'
	const last rune = 'P'
	const alpha rune = '\u03B1' // small alpha
	const pi rune = '\u03C0'    // small pi
	const omega rune = '\u03A9'
	var symbols [MaxElements]string
	var names [MaxElements]string
	symbols[0] = fmt.Sprintf("%v%c", LatinChars[0], omega)
	names[0] = fmt.Sprintf("%v.%v", LatinChars[0], GreekChars[0])
	var count = 1
	for j, g := alpha, 1; j <= pi && count < MaxElements; j, g = j+1, g+1 {
		for i, l := first, 1; i <= last && count < MaxElements; i, l, count = i+1, l+1, count+1 {
			symbols[count] = fmt.Sprintf("%c%c", i, j)
			names[count] = fmt.Sprintf("%v.%v", LatinChars[l], GreekChars[g])
		}
	}
	return symbols[:], names[:]
}

func buildConst(symbols []string) string {
	result := fmt.Sprintf("const (\n\t%v Symbol = iota\n", symbols[0])
	for i := 1; i < len(symbols); i++ {
		result += fmt.Sprintf("\t%v\n", symbols[i])
	}
	result += fmt.Sprintf(")\n")
	return result
}

func buildLinks() []int {
	var links [MaxElements]int
	piDigits := piMath.PiDigits(20)
	links[0] = 0
	var i, p = 0, 0
	for i < len(links) {
		links[i] = 0
		var j = 1
		for i+j < len(links) && j < int(piDigits[p]-'0') {
			links[i+j] = j
			j++
		}
		i += int(piDigits[p] - '0')
		p++
	}
	return links[:]
}

func buildPolarity(links []int) []string {
	var result = make([]string, len(links))
	//var sign string
	for i, n := range links {
		count := countBits(n)
		if (count+n)%2 == 1 {
			count = -1 * count
		}
		result[i] = fmt.Sprintf("%2d", count)
	}
	return result[:]
}

func countBits(n int) int {
	count := 0
	for n > 0 {
		count += n % 2
		n /= 2
	}
	return count
}

func buildDecay(links []int) []int {
	var decay = make([]int, len(links))
	for i := len(links) - 1; i >= 0; i-- {
		decay[i] = i
		for j := i - 1; j > 0; j-- {
			if links[i] == links[j] {
				decay[i] = j
				break
			}
		}
	}
	return decay[:]
}

func buildElementTable() {
	symbols, names := buildSymbols()
	maxnamelen := len(names[0])
	for _, name := range names {
		if l := len(name); l > maxnamelen {
			maxnamelen = l
		}
	}
	constants := buildConst(symbols)
	links := buildLinks()
	polar := buildPolarity(links)
	decay := buildDecay(links)
	table := fmt.Sprintf("var eTable = [...]Atom{\n")
	for i, n := range symbols {
		gap := strings.Repeat(" ", maxnamelen-len(names[i]))
		table += fmt.Sprintf("\t{symbol: %v,\tname: \"%v\",%v\tshort: \"%v\",\tenergy: %2d,\tlinks: %v,\tpolarity: %v,\tdecay: %v},\n", n, names[i], gap, n, i, links[i], polar[i], symbols[decay[i]])
	}
	table += fmt.Sprintf("}\n")
	content := fmt.Sprintf("package elmt\n// This file is generated by buildElementTable in builder.go\n\n%v\n%v", constants, table)
	writeFile("buffer/chemical/elmt/table.go", content)
}

func writeFile(path, content string) {
	if f, err := os.Create(path); err != nil {
		fmt.Println("ERROR creating file: ", err.Error())
	} else {
		defer f.Close()
		if n, err := f.WriteString(content); err != nil {
			fmt.Println("ERROR writing file: ", err.Error())
		} else {
			fmt.Printf("File saved: '%v'. Bytes written: %v\n", path, n)
		}
	}
}
