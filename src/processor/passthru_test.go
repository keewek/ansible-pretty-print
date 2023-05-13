// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package processor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Passthru(t *testing.T) {
	p := Passthru(" \t123 \n")

	t.Run("Implements Stringer interface", func(t *testing.T) {
		got := p.String()
		want := " \t123 \n"

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})
}
