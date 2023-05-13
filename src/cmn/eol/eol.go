// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package eol

type EndOfLine string

const (
	CRLF string = "\r\n"
	LF   string = "\n"
)

func (e EndOfLine) String() string {
	switch e {
	case "":
		return LF
	}

	return string(e)
}

func CrLf() EndOfLine {
	return EndOfLine(CRLF)
}

func Lf() EndOfLine {
	return EndOfLine(LF)
}
