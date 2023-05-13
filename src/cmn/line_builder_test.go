// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package cmn

import (
	"strings"
	"testing"

	"github.com/keewek/ansible-pretty-print/src/cmn/eol"
	"github.com/keewek/ansible-pretty-print/src/cmn/tst"
)

func Test_LineBuilder(t *testing.T) {

	t.Run("WriteLine", func(t *testing.T) {
		var b strings.Builder
		var lb LineBuilder

		b.WriteString("1\n")
		b.WriteString("2\n")
		b.WriteString("3\n")

		lb.WriteLine("1")
		lb.WriteLine("2")
		lb.WriteLine("3")

		want := b.String()
		got := lb.String()

		tst.DiffError(t, want, got)

	})

	t.Run("EoL", func(t *testing.T) {

		t.Run("CRLF", func(t *testing.T) {
			var b strings.Builder
			var lb LineBuilder

			b.WriteString("1\r\n")
			b.WriteString("2\r\n")
			b.WriteString("3\r\n")

			lb.EoL = eol.CrLf()
			lb.WriteLine("1")
			lb.WriteLine("2")
			lb.WriteLine("3")

			want := b.String()
			got := lb.String()

			tst.DiffError(t, want, got)
		})

		t.Run("Custom", func(t *testing.T) {
			var b strings.Builder
			var lb LineBuilder

			b.WriteString("1###\n")
			b.WriteString("2###\n")
			b.WriteString("3###\n")

			lb.EoL = "###\n"
			lb.WriteLine("1")
			lb.WriteLine("2")
			lb.WriteLine("3")

			want := b.String()
			got := lb.String()

			tst.DiffError(t, want, got)
		})

	})
}
