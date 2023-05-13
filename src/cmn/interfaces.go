// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package cmn

type LineWriter interface {
	WriteLine(value string) (int, error)
}

type Widther interface {
	Width(s string) int
}
