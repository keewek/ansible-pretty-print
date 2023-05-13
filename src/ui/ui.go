// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/keewek/ansible-pretty-print/src/cmn"
)

var (
	Box     cmn.BoxChars  = cmn.BoxCharsAscii()
	FnWidth cmn.WidthFunc = cmn.WidthRunes
)

func MsgBoxFunc(w io.Writer, input []string, box *cmn.BoxChars, fnWidth cmn.WidthFunc) {
	padding := 1
	maxLen := 0
	borderLen := 0

	for _, item := range input {
		itemLen := fnWidth(item)
		if itemLen > maxLen {
			maxLen = itemLen
		}
	}

	borderLen = maxLen + padding*2
	border := strings.Repeat(box.Hor, borderLen)
	pad := strings.Repeat(" ", padding)

	fmt.Fprint(w, box.CornerTL, border, box.CornerTR, "\n")

	for _, item := range input {
		val := cmn.PadRightFunc(item, ' ', maxLen, fnWidth)
		fmt.Fprint(w, box.Ver, pad, val, pad, box.Ver, "\n")
	}

	fmt.Fprint(w, box.CornerBL, border, box.CornerBR, "\n")
}

func MsgBoxTo(output io.Writer, input []string) {
	MsgBoxFunc(output, input, &Box, FnWidth)
}

func MsgBox(input []string) {
	MsgBoxFunc(os.Stdout, input, &Box, FnWidth)
}
