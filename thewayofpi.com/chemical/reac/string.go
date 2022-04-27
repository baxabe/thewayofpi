package reac

import (
	"thewayofpi.com/buffer/util/board"
)

const (
	S0 = '\u0020' //   SPACE
	H0 = '\u2500' // ─ BOX DRAWINGS LIGHT HORIZONTAL
	V0 = '\u2502' // │ BOX DRAWINGS LIGHT VERTICAL
	RC = '\u251C' // ├ BOX DRAWINGS LIGHT VERTICAL AND RIGHT
	LC = '\u2524' // ┤ BOX DRAWINGS LIGHT VERTICAL AND LEFT
	UL = '\u256D' // ╭ BOX DRAWINGS LIGHT ARC DOWN AND RIGHT
	UR = '\u256E' // ╮ BOX DRAWINGS LIGHT ARC DOWN AND LEFT
	BR = '\u256F' // ╯ BOX DRAWINGS LIGHT ARC UP AND LEFT
	BL = '\u2570' // ╰ BOX DRAWINGS LIGHT ARC UP AND RIGHT
)

type (
	box   [5][3]rune
)

var a21 = box{
	{UL, H0, UR},
	{LC, S0, V0},
	{V0, S0, RC},
	{LC, S0, V0},
	{BL, H0, BR}}

var a12 = box{
	{UL, H0, UR},
	{V0, S0, RC},
	{LC, S0, V0},
	{V0, S0, RC},
	{BL, H0, BR}}

var a11 = box{
	{S0, S0, S0},
	{S0, S0, S0},
	{H0, H0, H0},
	{S0, S0, S0},
	{S0, S0, S0}}

func (r Reaction) String() string {
	b := board.New(2, 3)
	return b.String()
}
