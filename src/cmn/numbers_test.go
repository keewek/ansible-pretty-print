// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package cmn

import (
	"testing"

	"github.com/keewek/ansible-pretty-print/src/cmn/tst"
)

func TestMax(t *testing.T) {
	tests := []struct {
		x    int
		y    int
		want int
	}{
		{1, 2, 2},
		{2, 1, 2},
		{2, 2, 2},
		{-2, 0, 0},
		{-2, -4, -2},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Max(tt.x, tt.y)
			tst.DiffError(t, tt.want, got)
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		x    int
		y    int
		want int
	}{
		{1, 2, 1},
		{2, 1, 1},
		{2, 2, 2},
		{-2, 0, -2},
		{-2, -4, -4},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Min(tt.x, tt.y)
			tst.DiffError(t, tt.want, got)
		})
	}
}
