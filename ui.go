//-----------------------------------------------------------------------------
// Copyright (C) Microsoft. All rights reserved.
// Licensed under the MIT license.
// See LICENSE.txt file in the project root for full license information.
//-----------------------------------------------------------------------------
package main

import (
	"fmt"
	tm "github.com/nsf/termbox-go"
	"math"
	"time"
)

const (
	lefttop      = 0
	horizontal   = 1
	righttop     = 2
	vertical     = 3
	leftbottom   = 4
	rightbottom  = 5
	middlebottom = 6
	middletop    = 7
	middleleft   = 8
	middleright  = 9
	middlemiddle = 10
	space        = 11
	box1         = 12
	box2         = 13
	box3         = 14
	box4         = 15
	uparrow      = 16
	dnarrow      = 17
)

var symbols = []rune{'┌', '─', '┐', '│', '└', '┘', '┴', '┬', '├', '┤', '┼', ' ', '░', '▒', '▓', '█', '↑', '↓'}

const (
	justifyLeft = iota
	justifyRight
	justifyCenter
)

const (
	border = iota
	noBorder
)

type table struct {
	ccount  int
	cwidth  []int
	x       int
	y       int
	cr      int
	justify int
	border  int
}

func (t *table) drawTblRow(ledge, redge, middle, spr rune, fg, bg tm.Attribute) {
	twidth := t.ccount + 1
	for _, w := range t.cwidth {
		twidth += w
	}

	for i := 0; i < twidth; i++ {
		tm.SetCell(t.x+i, t.y+t.cr, middle, fg, bg)
	}

	if t.border == border {
		tm.SetCell(t.x, t.y+t.cr, ledge, fg, bg)
		tm.SetCell(t.x+twidth, t.y+t.cr, redge, fg, bg)
	}

	o := 0
	for c, w := range t.cwidth {
		o += w + 1
		if c < t.ccount-1 {
			tm.SetCell(t.x+o, t.y+t.cr, spr, fg, bg)
		}
	}
	t.cr++
}

func (t *table) addTblRow(row []string) {
	t.drawTblRow(symbols[vertical], symbols[vertical], symbols[space],
		symbols[vertical], tm.ColorDefault, tm.ColorDefault)
	t.cr--

	o := 1
	align_offset := 0
	for i := 0; i < t.ccount; i++ {
		w := t.cwidth[i]
		var s string
		if t.justify == justifyLeft {
			s = fmt.Sprintf("%-*s", w, row[i])
			if i == 0 && t.border == noBorder {
				align_offset = -1
			}
		} else {
			s = fmt.Sprintf("%*s", w, row[i])
		}
		printText(t.x+o+align_offset, t.y+t.cr, w, s, tm.ColorDefault, tm.ColorDefault)
		o += w + 1
		align_offset = 0
	}

	t.cr++
}

func (t *table) addTblSpr() {
	t.drawTblRow(symbols[middleleft], symbols[middleright], symbols[horizontal],
		symbols[middlemiddle], tm.ColorDefault, tm.ColorDefault)
}

func (t *table) addTblHdr() {
	t.drawTblRow(symbols[lefttop], symbols[righttop], symbols[horizontal],
		symbols[middletop], tm.ColorDefault, tm.ColorDefault)
}

func (t *table) addTblFtr() {
	t.drawTblRow(symbols[leftbottom], symbols[rightbottom], symbols[horizontal],
		symbols[middlebottom], tm.ColorDefault, tm.ColorDefault)
}

func printHLineText(x, y int, w int, text string) {
	for i := 0; i < w; i++ {
		tm.SetCell(x+i, y, symbols[horizontal], tm.ColorWhite, tm.ColorDefault)
	}
	offset := (w - len(text)) / 2
	textArr := []rune(text)
	for i := 0; i < len(text); i++ {
		tm.SetCell(x+offset+i, y, textArr[i], tm.ColorWhite, tm.ColorDefault)
	}
}

func printVLine(x, y int, h int) {
	tm.SetCell(x, y, symbols[middletop], tm.ColorWhite, tm.ColorDefault)
	for i := 1; i < h; i++ {
		tm.SetCell(x, y+i, symbols[vertical], tm.ColorWhite, tm.ColorDefault)
	}
}

func printText(x, y, w int, text string, fg, bg tm.Attribute) {
	textArr := []rune(text)
	for i := 0; i < w; i++ {
		tm.SetCell(x+i, y, ' ', fg, bg)
	}
	for i := 0; i < len(textArr); i++ {
		tm.SetCell(x+i, y, textArr[i], fg, bg)
	}
}

func printCenterText(x, y, w int, text string, fg, bg tm.Attribute) {
	offset := (w - len(text)) / 2
	textArr := []rune(text)
	for i := 0; i < w; i++ {
		tm.SetCell(x+i, y, ' ', fg, bg)
	}
	for i := 0; i < len(textArr); i++ {
		tm.SetCell(x+offset+i, y, textArr[i], fg, bg)
	}
}

func printHLine(x, y int, w int) {
	for i := 0; i < w; i++ {
		tm.SetCell(x+i, y, symbols[horizontal], tm.ColorWhite, tm.ColorDefault)
	}
}

func printUsageBar(x, y, w int, usage, scale uint64, clr tm.Attribute) {
	barw := int(math.Log10(float64(uint64((usage + scale - 1) / (scale / 10)))))
	if barw > w {
		barw = w
	} else if barw < 0 {
		barw = 0
	}
	for j := 0; j < w; j++ {
		tm.SetCell(x+j, y, symbols[box3], clr, tm.ColorDefault)
	}
	for j := 0; j < barw; j++ {
		tm.SetCell(x+j, y, symbols[box3], clr|tm.AttrBold, clr)
	}
}

func printDivider() {
	ui.printMsg("-----------------------------------------------------------")
}
func printDivider2() {
	ui.printMsg("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
}

type ethrUi interface {
	fini()
	printMsg(format string, a ...interface{})
	printErr(format string, a ...interface{})
	printDbg(format string, a ...interface{})
	paint()
	emitTestHdr()
	emitLatencyHdr()
	emitLatencyResults(remote, proto string, avg, min, max, p50, p90, p95, p99, p999, p9999 time.Duration)
	emitTestResultBegin()
	emitTestResult(s *ethrSession, proto EthrProtocol)
	printTestResults(s []string)
	emitTestResultEnd()
	emitStats(ethrNetStat)
}

var ui ethrUi
