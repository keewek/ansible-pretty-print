// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/keewek/ansible-pretty-print/src/cmn"
	"github.com/keewek/ansible-pretty-print/src/cmn/tst"
	"github.com/keewek/ansible-pretty-print/src/processor"
)

func TestNewTablePrinter(t *testing.T) {
	wtp := &TablePrinter{
		fnChopMarkLine: cmn.ChopMarkLine,
		widther:        cmn.RunesWidther{},
		indentPlay:     defaultIndentPlay,
		indentTask:     defaultIndentTask,
		padPlay:        strings.Repeat(" ", defaultIndentPlay),
		padTask:        strings.Repeat(" ", defaultIndentTask),
		box:            cmn.BoxCharsAscii(),
	}

	gtp := NewTablePrinter()

	want := fmt.Sprintf("%#v", wtp)
	got := fmt.Sprintf("%#v", gtp)

	tst.DiffError(t, want, got)
}

func Test_TablePrinterSetWidther(t *testing.T) {

	w := cmn.MonospaceWidther{}

	tp := NewTablePrinter()
	tp.SetWidther(w)

	want := fmt.Sprintf("%#v", w)
	got := fmt.Sprintf("%#v", tp.widther)

	tst.DiffError(t, want, got)
}

func Test_TablePrinterSetMaxLineWidth(t *testing.T) {

	w := 42

	tp := NewTablePrinter()
	tp.SetMaxLineWidth(w)

	want := w
	got := tp.maxLineWidth

	tst.DiffError(t, want, got)
}

func Test_TablePrinterSetBoxChars(t *testing.T) {

	w := cmn.BoxCharsDos()

	tp := NewTablePrinter()
	tp.SetBoxChars(w)

	want := w
	got := tp.box

	tst.DiffError(t, want, got)
}

func Test_TablePrinter_makeBorders(t *testing.T) {

	t.Run("BoxCharsAscii", func(t *testing.T) {
		tests := []struct {
			// want   []string
			top    string
			middle string
			bottom string
			width  tableWidth
		}{
			{
				top:    "      +--+--+--+",
				middle: "      +--+--+--+",
				bottom: "      +--+--+--+",
				width:  tableWidth{0, 0, 0}},
			{
				top:    "      +---+----+-----+",
				middle: "      +---+----+-----+",
				bottom: "      +---+----+-----+",
				width:  tableWidth{1, 2, 3}},
		}

		for _, tt := range tests {
			t.Run("", func(t *testing.T) {
				tp := NewTablePrinter()

				top, middle, bottom := tp.makeBorders(&tt.width)

				tst.DiffError(t, tt.top, top)
				tst.DiffError(t, tt.middle, middle)
				tst.DiffError(t, tt.bottom, bottom)

			})
		}
	})

	t.Run("BoxCharsDos", func(t *testing.T) {
		tests := []struct {
			// want   []string
			top    string
			middle string
			bottom string
			width  tableWidth
		}{
			{
				top:    "      ┌──┬──┬──┐",
				middle: "      ├──┼──┼──┤",
				bottom: "      └──┴──┴──┘",
				width:  tableWidth{0, 0, 0}},
			{
				top:    "      ┌───┬────┬─────┐",
				middle: "      ├───┼────┼─────┤",
				bottom: "      └───┴────┴─────┘",
				width:  tableWidth{1, 2, 3}},
		}

		for _, tt := range tests {
			t.Run("", func(t *testing.T) {
				tp := NewTablePrinter()
				tp.SetBoxChars(cmn.BoxCharsDos())

				top, middle, bottom := tp.makeBorders(&tt.width)

				tst.DiffError(t, tt.top, top)
				tst.DiffError(t, tt.middle, middle)
				tst.DiffError(t, tt.bottom, bottom)

			})
		}
	})
}

func Test_TablePrinter_fitTable(t *testing.T) {
	tests := []struct {
		stats        *processor.Stats
		maxLineWidth int
		want         *tableWidth
	}{
		{
			stats:        &processor.Stats{},
			maxLineWidth: 80,
			want:         &tableWidth{block: 5, name: 4, tags: 4},
		},
		{
			stats: &processor.Stats{
				LongestTaskBlockLength: 10,
				LongestTaskNameLength:  20,
				LongestTaskTagsLength:  30,
			},
			maxLineWidth: 60,
			want:         &tableWidth{block: 10, name: 20, tags: 14},
		},
		{
			stats: &processor.Stats{
				LongestTaskBlockLength: 10,
				LongestTaskNameLength:  20,
				LongestTaskTagsLength:  30,
			},
			maxLineWidth: 30,
			want:         &tableWidth{block: 6, name: 4, tags: 4},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			tp := NewTablePrinter()
			tp.SetMaxLineWidth(tt.maxLineWidth)

			tw := tp.fitTable(tt.stats)

			want := tst.ReprValue(tt.want) // fmt.Sprintf("%#v", tt.want)
			got := tst.ReprValue(tw)       // fmt.Sprintf("%#v", tw)

			tst.DiffError(t, want, got)
		})
	}
}

func Test_TablePrinter_printLine(t *testing.T) {
	tests := []struct {
		maxLineWidth int
		value        string
		want         string
	}{
		{0, "12345", "\n"},
		{80, "12345", "12345\n"},
		{4, "12345", "123▒\n"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			var lb cmn.LineBuilder

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(tt.maxLineWidth)
			tp.printLine(&lb, tt.value)

			got := lb.String()

			tst.DiffError(t, tt.want, got)
		})
	}
}

func Test_TablePrinter_printTable(t *testing.T) {
	t.Run("BoxCharsAscii", func(t *testing.T) {
		t.Run("", func(t *testing.T) {
			var out, want cmn.LineBuilder

			tasks := &processor.Tasks{
				PlayNumber: 1,
				Tasks:      []*processor.Task{},
			}

			stats := &processor.Stats{}
			maxLineWidth := 80
			want.WriteLine("      +-------+------+------+")
			want.WriteLine("      | Block | Name | Tags |")
			want.WriteLine("      +-------+------+------+")
			want.WriteLine("      +-------+------+------+")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(maxLineWidth)
			tp.printTable(&out, tasks, stats)

			tst.DiffError(t, want.String(), out.String())
		})

		t.Run("", func(t *testing.T) {
			var out, want cmn.LineBuilder

			tasks := &processor.Tasks{
				PlayNumber: 1,
				Tasks:      []*processor.Task{{}},
			}

			stats := &processor.Stats{}
			maxLineWidth := 80
			want.WriteLine("      +-------+------+------+")
			want.WriteLine("      | Block | Name | Tags |")
			want.WriteLine("      +-------+------+------+")
			want.WriteLine("      |       |      |      |")
			want.WriteLine("      +-------+------+------+")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(maxLineWidth)
			tp.printTable(&out, tasks, stats)

			tst.DiffError(t, want.String(), out.String())
		})

		t.Run("", func(t *testing.T) {
			var out, want cmn.LineBuilder

			tasks := &processor.Tasks{
				PlayNumber: 1,
				Tasks: []*processor.Task{
					{
						Block: "0123456789",
						Name:  "0123456789",
						Tags:  "0123456789",
					},
					{
						Block: "ABCDEFGHIJ",
						Name:  "ABCDEFGHIJ",
						Tags:  "ABCDEFGHIJ",
					},
				},
			}

			stats := &processor.Stats{
				LongestTaskBlockLength: 10,
				LongestTaskNameLength:  10,
				LongestTaskTagsLength:  10,
			}

			maxLineWidth := 80
			want.WriteLine("      +------------+------------+------------+")
			want.WriteLine("      | Block      | Name       | Tags       |")
			want.WriteLine("      +------------+------------+------------+")
			want.WriteLine("      | 0123456789 | 0123456789 | 0123456789 |")
			want.WriteLine("      | ABCDEFGHIJ | ABCDEFGHIJ | ABCDEFGHIJ |")
			want.WriteLine("      +------------+------------+------------+")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(maxLineWidth)
			tp.printTable(&out, tasks, stats)

			tst.DiffError(t, want.String(), out.String())
		})

		t.Run("", func(t *testing.T) {
			var out, want cmn.LineBuilder

			tasks := &processor.Tasks{
				PlayNumber: 1,
				Tasks: []*processor.Task{
					{
						Block: "0123456789",
						Name:  "0123456789",
						Tags:  "0123456789",
					},
					{
						Block: "ABCDEFGHIJ",
						Name:  "ABCDEFGHIJ",
						Tags:  "ABCDEFGHIJ",
					},
				},
			}

			stats := &processor.Stats{
				LongestTaskBlockLength: 10,
				LongestTaskNameLength:  10,
				LongestTaskTagsLength:  10,
			}

			maxLineWidth := 40
			want.WriteLine("      +------------+------------+------+")
			want.WriteLine("      | Block      | Name       | Tags |")
			want.WriteLine("      +------------+------------+------+")
			want.WriteLine("      | 0123456789 | 0123456789 | 012▒ |")
			want.WriteLine("      | ABCDEFGHIJ | ABCDEFGHIJ | ABC▒ |")
			want.WriteLine("      +------------+------------+------+")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(maxLineWidth)
			tp.printTable(&out, tasks, stats)

			tst.DiffError(t, want.String(), out.String())
		})

		t.Run("", func(t *testing.T) {
			var out, want cmn.LineBuilder

			tasks := &processor.Tasks{
				PlayNumber: 1,
				Tasks: []*processor.Task{
					{
						Block: "0123456789",
						Name:  "0123456789",
						Tags:  "0123456789",
					},
					{
						Block: "ABCDEFGHIJ",
						Name:  "ABCDEFGHIJ",
						Tags:  "ABCDEFGHIJ",
					},
				},
			}

			stats := &processor.Stats{
				LongestTaskBlockLength: 10,
				LongestTaskNameLength:  10,
				LongestTaskTagsLength:  10,
			}

			maxLineWidth := 30
			want.WriteLine("      +--------+------+------+")
			want.WriteLine("      | Block  | Name | Tags |")
			want.WriteLine("      +--------+------+------+")
			want.WriteLine("      | 01234▒ | 012▒ | 012▒ |")
			want.WriteLine("      | ABCDE▒ | ABC▒ | ABC▒ |")
			want.WriteLine("      +--------+------+------+")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(maxLineWidth)
			tp.printTable(&out, tasks, stats)

			tst.DiffError(t, want.String(), out.String())
		})
	})

	t.Run("BoxCharsDos", func(t *testing.T) {
		t.Run("", func(t *testing.T) {
			var out, want cmn.LineBuilder

			tasks := &processor.Tasks{
				PlayNumber: 1,
				Tasks:      []*processor.Task{},
			}

			stats := &processor.Stats{}
			maxLineWidth := 80
			want.WriteLine("      ┌───────┬──────┬──────┐")
			want.WriteLine("      │ Block │ Name │ Tags │")
			want.WriteLine("      ├───────┼──────┼──────┤")
			want.WriteLine("      └───────┴──────┴──────┘")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(maxLineWidth)
			tp.SetBoxChars(cmn.BoxCharsDos())
			tp.printTable(&out, tasks, stats)

			tst.DiffError(t, want.String(), out.String())
		})

		t.Run("", func(t *testing.T) {
			var out, want cmn.LineBuilder

			tasks := &processor.Tasks{
				PlayNumber: 1,
				Tasks:      []*processor.Task{{}},
			}

			stats := &processor.Stats{}
			maxLineWidth := 80
			want.WriteLine("      ┌───────┬──────┬──────┐")
			want.WriteLine("      │ Block │ Name │ Tags │")
			want.WriteLine("      ├───────┼──────┼──────┤")
			want.WriteLine("      │       │      │      │")
			want.WriteLine("      └───────┴──────┴──────┘")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(maxLineWidth)
			tp.SetBoxChars(cmn.BoxCharsDos())
			tp.printTable(&out, tasks, stats)

			tst.DiffError(t, want.String(), out.String())
		})

		t.Run("", func(t *testing.T) {
			var out, want cmn.LineBuilder

			tasks := &processor.Tasks{
				PlayNumber: 1,
				Tasks: []*processor.Task{
					{
						Block: "0123456789",
						Name:  "0123456789",
						Tags:  "0123456789",
					},
					{
						Block: "ABCDEFGHIJ",
						Name:  "ABCDEFGHIJ",
						Tags:  "ABCDEFGHIJ",
					},
				},
			}

			stats := &processor.Stats{
				LongestTaskBlockLength: 10,
				LongestTaskNameLength:  10,
				LongestTaskTagsLength:  10,
			}

			maxLineWidth := 80
			want.WriteLine("      ┌────────────┬────────────┬────────────┐")
			want.WriteLine("      │ Block      │ Name       │ Tags       │")
			want.WriteLine("      ├────────────┼────────────┼────────────┤")
			want.WriteLine("      │ 0123456789 │ 0123456789 │ 0123456789 │")
			want.WriteLine("      │ ABCDEFGHIJ │ ABCDEFGHIJ │ ABCDEFGHIJ │")
			want.WriteLine("      └────────────┴────────────┴────────────┘")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(maxLineWidth)
			tp.SetBoxChars(cmn.BoxCharsDos())
			tp.printTable(&out, tasks, stats)

			tst.DiffError(t, want.String(), out.String())
		})

		t.Run("", func(t *testing.T) {
			var out, want cmn.LineBuilder

			tasks := &processor.Tasks{
				PlayNumber: 1,
				Tasks: []*processor.Task{
					{
						Block: "0123456789",
						Name:  "0123456789",
						Tags:  "0123456789",
					},
					{
						Block: "ABCDEFGHIJ",
						Name:  "ABCDEFGHIJ",
						Tags:  "ABCDEFGHIJ",
					},
				},
			}

			stats := &processor.Stats{
				LongestTaskBlockLength: 10,
				LongestTaskNameLength:  10,
				LongestTaskTagsLength:  10,
			}

			maxLineWidth := 40
			want.WriteLine("      ┌────────────┬────────────┬──────┐")
			want.WriteLine("      │ Block      │ Name       │ Tags │")
			want.WriteLine("      ├────────────┼────────────┼──────┤")
			want.WriteLine("      │ 0123456789 │ 0123456789 │ 012▒ │")
			want.WriteLine("      │ ABCDEFGHIJ │ ABCDEFGHIJ │ ABC▒ │")
			want.WriteLine("      └────────────┴────────────┴──────┘")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(maxLineWidth)
			tp.SetBoxChars(cmn.BoxCharsDos())
			tp.printTable(&out, tasks, stats)

			tst.DiffError(t, want.String(), out.String())
		})

		t.Run("", func(t *testing.T) {
			var out, want cmn.LineBuilder

			tasks := &processor.Tasks{
				PlayNumber: 1,
				Tasks: []*processor.Task{
					{
						Block: "0123456789",
						Name:  "0123456789",
						Tags:  "0123456789",
					},
					{
						Block: "ABCDEFGHIJ",
						Name:  "ABCDEFGHIJ",
						Tags:  "ABCDEFGHIJ",
					},
				},
			}

			stats := &processor.Stats{
				LongestTaskBlockLength: 10,
				LongestTaskNameLength:  10,
				LongestTaskTagsLength:  10,
			}

			maxLineWidth := 30
			want.WriteLine("      ┌────────┬──────┬──────┐")
			want.WriteLine("      │ Block  │ Name │ Tags │")
			want.WriteLine("      ├────────┼──────┼──────┤")
			want.WriteLine("      │ 01234▒ │ 012▒ │ 012▒ │")
			want.WriteLine("      │ ABCDE▒ │ ABC▒ │ ABC▒ │")
			want.WriteLine("      └────────┴──────┴──────┘")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(maxLineWidth)
			tp.SetBoxChars(cmn.BoxCharsDos())
			tp.printTable(&out, tasks, stats)

			tst.DiffError(t, want.String(), out.String())
		})
	})

}

func Test_TablePrinterPrintTo(t *testing.T) {

	t.Run("row is fmt.Stringer", func(t *testing.T) {

		t.Run("maxLineWidth is 0", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: blockquote("Indent 0")},
					{Indent: 2, Data: blockquote("Indent 2")},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("")
			lb.WriteLine("")

			tp := NewTablePrinter()
			tp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("maxLineWidth is 5", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: blockquote("Indent 0")},
					{Indent: 2, Data: blockquote("Indent 2")},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("| In▒")
			lb.WriteLine("| In▒")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(5)
			tp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

	})

	t.Run("row is processor.Play", func(t *testing.T) {

		t.Run("maxLineWidth is 0", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Play{Name: "play #1 (demo): Demo play", Tags: "[p1, demo]"}},
					{Indent: 0, Data: &processor.Play{Name: "play #2 (demo): Demo play", Tags: "[p2, demo]"}},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("")
			lb.WriteLine("")

			tp := NewTablePrinter()
			tp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("maxLineWidth is 30", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Play{Name: "play #1 (demo): Demo play", Tags: "[p1, demo]"}},
					{Indent: 0, Data: &processor.Play{Name: "play #2 (demo): Demo play", Tags: "[p2, demo]"}},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("  play #1 (demo): Demo play  ▒")
			lb.WriteLine("  play #2 (demo): Demo play  ▒")

			tp := NewTablePrinter()
			tp.maxLineWidth = 30
			tp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})
	})

	t.Run("row is processor.Tasks", func(t *testing.T) {
		t.Run("maxLineWidth is 0", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Tasks{
						PlayNumber: 1,
						Tasks: []*processor.Task{
							{Block: "0123456789", Name: "0123456789", Tags: "0123456789"},
							{Block: "ABCDEFGHIJ", Name: "ABCDEFGHIJ", Tags: "ABCDEFGHIJ"},
						},
					}},
				},
				Stats: &processor.Stats{
					LongestTaskBlockLength: 10,
					LongestTaskNameLength:  10,
					LongestTaskTagsLength:  10,
				},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("")
			lb.WriteLine("")
			lb.WriteLine("")
			lb.WriteLine("")
			lb.WriteLine("")
			lb.WriteLine("")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(0)
			tp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("maxLineWidth is 30", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Tasks{
						PlayNumber: 1,
						Tasks: []*processor.Task{
							{Block: "0123456789", Name: "0123456789", Tags: "0123456789"},
							{Block: "ABCDEFGHIJ", Name: "ABCDEFGHIJ", Tags: "ABCDEFGHIJ"},
						},
					}},
				},
				Stats: &processor.Stats{
					LongestTaskBlockLength: 10,
					LongestTaskNameLength:  10,
					LongestTaskTagsLength:  10,
				},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("      +--------+------+------+")
			lb.WriteLine("      | Block  | Name | Tags |")
			lb.WriteLine("      +--------+------+------+")
			lb.WriteLine("      | 01234▒ | 012▒ | 012▒ |")
			lb.WriteLine("      | ABCDE▒ | ABC▒ | ABC▒ |")
			lb.WriteLine("      +--------+------+------+")

			tp := NewTablePrinter()
			tp.SetMaxLineWidth(30)
			tp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})
	})

	t.Run("mixed rows", func(t *testing.T) {
		// t.Run("maxLineWidth is 30", func(t *testing.T) {
		r := processor.Result{
			Rows: []*processor.Row{
				{Indent: 0, Data: processor.Passthru("Passthru")},
				{Indent: 0, Data: &processor.Play{Name: "play #1 (demo): Demo play", Tags: "[p1, demo]"}},
				{Indent: 0, Data: blockquote("  blockquote")},
				{Indent: 0, Data: &processor.Tasks{
					PlayNumber: 1,
					Tasks: []*processor.Task{
						{Block: "0123456789", Name: "0123456789", Tags: "0123456789"},
						{Block: "ABCDEFGHIJ", Name: "ABCDEFGHIJ", Tags: "ABCDEFGHIJ"},
					},
				}},
			},
			Stats: &processor.Stats{
				LongestTaskBlockLength: 10,
				LongestTaskNameLength:  10,
				LongestTaskTagsLength:  10,
			},
		}

		var lb, out cmn.LineBuilder
		lb.WriteLine("Passthru")
		lb.WriteLine("  play #1 (demo): Demo play  ▒")
		lb.WriteLine("|   blockquote")
		lb.WriteLine("      +--------+------+------+")
		lb.WriteLine("      | Block  | Name | Tags |")
		lb.WriteLine("      +--------+------+------+")
		lb.WriteLine("      | 01234▒ | 012▒ | 012▒ |")
		lb.WriteLine("      | ABCDE▒ | ABC▒ | ABC▒ |")
		lb.WriteLine("      +--------+------+------+")

		tp := NewTablePrinter()
		tp.SetMaxLineWidth(30)
		tp.PrintTo(&out, &r)

		want := lb.String()
		got := out.String()

		tst.DiffError(t, want, got)
		// })
	})

}
