package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
)

const (
	strJug =	"Players"
	strDims =	"Dim"
	strTodo =	"Todo"
	strDone =	"Done"
	strRoots =	"Roots"
	strCount =	"nPulses"
	strRate =	"Rate"
	strSleep =	"Sleep"
	strSteep =	"Steep"
)

type (
	point struct {
		x	int
		y	int
	}
	attrib struct {
		fg	tb.Attribute
		bg	tb.Attribute
	}
)

type (
	screen struct {
		wx			int
		wy			int
		posJug		point
		posDims		point
		posTodo		point
		posDone		point
		posRoots	point
		posCount	point
		posConf		point
	}
)

var (
	scr	screen
	rate point
)

func showConf() {
	atr := attrib{tb.ColorDefault, tb.ColorDefault}
	place(point{scr.posConf.x + len(strRate) + 3, scr.posConf.y}, atr, fmt.Sprintf("%d/%d", rate.x, rate.y))
	place(point{scr.posConf.x + len(strSteep) + 2, scr.posConf.y + 1}, atr, fmt.Sprintf("%8d", *steep))
	place(point{scr.posConf.x + len(strSleep) + 2, scr.posConf.y + 2}, atr, fmt.Sprintf("%8d", *sleep))
	flush()
}

func showPulses(pulse int) {
	if show {
		atr := attrib{tb.ColorDefault, tb.ColorDefault}
		place(point{scr.posCount.x + len(strCount) + 2, scr.posCount.y}, atr, fmt.Sprintf("%15d", pulse))
		flush()
	}
}

func showPlayers() {
	if show {
		atr := attrib{tb.ColorDefault, tb.ColorDefault}
		place(point{scr.posJug.x + len(strJug) + 2, scr.posJug.y}, atr, fmt.Sprintf("%15d", plCount))
		flush()
	}
}

func showData() {
	if show {
		atr := attrib{tb.ColorDefault, tb.ColorDefault}
		for y, idx := scr.posTodo.y + 1, 0; y < scr.wy - 4; y, idx = y + 1, idx + 1 {
			place(point{scr.posTodo.x, y}, atr, fmt.Sprintf("%10d", treelist[idx].todo.Len()))
			place(point{scr.posDone.x, y}, atr, fmt.Sprintf("%10d", treelist[idx].done.Len()))
			place(point{scr.posRoots.x, y}, atr, fmt.Sprintf("%10d", treelist[idx].roots))
		}
		flush()
	}
}

func showLabels() {
	clear()
	atr := attrib{tb.ColorDefault, tb.ColorDefault}
	place(scr.posJug, atr, strJug+":")
	place(scr.posDims, atr, strDims+":")
	for y:= scr.posDims.y; y < scr.wy - 4; y++ {
		place(point{scr.posDims.x, y + 1}, atr, fmt.Sprintf("%3d", y - scr.posDims.y))
	}
	place(scr.posTodo, atr, strTodo+":")
	place(scr.posDone, atr, strDone+":")
	place(scr.posRoots, atr, strRoots+":")
	place(scr.posCount, atr, strCount+":")
	place(scr.posConf, atr, strRate+":")
	place(point{scr.posConf.x, scr.posConf.y + 1}, atr, fmt.Sprintf(strSteep+":"))
	place(point{scr.posConf.x, scr.posConf.y + 2}, atr, fmt.Sprintf(strSleep+":"))
	flush()
}

func clear() {
	if err := tb.Clear(tb.ColorDefault, tb.ColorDefault); err != nil {
		fmt.Println("Clear error")
	}
}

func flush() {
	if err := tb.Flush(); err != nil {
		fmt.Println("Flush error")
	}
}

func place(p point, a attrib, msg string) {
	x := p.x
	for _, car := range msg {
		tb.SetCell(x, p.y, car, a.fg, a.bg)
		x++
	}
}

func logs(s string) {
	place(point{scr.wx/2 - len(s)/2, scr.wy - 1}, attrib{tb.ColorDefault, tb.ColorDefault}, s)
	flush()
}

func logn(n int) {
	s := fmt.Sprintf("%010d", n)
	place(point{scr.wx/2 - len(s)/2, scr.wy - 1}, attrib{tb.ColorDefault, tb.ColorDefault}, s)
	flush()
}
