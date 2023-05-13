// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package app

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/keewek/ansible-pretty-print/src/cmn"
	"github.com/keewek/ansible-pretty-print/src/cmn/flags"
)

func fileComparer(x, y *os.File) bool {
	return x.Name() == y.Name()
}

func TestDefaultConfig(t *testing.T) {
	want := &Config{
		TermWidth: 80,
		Widther:   cmn.RunesWidther{},
		Out:       os.Stdout,
		OutErr:    os.Stderr,
	}

	got := DefaultConfig()

	if diff := cmp.Diff(want, got, cmp.Comparer(fileComparer)); diff != "" {
		t.Errorf("(-want +got): \n%s", diff)
	}

	// if diff := cmp.Diff(want, got, cmp.Exporter(func(rt reflect.Type) bool { return true })); diff != "" {
	// 	t.Errorf("(-want +got): \n%s", diff)
	// }
}

func Test_ConfigInit(t *testing.T) {

	fnTermSize := func(cols int, lines int, err error) TermSizeFunc {

		return func() (int, int, error) {
			return cols, lines, err
		}
	}

	t.Run("MonospaceWidther", func(t *testing.T) {
		c := &Config{
			IsMono: true,
		}

		c.Init(fnTermSize(0, 0, nil))
		want := "cmn.MonospaceWidther"
		got := fmt.Sprintf("%T", c.Widther)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("RunesWidther", func(t *testing.T) {
		c := &Config{
			IsMono: false,
		}

		c.Init(fnTermSize(0, 0, nil))
		want := "cmn.RunesWidther"
		got := fmt.Sprintf("%T", c.Widther)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("Term size", func(t *testing.T) {

		var lb cmn.LineBuilder

		lb.WriteLine("")
		lb.WriteLine("!!! Can't determine terminal width!")
		lb.WriteLine("!!!    [width: 0; err: Forced test error]")
		lb.WriteLine("!!!")
		lb.WriteLine("!!! Falling back to default width of 80 chars.")
		lb.WriteLine("!!! Use --width to specify custom value.")
		lb.WriteLine("!!!    e.g., --width $(tput cols)")
		lb.WriteLine("")

		errMsg := lb.String()

		tests := []struct {
			name       string
			want       int
			isChop     bool
			isTable    bool
			forceWidth bool
			fn         TermSizeFunc
			errMsg     string
		}{
			{"Default", 80, false, false, false, fnTermSize(0, 0, nil), ""},
			{"IsChop", 10, true, false, false, fnTermSize(10, 0, nil), ""},
			{"IsTable", 10, false, true, false, fnTermSize(10, 0, nil), ""},
			{"Default with width flag", 123, false, false, true, fnTermSize(10, 0, nil), ""},
			{"IsChop with width flag", 123, true, false, true, fnTermSize(10, 0, nil), ""},
			{"IsTable with width flag", 123, false, true, true, fnTermSize(10, 0, nil), ""},
			{"TermSizeError", 80, false, true, false, fnTermSize(0, 0, errors.New("Forced test error")), errMsg},
		}

		flag.Set(kFlagWidth, "123")

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				lb.Reset()

				c := &Config{
					TermWidth: 80,
					IsChop:    tt.isChop,
					IsTable:   tt.isTable,
					OutErr:    &lb,
				}

				if tt.forceWidth {
					flags.Enable(kFlagWidth)
				} else {
					flags.Disable(kFlagWidth)
				}

				c.ApplyFlags()
				c.Init(tt.fn)

				got := c.TermWidth

				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("(-want +got): \n%s", diff)
				}

				if diff := cmp.Diff(tt.errMsg, lb.String()); diff != "" {
					t.Errorf("(-want +got): \n%s", diff)
				}

			})
		}

	})
}

func Test_ConfigAcquireBoxChars(t *testing.T) {
	t.Run("Dos", func(t *testing.T) {
		c := &Config{
			IsDos: true,
		}

		want := cmn.BoxCharsDos()
		got := c.AcquireBoxChars()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("Ascii", func(t *testing.T) {
		c := &Config{
			IsDos: false,
		}

		want := cmn.BoxCharsAscii()
		got := c.AcquireBoxChars()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})
}

func Test_ConfigAcquirePrinter(t *testing.T) {
	t.Run("TablePrinter", func(t *testing.T) {
		c := &Config{
			IsTable: true,
		}

		p := c.AcquirePrinter()
		want := "*printer.TablePrinter"
		got := fmt.Sprintf("%T", p)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("ColumnPrinter", func(t *testing.T) {
		c := &Config{
			IsTable: false,
		}

		p := c.AcquirePrinter()
		want := "*printer.ColumnPrinter"
		got := fmt.Sprintf("%T", p)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})
}

func Test_ConfigAcquireScanner(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		c := &Config{}

		s, closer, err := c.AcquireScanner()

		defer closer()

		want := "==nil: s=true; err=true"
		got := fmt.Sprintf("==nil: s=%v; err=%v", s == nil, err == nil)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("Stdin", func(t *testing.T) {
		c := &Config{
			IsStdin:  true,
			Filepath: "testdata/test.txt",
		}

		s, closer, err := c.AcquireScanner()

		defer closer()

		want := "==nil: s=false; err=true"
		got := fmt.Sprintf("==nil: s=%v; err=%v", s == nil, err == nil)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("File", func(t *testing.T) {
		c := &Config{
			IsStdin:  false,
			Filepath: "testdata/test.txt",
		}

		s, closer, err := c.AcquireScanner()

		defer closer()

		want := "==nil: s=false; err=true"
		got := fmt.Sprintf("==nil: s=%v; err=%v", s == nil, err == nil)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("File not found", func(t *testing.T) {
		c := &Config{
			IsStdin:  false,
			Filepath: "testdata/file-not-found.txt",
		}

		s, closer, err := c.AcquireScanner()

		defer closer()

		want := "==nil: s=true; err=false"
		got := fmt.Sprintf("==nil: s=%v; err=%v", s == nil, err == nil)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})
}

func Test_ConfigApplyFlags(t *testing.T) {
	flag.Set(kFlagIsChop, "1")
	flag.Set(kFlagIsDos, "1")
	flag.Set(kFlagIsIndent, "1")
	flag.Set(kFlagIsMono, "1")
	flag.Set(kFlagIsStats, "1")
	flag.Set(kFlagIsStdin, "1")
	flag.Set(kFlagIsTable, "1")
	flag.Set(kFlagIsVersion, "1")
	flag.Set(kFlagWidth, "40")

	flags.EnableAll()

	want := &Config{
		IsChop:    true,
		IsDos:     true,
		IsIndent:  true,
		IsMono:    true,
		IsStats:   true,
		IsStdin:   true,
		IsTable:   true,
		IsVersion: true,
		TermWidth: 40,
		Widther:   nil,
	}

	got := &Config{}
	got.ApplyFlags()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got): \n%s", diff)
	}
}
