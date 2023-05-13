// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"testing"

	"github.com/keewek/ansible-pretty-print/src/cmn"
	"github.com/keewek/ansible-pretty-print/src/cmn/tst"
	"github.com/keewek/ansible-pretty-print/src/processor"
)

type blockquote string

func (b blockquote) String() string {
	return "| " + string(b)
}

func TestNewColumnPrinter(t *testing.T) {
	cp := &ColumnPrinter{
		widther:         cmn.RunesWidther{},
		columnSeparator: defaultColumnSeparator,
		blockSeparator:  defaultBlockSeparator,
		indentPlay:      defaultIndentPlay,
		indentTask:      defaultIndentTask,
	}

	r := NewColumnPrinter()

	want := fmt.Sprintf("%#v", cp)
	got := fmt.Sprintf("%#v", r)

	tst.DiffError(t, want, got)
}

func Test_ColumnPrinterSetWidther(t *testing.T) {
	w := cmn.MonospaceWidther{}
	cp := NewColumnPrinter()

	cp.SetWidther(w)

	want := fmt.Sprintf("%#v", w)
	got := fmt.Sprintf("%#v", cp.widther)

	tst.DiffError(t, want, got)
}

func Test_ColumnPrinterSetIsIndentBlock(t *testing.T) {
	isIndentBlock := true
	cp := NewColumnPrinter()

	cp.SetIsIndentBlock(isIndentBlock)

	want := isIndentBlock
	got := cp.isIndentBlock

	tst.DiffError(t, want, got)
}

func Test_ColumnPrinterSetIsChopLines(t *testing.T) {
	isChopLines := true
	cp := NewColumnPrinter()

	cp.SetIsChopLines(isChopLines)

	want := isChopLines
	got := cp.isChopLines

	tst.DiffError(t, want, got)
}

func Test_ColumnPrinterSetMaxLineWidth(t *testing.T) {
	maxLineWidth := 40
	cp := NewColumnPrinter()

	cp.SetMaxLineWidth(maxLineWidth)

	want := maxLineWidth
	got := cp.maxLineWidth

	tst.DiffError(t, want, got)
}

func Test_ColumnPrinter_calcCol1Witdh(t *testing.T) {

	tests := []struct {
		name           string
		indentPlay     int
		indentTask     int
		blockSeparator string
		isIndentBlock  bool
		stats          processor.Stats
		want           int
	}{
		{
			name:           "default---zero-value-stats---indent_false",
			indentPlay:     defaultIndentPlay,
			indentTask:     defaultIndentTask,
			blockSeparator: defaultBlockSeparator,
			isIndentBlock:  false,
			stats: processor.Stats{
				LongestPlayDescriptionLength: 0,
				LongestTaskBlockLength:       0,
				LongestTaskNameLength:        0,
				LongestTaskDescriptionLength: 0,
			},
			want: 6,
		},
		{
			name:           "default---zero-value-stats---indent_true",
			indentPlay:     defaultIndentPlay,
			indentTask:     defaultIndentTask,
			blockSeparator: defaultBlockSeparator,
			isIndentBlock:  true,
			stats: processor.Stats{
				LongestPlayDescriptionLength: 0,
				LongestTaskBlockLength:       0,
				LongestTaskNameLength:        0,
				LongestTaskDescriptionLength: 0,
			},
			want: 8,
		},
		{
			name:           "",
			indentPlay:     defaultIndentPlay,
			indentTask:     defaultIndentTask,
			blockSeparator: defaultBlockSeparator,
			isIndentBlock:  false,
			stats: processor.Stats{
				LongestPlayDescriptionLength: 40,
				LongestTaskBlockLength:       10,
				LongestTaskNameLength:        18,
				LongestTaskDescriptionLength: 30,
			},
			want: 42,
		},
		{
			name:           "",
			indentPlay:     defaultIndentPlay,
			indentTask:     defaultIndentTask,
			blockSeparator: defaultBlockSeparator,
			isIndentBlock:  false,
			stats: processor.Stats{
				LongestPlayDescriptionLength: 30,
				LongestTaskBlockLength:       10,
				LongestTaskNameLength:        18,
				LongestTaskDescriptionLength: 30,
			},
			want: 36,
		},
		{
			name:           "",
			indentPlay:     defaultIndentPlay,
			indentTask:     defaultIndentTask,
			blockSeparator: defaultBlockSeparator,
			isIndentBlock:  true,
			stats: processor.Stats{
				LongestPlayDescriptionLength: 30,
				LongestTaskBlockLength:       10,
				LongestTaskNameLength:        18,
				LongestTaskDescriptionLength: 30,
			},
			want: 36,
		},
		{
			name:           "",
			indentPlay:     defaultIndentPlay,
			indentTask:     defaultIndentTask,
			blockSeparator: defaultBlockSeparator,
			isIndentBlock:  false,
			stats: processor.Stats{
				LongestPlayDescriptionLength: 41,
				LongestTaskBlockLength:       5,
				LongestTaskNameLength:        70,
				LongestTaskDescriptionLength: 11,
			},
			want: 43,
		},
		{
			name:           "",
			indentPlay:     defaultIndentPlay,
			indentTask:     defaultIndentTask,
			blockSeparator: defaultBlockSeparator,
			isIndentBlock:  true,
			stats: processor.Stats{
				LongestPlayDescriptionLength: 41,
				LongestTaskBlockLength:       5,
				LongestTaskNameLength:        70,
				LongestTaskDescriptionLength: 11,
			},
			want: 83,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := NewColumnPrinter()
			cp.indentPlay = tt.indentPlay
			cp.indentTask = tt.indentTask
			cp.blockSeparator = tt.blockSeparator
			cp.isIndentBlock = tt.isIndentBlock

			got := cp.calcCol1Width(&tt.stats)

			tst.DiffError(t, tt.want, got)
		})
	}
}

func Test_ColumnPrinterPrintTo(t *testing.T) {

	t.Run("row is fmt.Stringer", func(t *testing.T) {
		t.Run("nochop", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: blockquote("Indent 0")},
					{Indent: 2, Data: blockquote("Indent 2")},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("| Indent 0    ")
			lb.WriteLine("| Indent 2    ")

			cp := NewColumnPrinter()
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("chop maxLineWidth is 0", func(t *testing.T) {
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

			cp := NewColumnPrinter()
			cp.isChopLines = true
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("chop maxLineWidth is 5", func(t *testing.T) {
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

			cp := NewColumnPrinter()
			cp.isChopLines = true
			cp.maxLineWidth = 5
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})
	})

	t.Run("row is processor.Play", func(t *testing.T) {
		t.Run("nochop", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Play{Name: "play #1 (demo): Demo play", Tags: "[p1, demo]"}},
					{Indent: 0, Data: &processor.Play{Name: "play #2 (demo): Demo play", Tags: "[p2, demo]"}},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("  play #1 (demo): Demo play    TAGS: [p1, demo]")
			lb.WriteLine("  play #2 (demo): Demo play    TAGS: [p2, demo]")

			cp := NewColumnPrinter()
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("chop maxLineWidth is 0", func(t *testing.T) {
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

			cp := NewColumnPrinter()
			cp.isChopLines = true
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("chop maxLineWidth is 5", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Play{Name: "play #1 (demo): Demo play", Tags: "[p1, demo]"}},
					{Indent: 0, Data: &processor.Play{Name: "play #2 (demo): Demo play", Tags: "[p2, demo]"}},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("  pl▒")
			lb.WriteLine("  pl▒")

			cp := NewColumnPrinter()
			cp.isChopLines = true
			cp.maxLineWidth = 5
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})
	})

	t.Run("row is processor.Tasks", func(t *testing.T) {
		t.Run("nochop", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Tasks{PlayNumber: 1, Tasks: []*processor.Task{
						{Block: "Block 1", Name: "Task 1", Tags: "[t1]"},
						{Block: "Block 2", Name: "Task 2", Tags: "[t2]"},
					}}},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("      Block 1: Task 1    TAGS: [t1]")
			lb.WriteLine("      Block 2: Task 2    TAGS: [t2]")

			cp := NewColumnPrinter()
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("chop maxLineWidth is 0", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Tasks{PlayNumber: 1, Tasks: []*processor.Task{
						{Block: "Block 1", Name: "Task 1", Tags: "[t1]"},
						{Block: "Block 2", Name: "Task 2", Tags: "[t2]"},
					}}},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("")
			lb.WriteLine("")

			cp := NewColumnPrinter()
			cp.isChopLines = true
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("chop maxLineWidth is 5", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Tasks{PlayNumber: 1, Tasks: []*processor.Task{
						{Block: "Block 1", Name: "Task 1", Tags: "[t1]"},
						{Block: "Block 2", Name: "Task 2", Tags: "[t2]"},
					}}},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("    ▒")
			lb.WriteLine("    ▒")

			cp := NewColumnPrinter()
			cp.isChopLines = true
			cp.maxLineWidth = 5
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("noindent", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Tasks{PlayNumber: 1, Tasks: []*processor.Task{
						{Block: "", Name: "Task 1", Tags: "[t1]"},
						{Block: "Block 2", Name: "Task 2", Tags: "[t2]"},
					}}},
				},
				Stats: &processor.Stats{},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("      Task 1    TAGS: [t1]")
			lb.WriteLine("      Block 2: Task 2    TAGS: [t2]")

			cp := NewColumnPrinter()
			cp.isIndentBlock = false
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("indent", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: &processor.Tasks{PlayNumber: 1, Tasks: []*processor.Task{
						{Block: "", Name: "Task 1", Tags: "[t1]"},
						{Block: "Block 2", Name: "Task 2", Tags: "[t2]"},
					}}},
				},
				Stats: &processor.Stats{
					LongestTaskBlockLength: 7,
				},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("             : Task 1    TAGS: [t1]")
			lb.WriteLine("      Block 2: Task 2    TAGS: [t2]")

			cp := NewColumnPrinter()
			cp.isIndentBlock = true
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})
	})

	t.Run("mixed rows", func(t *testing.T) {
		t.Run("nochop", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: processor.Passthru("Passthru")},
					{Indent: 0, Data: &processor.Play{Name: "play #1 (demo): Demo play", Tags: "[p1, demo]"}},
					{Indent: 0, Data: blockquote("  blockquote")},
					{Indent: 0, Data: &processor.Tasks{PlayNumber: 1, Tasks: []*processor.Task{
						{Block: "Block 1", Name: "Task 1", Tags: "[t1]"},
						{Block: "Block 2", Name: "Task 2", Tags: "[t2]"},
					}}},
				},
				Stats: &processor.Stats{
					LongestPlayDescriptionLength: 25,
				},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("Passthru                       ")
			lb.WriteLine("  play #1 (demo): Demo play    TAGS: [p1, demo]")
			lb.WriteLine("|   blockquote                 ")
			lb.WriteLine("      Block 1: Task 1          TAGS: [t1]")
			lb.WriteLine("      Block 2: Task 2          TAGS: [t2]")

			cp := NewColumnPrinter()
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})

		t.Run("chop maxLineWidth is 10", func(t *testing.T) {
			r := processor.Result{
				Rows: []*processor.Row{
					{Indent: 0, Data: processor.Passthru("Passthru")},
					{Indent: 0, Data: &processor.Play{Name: "play #1 (demo): Demo play", Tags: "[p1, demo]"}},
					{Indent: 0, Data: blockquote("  blockquote")},
					{Indent: 0, Data: &processor.Tasks{PlayNumber: 1, Tasks: []*processor.Task{
						{Block: "Block 1", Name: "Task 1", Tags: "[t1]"},
						{Block: "Block 2", Name: "Task 2", Tags: "[t2]"},
					}}},
				},
				Stats: &processor.Stats{
					LongestPlayDescriptionLength: 25,
				},
			}

			var lb, out cmn.LineBuilder
			lb.WriteLine("Passthru ▒")
			lb.WriteLine("  play #1▒")
			lb.WriteLine("|   block▒")
			lb.WriteLine("      Blo▒")
			lb.WriteLine("      Blo▒")

			cp := NewColumnPrinter()
			cp.isChopLines = true
			cp.maxLineWidth = 10
			cp.PrintTo(&out, &r)

			want := lb.String()
			got := out.String()

			tst.DiffError(t, want, got)
		})
	})

}
