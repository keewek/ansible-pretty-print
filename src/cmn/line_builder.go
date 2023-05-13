// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package cmn

import (
	"strings"

	"github.com/keewek/ansible-pretty-print/src/cmn/eol"
)

type LineBuilder struct {
	strings.Builder
	EoL eol.EndOfLine
}

func (lb *LineBuilder) WriteLine(value string) (int, error) {
	len1, err := lb.WriteString(value)
	if err != nil {
		return len1, err
	}

	len2, err := lb.WriteString(lb.EoL.String())

	return len1 + len2, err
}
