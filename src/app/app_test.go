// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package app

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/keewek/ansible-pretty-print/src/cmn"
	"github.com/keewek/ansible-pretty-print/src/cmn/tst"
)

var fnTermSize = func(cols int, lines int, err error) TermSizeFunc {

	return func() (int, int, error) {
		return cols, lines, err
	}
}

func TestRun(t *testing.T) {

	t.Run("scanner error", func(t *testing.T) {
		var (
			out    cmn.LineBuilder
			outErr cmn.LineBuilder
		)

		c := &Config{
			TermWidth: DefaultTermWidth,
			Out:       &out,
			OutErr:    &outErr,
			Filepath:  "testdata/file-not-found.txt",
		}

		c.Init(fnTermSize(80, 0, nil))
		r := Run(c)

		want := "Exit code: 1"
		got := fmt.Sprintf("Exit code: %v", r)

		tst.DiffError(t, want, got)

	})

	t.Run("processor error", func(t *testing.T) {
		var (
			out    cmn.LineBuilder
			outErr cmn.LineBuilder
		)

		c := &Config{
			TermWidth: DefaultTermWidth,
			Out:       &out,
			OutErr:    &outErr,
			Filepath:  "testdata/list-tasks-err-play.txt",
		}

		c.Init(fnTermSize(80, 0, nil))
		r := Run(c)

		want := "Exit code: 1"
		got := fmt.Sprintf("Exit code: %v", r)

		tst.DiffError(t, want, got)

	})

	t.Run("usage", func(t *testing.T) {
		var (
			out     cmn.LineBuilder
			outErr  cmn.LineBuilder
			outFlag cmn.LineBuilder
			lb      cmn.LineBuilder
		)

		flag.CommandLine.SetOutput(&outFlag)
		defer func() {
			flag.CommandLine.SetOutput(nil)
		}()

		c := &Config{
			TermWidth: DefaultTermWidth,
			Out:       &out,
			OutErr:    &outErr,
			Filepath:  "",
		}

		c.Init(fnTermSize(80, 0, nil))
		r := Run(c)

		lb.WriteLine(fmt.Sprintf("Usage: %v [OPTION]... [FILE]", os.Args[0]))
		lb.WriteLine("Pretty-print Ansible's --list-tasks output")
		lb.WriteLine("")
		lb.WriteString(outFlag.String())
		outErr.WriteString(outFlag.String())

		want := lb.String()
		got := outErr.String()

		tst.DiffError(t, string(want), got)

		want = "Exit code: 0"
		got = fmt.Sprintf("Exit code: %v", r)

		tst.DiffError(t, want, got)

	})

	t.Run("output", func(t *testing.T) {
		var out cmn.LineBuilder
		// var outErr cmn.LineBuilder

		type testItem struct {
			name     string
			isMono   bool
			isTable  bool
			isDos    bool
			isStats  bool
			isChop   bool
			isIndent bool
		}

		tests := []testItem{
			{name: "runes"},
			{name: "runes-indent", isIndent: true},
			{name: "runes-chop_80", isChop: true},
			{name: "runes-chop_80-indent", isChop: true, isIndent: true},
			{name: "runes-table_80-ascii", isTable: true},
			{name: "runes-table_80-dos", isTable: true, isDos: true},
			{name: "mono", isMono: true},
			{name: "mono-indent", isIndent: true, isMono: true},
			{name: "mono-chop_80", isChop: true, isMono: true},
			{name: "mono-chop_80-indent", isChop: true, isIndent: true, isMono: true},
			{name: "mono-table_80-ascii", isTable: true, isMono: true},
			{name: "mono-table_80-dos", isTable: true, isDos: true, isMono: true},
		}

		for _, ti := range tests {

			file := "testdata/out-" + ti.name + ".txt"

			t.Run(ti.name, func(t *testing.T) {
				c := &Config{
					TermWidth: DefaultTermWidth,
					// Widther:   cmn.MonospaceWidther{},
					IsMono:   ti.isMono,
					IsTable:  ti.isTable,
					IsDos:    ti.isDos,
					IsStats:  ti.isStats,
					IsChop:   ti.isChop,
					IsIndent: ti.isIndent,
					Out:      &out,
					OutErr:   os.Stderr,
					Filepath: "testdata/list-tasks-1.txt",
				}

				c.Init(fnTermSize(80, 0, nil))
				Run(c)

				want, err := os.ReadFile(file)
				if err != nil {
					t.Fatal(err)
				}
				got := out.String()

				tst.DiffError(t, string(want), got)

				out.Reset()
			})
		}
	})

	t.Run("stats", func(t *testing.T) {
		var out cmn.LineBuilder
		// var outErr cmn.LineBuilder

		type testItem struct {
			name   string
			isMono bool
			// isTable  bool
			isDos   bool
			isStats bool
			// isChop   bool
			// isIndent bool
		}

		tests := []testItem{
			{name: "runes-ascii", isStats: true},
			{name: "runes-dos", isStats: true, isDos: true},
			{name: "mono-ascii", isStats: true, isMono: true},
			{name: "mono-dos", isStats: true, isMono: true, isDos: true},
		}

		for _, ti := range tests {

			file := "testdata/out-stats-" + ti.name + ".txt"

			t.Run(ti.name, func(t *testing.T) {
				c := &Config{
					TermWidth: DefaultTermWidth,
					// Widther:   cmn.MonospaceWidther{},
					IsMono:   ti.isMono,
					IsDos:    ti.isDos,
					IsStats:  ti.isStats,
					Out:      &out,
					OutErr:   os.Stderr,
					Filepath: "testdata/list-tasks-stats.txt",
				}

				c.Init(fnTermSize(80, 0, nil))
				Run(c)

				want, err := os.ReadFile(file)
				if err != nil {
					t.Fatal(err)
				}
				got := out.String()

				tst.DiffError(t, string(want), got)

				out.Reset()
			})
		}
	})

}
