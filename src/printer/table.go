// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/keewek/ansible-pretty-print/src/cmn"
	"github.com/keewek/ansible-pretty-print/src/processor"
)

type TablePrinter struct {
	fnChopMarkLine func(line string, maxWidth int, chopMark string) string
	widther        cmn.Widther
	indentPlay     int
	indentTask     int
	padPlay        string
	padTask        string
	maxLineWidth   int
	box            cmn.BoxChars
}

type tableWidth struct {
	block int
	name  int
	tags  int
}

func NewTablePrinter() *TablePrinter {
	return &TablePrinter{
		fnChopMarkLine: cmn.ChopMarkLine,
		widther:        cmn.RunesWidther{},
		indentPlay:     defaultIndentPlay,
		indentTask:     defaultIndentTask,
		padPlay:        strings.Repeat(" ", defaultIndentPlay),
		padTask:        strings.Repeat(" ", defaultIndentTask),
		box:            cmn.BoxCharsAscii(),
	}
}

func (tp *TablePrinter) SetWidther(value cmn.Widther) *TablePrinter {
	tp.widther = value
	tp.fnChopMarkLine = cmn.ChopMarkLineSelector(value)

	return tp
}

func (tp *TablePrinter) SetMaxLineWidth(value int) *TablePrinter {
	tp.maxLineWidth = value

	return tp
}

func (tp *TablePrinter) SetBoxChars(value cmn.BoxChars) *TablePrinter {
	tp.box = value

	return tp
}

func (tp *TablePrinter) makeBorders(w *tableWidth) (top string, middle string, bottom string) {

	block := strings.Repeat(tp.box.Hor, w.block+2)
	name := strings.Repeat(tp.box.Hor, w.name+2)
	tags := strings.Repeat(tp.box.Hor, w.tags+2)

	top = fmt.Sprint(tp.padTask, tp.box.CornerTL, block, tp.box.Top, name, tp.box.Top, tags, tp.box.CornerTR)
	middle = fmt.Sprint(tp.padTask, tp.box.Left, block, tp.box.Cross, name, tp.box.Cross, tags, tp.box.Right)
	bottom = fmt.Sprint(tp.padTask, tp.box.CornerBL, block, tp.box.Bottom, name, tp.box.Bottom, tags, tp.box.CornerBR)

	return top, middle, bottom
}

func (tp *TablePrinter) fitTable(s *processor.Stats) *tableWidth {
	width := &tableWidth{}

	fit := func(dW, width, minWidth int) (fit int, r int) {
		fit = cmn.Max(minWidth, width-dW)
		r = dW - (width - fit)

		return fit, r
	}

	minBlock := len("Block")
	minName := len("Name")
	minTags := len("Tags")

	width.block = cmn.Max(s.LongestTaskBlockLength, minBlock)
	width.name = cmn.Max(s.LongestTaskNameLength, minName)
	width.tags = cmn.Max(s.LongestTaskTagsLength, minTags)

	borderTop, _, _ := tp.makeBorders(width)
	lineWidth := tp.widther.Width(borderTop)

	if lineWidth > tp.maxLineWidth {

		dW := lineWidth - tp.maxLineWidth
		width.tags, dW = fit(dW, width.tags, minTags)

		if dW > 0 {
			width.name, dW = fit(dW, width.name, minName)

			if dW > 0 {
				width.block, _ = fit(dW, width.block, minBlock)
			}
		}

	}

	return width
}

func (tp *TablePrinter) printTable(output io.Writer, t *processor.Tasks, s *processor.Stats) {
	width := tp.fitTable(s)
	borderTop, borderMiddle, borderBottom := tp.makeBorders(width)

	fnPrintHeader := func() {
		block := fmt.Sprintf("%-*s", width.block, "Block")
		name := fmt.Sprintf("%-*s", width.name, "Name")
		tags := fmt.Sprintf("%-*s", width.tags, "Tags")

		tp.printLine(output, borderTop)
		tp.printLine(output, fmt.Sprintf("%[1]s%[2]s %[3]s %[2]s %[4]s %[2]s %[5]s %[2]s", tp.padTask, tp.box.Ver, block, name, tags))
		tp.printLine(output, borderMiddle)
	}

	fnPrintRow := func(t *processor.Task) {
		block := cmn.PadLeftFunc(tp.fnChopMarkLine(t.Block, width.block, "▒"), ' ', width.block, tp.widther.Width)
		name := cmn.PadRightFunc(tp.fnChopMarkLine(t.Name, width.name, "▒"), ' ', width.name, tp.widther.Width)
		tags := cmn.PadRightFunc(tp.fnChopMarkLine(t.Tags, width.tags, "▒"), ' ', width.tags, tp.widther.Width)

		tp.printLine(output, fmt.Sprintf("%[1]s%[2]s %[3]s %[2]s %[4]s %[2]s %[5]s %[2]s", tp.padTask, tp.box.Ver, block, name, tags))
	}

	fnPrintHeader()
	for _, task := range t.Tasks {
		fnPrintRow(task)
	}
	tp.printLine(output, borderBottom)
}

func (tp *TablePrinter) printLine(output io.Writer, value string) {
	fmt.Fprintln(output, tp.fnChopMarkLine(value, tp.maxLineWidth, "▒"))
}

func (tp *TablePrinter) PrintTo(output io.Writer, data *processor.Result) {
	var line string

	padPlay := strings.Repeat(" ", tp.indentPlay)

	for _, row := range data.Rows {

		switch t := row.Data.(type) {
		case *processor.Play:
			line = fmt.Sprintf("%s%s    TAGS: %s", padPlay, t.Description(), t.Tags)
			tp.printLine(output, line)

		case *processor.Tasks:
			tp.printTable(output, t, data.Stats)

		default:
			tp.printLine(output, t.String())
		}

	}
}

func (tp *TablePrinter) Print(data *processor.Result) {
	tp.PrintTo(os.Stdout, data)
}
