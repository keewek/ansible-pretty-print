// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package cmn

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

type BoxChars struct {
	CornerTL string
	CornerTR string
	CornerBL string
	CornerBR string
	Left     string
	Right    string
	Top      string
	Bottom   string
	Cross    string
	Hor      string
	Ver      string
}

func BoxCharsAscii() BoxChars {
	return BoxChars{
		CornerTL: "+",
		CornerTR: "+",
		CornerBL: "+",
		CornerBR: "+",
		Left:     "+",
		Right:    "+",
		Top:      "+",
		Bottom:   "+",
		Cross:    "+",
		Hor:      "-",
		Ver:      "|",
	}
}

func BoxCharsDos() BoxChars {
	return BoxChars{
		CornerTL: "┌",
		CornerTR: "┐",
		CornerBL: "└",
		CornerBR: "┘",
		Left:     "├",
		Right:    "┤",
		Top:      "┬",
		Bottom:   "┴",
		Cross:    "┼",
		Hor:      "─",
		Ver:      "│",
	}
}

// ---

type WidthFunc func(string) int

// WidthBytes returns the number of bytes in s
//
// Equivalent to len(s)
func WidthBytes(s string) int {
	return len(s)
}

// WidthRunes returns the number of runes in s
//
// Equivalent to utf8.RuneCountInString(s)
func WidthRunes(s string) int {
	return utf8.RuneCountInString(s)
}

// WidthMonospace returns the monospace width for the given string, that is, the number of same-size cells to be
// occupied by the string.
//
// Equivalent to uniseg.StringWidth(s)
func WidthMonospace(s string) int {
	return uniseg.StringWidth(s)
}

// ---

func PadLeftFunc(s string, padWith rune, maxWidth int, fnWidth WidthFunc) string {
	dW := maxWidth - fnWidth(s)

	if dW <= 0 {
		return s
	}

	pad := strings.Repeat(string(padWith), dW)
	return pad + s
}

func PadLeft(s string, padWith rune, maxRunes int) string {
	return PadLeftFunc(s, padWith, maxRunes, utf8.RuneCountInString)
}

func PadRightFunc(s string, padWith rune, maxWidth int, fnWidth WidthFunc) string {
	dW := maxWidth - fnWidth(s)

	if dW <= 0 {
		return s
	}

	pad := strings.Repeat(string(padWith), dW)
	return s + pad
}

func PadRight(s string, padWith rune, maxRunes int) string {
	return PadRightFunc(s, padWith, maxRunes, utf8.RuneCountInString)
}

// ---

type BytesWidther struct{}
type RunesWidther struct{}
type MonospaceWidther struct{}

func (bw BytesWidther) Width(s string) int {
	return len(s)
}

func (bw BytesWidther) String() string {
	return "BytesWidther{}"
}

func (rw RunesWidther) Width(s string) int {
	return utf8.RuneCountInString(s)
}

func (rw RunesWidther) String() string {
	return "RunesWidther{}"
}

func (mw MonospaceWidther) Width(s string) int {
	return uniseg.StringWidth(s)
}

func (mw MonospaceWidther) String() string {
	return "MonospaceWidther{}"
}

// ---

type FnChopMarkLine func(line string, maxWidth int, chopMark string) string

func ChopLineFunc(line string, maxWidth int, fnWidth WidthFunc) string {
	result := line

	if maxWidth <= 0 {
		result = ""
	} else if fnWidth(line) > maxWidth {

		var b strings.Builder
		width := 0

		for _, r := range line {
			width += fnWidth(string(r))
			if width > maxWidth {
				break
			}
			b.WriteRune(r)
		}

		result = b.String()

	}

	return result
}

func ChopMarkLineFunc(line string, maxWidth int, chopMark string, fnWidth WidthFunc) string {
	result := line

	if maxWidth <= 0 {
		result = ""
	} else if fnWidth(line) > maxWidth {

		if chopMark != "" {
			chopRune, _ := utf8.DecodeRuneInString(chopMark)
			chopMark = string(chopRune)
			chopMarkWidth := fnWidth(chopMark)
			if chopMarkWidth > maxWidth {
				panic(fmt.Errorf("cmn: can't fit `chopMark` with width %d when `maxWidth` is %d", chopMarkWidth, maxWidth))
			}
			maxWidth -= chopMarkWidth
		}

		result = ChopLineFunc(line, maxWidth, fnWidth) + chopMark

	}

	return result
}

func ChopMarkLine(line string, maxWidth int, chopMark string) string {
	result := line

	if maxWidth <= 0 {
		result = ""
	} else if WidthRunes(line) > maxWidth {

		if chopMark != "" {
			result = fmt.Sprintf("%.*s%.1s", maxWidth-1, line, chopMark)
		} else {
			result = fmt.Sprintf("%.*s", maxWidth, line)
		}

	}

	return result
}

func ChopMarkLineSelector(w Widther) (fn FnChopMarkLine) {

	switch w.(type) {
	case RunesWidther:
		fn = ChopMarkLine
	default:
		fn = func(line string, maxWidth int, chopMark string) string {
			return ChopMarkLineFunc(line, maxWidth, chopMark, w.Width)
		}
	}

	return fn
}
