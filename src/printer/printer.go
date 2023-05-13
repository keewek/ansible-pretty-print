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

const (
	defaultColumnSeparator = "    "
	defaultBlockSeparator  = ": "
	defaultIndentPlay      = 2
	defaultIndentTask      = 6
)

type ColumnPrinter struct {
	widther         cmn.Widther
	columnSeparator string
	blockSeparator  string
	indentPlay      int
	indentTask      int
	maxLineWidth    int
	isIndentBlock   bool
	isChopLines     bool
}

func NewColumnPrinter() *ColumnPrinter {
	return &ColumnPrinter{
		widther:         cmn.RunesWidther{},
		columnSeparator: defaultColumnSeparator,
		blockSeparator:  defaultBlockSeparator,
		indentPlay:      defaultIndentPlay,
		indentTask:      defaultIndentTask,
	}
}

func (cp *ColumnPrinter) calcCol1Width(stats *processor.Stats) int {
	var play, task int

	play = cp.indentPlay + stats.LongestPlayDescriptionLength
	if cp.isIndentBlock {
		task = cp.indentTask + stats.LongestTaskBlockLength + cp.widther.Width(cp.blockSeparator) + stats.LongestTaskNameLength
	} else {
		task = cp.indentTask + stats.LongestTaskDescriptionLength
	}

	return cmn.Max(play, task)
}

func (cp *ColumnPrinter) SetWidther(value cmn.Widther) *ColumnPrinter {
	cp.widther = value

	return cp
}

func (cp *ColumnPrinter) SetIsIndentBlock(value bool) *ColumnPrinter {
	cp.isIndentBlock = value

	return cp
}

func (cp *ColumnPrinter) SetIsChopLines(value bool) *ColumnPrinter {
	cp.isChopLines = value

	return cp
}

func (cp *ColumnPrinter) SetMaxLineWidth(value int) *ColumnPrinter {
	cp.maxLineWidth = value

	return cp
}

func (cp *ColumnPrinter) PrintTo(output io.Writer, data *processor.Result) {
	var (
		col1, col2  string
		fnPrintLine func(line string)
		fnFormCol1  func(task *processor.Task) string
	)

	padPlay := strings.Repeat(" ", cp.indentPlay)
	padTask := strings.Repeat(" ", cp.indentTask)

	fnFormLine := func(col1 string, col2 string) string {
		col1Width := cp.calcCol1Width(data.Stats)
		col1Padded := cmn.PadRightFunc(col1, ' ', col1Width, cp.widther.Width)

		return fmt.Sprint(col1Padded, cp.columnSeparator, col2)
	}

	if cp.isIndentBlock {
		fnFormCol1 = func(t *processor.Task) string {
			blockPadded := cmn.PadLeftFunc(t.Block, ' ', data.Stats.LongestTaskBlockLength, cp.widther.Width)

			return fmt.Sprint(padTask, blockPadded, cp.blockSeparator, t.Name)
		}

	} else {
		fnFormCol1 = func(t *processor.Task) string {
			return padTask + t.Description()
		}
	}

	if cp.isChopLines {
		fnChopMarkLine := cmn.ChopMarkLineSelector(cp.widther)

		fnPrintLine = func(line string) {
			fmt.Fprintln(output, fnChopMarkLine(line, cp.maxLineWidth, "▒"))
		}

	} else {
		fnPrintLine = func(line string) {
			fmt.Fprintln(output, line)
		}
	}

	for _, row := range data.Rows {

		switch t := row.Data.(type) {
		case *processor.Play:
			col1 = padPlay + t.Name
			col2 = "TAGS: " + t.Tags
			fnPrintLine(fnFormLine(col1, col2))

		case *processor.Tasks:
			for _, task := range t.Tasks {
				col1 = fnFormCol1(task)
				col2 = "TAGS: " + task.Tags
				fnPrintLine(fnFormLine(col1, col2))
			}

		default:
			col1 = t.String()
			col2 = ""
			fnPrintLine(fnFormLine(col1, col2))
		}

	}
}

func (cp *ColumnPrinter) Print(data *processor.Result) {
	cp.PrintTo(os.Stdout, data)
}
