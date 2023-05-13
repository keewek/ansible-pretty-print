// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package cmn

import (
	"fmt"
	"testing"

	"github.com/keewek/ansible-pretty-print/src/cmn/tst"
)

func Test_WidthFunc(t *testing.T) {
	var fn WidthFunc = func(s string) int {
		return 0
	}

	_ = fn
}

func TestWidthBytes(t *testing.T) {
	tests := []struct {
		value string
		want  int
	}{
		{"Hello", 5},
		{"Привет", 12},
		{"你好", 6},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := WidthBytes(tt.value)
			tst.DiffError(t, tt.want, got)
		})
	}
}

func TestWidthRunes(t *testing.T) {
	tests := []struct {
		value string
		want  int
	}{
		{"Hello", 5},
		{"Привет", 6},
		{"你好", 2},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := WidthRunes(tt.value)
			tst.DiffError(t, tt.want, got)
		})
	}
}

func TestWidthMonospace(t *testing.T) {
	tests := []struct {
		value string
		want  int
	}{
		{"Hello", 5},
		{"Привет", 6},
		{"你好", 4},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := WidthMonospace(tt.value)
			tst.DiffError(t, tt.want, got)
		})
	}
}

func Test_BytesWidther(t *testing.T) {
	t.Run("Width", func(t *testing.T) {
		tests := []struct {
			value string
			want  int
		}{
			{"Hello", 5},
			{"Привет", 12},
			{"你好", 6},
		}

		for _, tt := range tests {
			t.Run("", func(t *testing.T) {
				got := BytesWidther{}.Width(tt.value)
				tst.DiffError(t, tt.want, got)
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		want := "BytesWidther{}"
		got := BytesWidther{}.String()

		tst.DiffError(t, want, got)
	})
}

func Test_RunesWidther(t *testing.T) {
	t.Run("Width", func(t *testing.T) {
		tests := []struct {
			value string
			want  int
		}{
			{"Hello", 5},
			{"Привет", 6},
			{"你好", 2},
		}

		for _, tt := range tests {
			t.Run("", func(t *testing.T) {
				got := RunesWidther{}.Width(tt.value)
				tst.DiffError(t, tt.want, got)
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		want := "RunesWidther{}"
		got := RunesWidther{}.String()

		tst.DiffError(t, want, got)
	})
}

func Test_MonospaceWidther(t *testing.T) {
	t.Run("Width", func(t *testing.T) {
		tests := []struct {
			value string
			want  int
		}{
			{"Hello", 5},
			{"Привет", 6},
			{"你好", 4},
		}

		for _, tt := range tests {
			t.Run("", func(t *testing.T) {
				got := MonospaceWidther{}.Width(tt.value)
				tst.DiffError(t, tt.want, got)
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		want := "MonospaceWidther{}"
		got := MonospaceWidther{}.String()

		tst.DiffError(t, want, got)
	})
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"Hello", "..........Hello"},
		{"Привет", ".........Привет"},
		{"你好", ".............你好"},
		{"###############", "###############"},
		{"################", "################"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := PadLeft(tt.value, '.', 15)
			tst.DiffError(t, tt.want, got)
		})
	}

}
func TestPadLeftFunc(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    string
		fnWidth WidthFunc
	}{
		{"WidthBytes", "Hello", "..........Hello", WidthBytes},
		{"WidthBytes", "Привет", "...Привет", WidthBytes},
		{"WidthBytes", "你好", ".........你好", WidthBytes},
		{"WidthRunes", "Hello", "..........Hello", WidthRunes},
		{"WidthRunes", "Привет", ".........Привет", WidthRunes},
		{"WidthRunes", "你好", ".............你好", WidthRunes},
		{"WidthMonospace", "Hello", "..........Hello", WidthMonospace},
		{"WidthMonospace", "Привет", ".........Привет", WidthMonospace},
		{"WidthMonospace", "你好", "...........你好", WidthMonospace},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PadLeftFunc(tt.value, '.', 15, tt.fnWidth)
			tst.DiffError(t, tt.want, got)
		})
	}

}

func TestPadRight(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{

		{"Hello", "Hello.........."},
		{"Привет", "Привет........."},
		{"你好", "你好............."},
		{"###############", "###############"},
		{"################", "################"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := PadRight(tt.value, '.', 15)
			tst.DiffError(t, tt.want, got)
		})
	}

}

func TestPadRightFunc(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    string
		fnWidth WidthFunc
	}{
		{"WidthBytes", "Hello", "Hello..........", WidthBytes},
		{"WidthBytes", "Привет", "Привет...", WidthBytes},
		{"WidthBytes", "你好", "你好.........", WidthBytes},
		{"WidthRunes", "Hello", "Hello..........", WidthRunes},
		{"WidthRunes", "Привет", "Привет.........", WidthRunes},
		{"WidthRunes", "你好", "你好.............", WidthRunes},
		{"WidthMonospace", "Hello", "Hello..........", WidthMonospace},
		{"WidthMonospace", "Привет", "Привет.........", WidthMonospace},
		{"WidthMonospace", "你好", "你好...........", WidthMonospace},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PadRightFunc(tt.value, '.', 15, tt.fnWidth)
			tst.DiffError(t, tt.want, got)
		})
	}

}

func TestChopLineFunc(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		want     string
		maxWidth int
		fnWidth  WidthFunc
	}{
		{"maxWidth==0", "#####", "", 0, WidthBytes},
		{"maxWidth<0", "#####", "", -1, WidthBytes},
		{"WidthBytes", "Hello", "H", 1, WidthBytes},
		{"WidthBytes", "Привет", "", 1, WidthBytes},
		{"WidthBytes", "你好", "", 1, WidthBytes},
		{"WidthBytes", "Hello", "He", 2, WidthBytes},
		{"WidthBytes", "Привет", "П", 2, WidthBytes},
		{"WidthBytes", "你好", "", 2, WidthBytes},
		{"WidthBytes", "Hello", "Hel", 3, WidthBytes},
		{"WidthBytes", "Привет", "П", 3, WidthBytes},
		{"WidthBytes", "你好", "你", 3, WidthBytes},
		{"WidthRunes", "Hello", "H", 1, WidthRunes},
		{"WidthRunes", "Привет", "П", 1, WidthRunes},
		{"WidthRunes", "你好", "你", 1, WidthRunes},
		{"WidthRunes", "Hello", "He", 2, WidthRunes},
		{"WidthRunes", "Привет", "Пр", 2, WidthRunes},
		{"WidthRunes", "你好", "你好", 2, WidthRunes},
		{"WidthRunes", "Hello", "Hel", 3, WidthRunes},
		{"WidthRunes", "Привет", "При", 3, WidthRunes},
		{"WidthRunes", "你好", "你好", 3, WidthRunes},
		{"WidthMonospace", "Hello", "H", 1, WidthMonospace},
		{"WidthMonospace", "Привет", "П", 1, WidthMonospace},
		{"WidthMonospace", "你好", "", 1, WidthMonospace},
		{"WidthMonospace", "Hello", "He", 2, WidthMonospace},
		{"WidthMonospace", "Привет", "Пр", 2, WidthMonospace},
		{"WidthMonospace", "你好", "你", 2, WidthMonospace},
		{"WidthMonospace", "Hello", "Hel", 3, WidthMonospace},
		{"WidthMonospace", "Привет", "При", 3, WidthMonospace},
		{"WidthMonospace", "你好", "你", 3, WidthMonospace},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChopLineFunc(tt.value, tt.maxWidth, tt.fnWidth)
			tst.DiffError(t, tt.want, got)
		})
	}

}

func TestChopMarkLineFunc(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		want     string
		maxWidth int
		fnWidth  WidthFunc
	}{
		{"maxWidth==0", "#####", "", 0, WidthBytes},
		{"maxWidth<0", "#####", "", -1, WidthBytes},
		{"WidthBytes", "Hello Hello", "▒", 3, WidthBytes},
		{"WidthBytes", "Привет Привет", "▒", 3, WidthBytes},
		{"WidthBytes", "你好 你好", "▒", 3, WidthBytes},
		{"WidthBytes", "Hello Hello", "H▒", 4, WidthBytes},
		{"WidthBytes", "Привет Привет", "▒", 4, WidthBytes},
		{"WidthBytes", "你好 你好", "▒", 4, WidthBytes},
		{"WidthBytes", "Hello Hello", "He▒", 5, WidthBytes},
		{"WidthBytes", "Привет Привет", "П▒", 5, WidthBytes},
		{"WidthBytes", "你好 你好", "▒", 5, WidthBytes},
		{"WidthBytes", "Hello Hello", "Hel▒", 6, WidthBytes},
		{"WidthBytes", "Привет Привет", "П▒", 6, WidthBytes},
		{"WidthBytes", "你好 你好", "你▒", 6, WidthBytes},
		{"WidthRunes", "Hello", "▒", 1, WidthRunes},
		{"WidthRunes", "Привет", "▒", 1, WidthRunes},
		{"WidthRunes", "你好", "▒", 1, WidthRunes},
		{"WidthRunes", "Hello", "H▒", 2, WidthRunes},
		{"WidthRunes", "Привет", "П▒", 2, WidthRunes},
		{"WidthRunes", "你好", "你好", 2, WidthRunes},
		{"WidthRunes", "Hello", "He▒", 3, WidthRunes},
		{"WidthRunes", "Привет", "Пр▒", 3, WidthRunes},
		{"WidthRunes", "你好", "你好", 3, WidthRunes},
		{"WidthMonospace", "Hello", "▒", 1, WidthMonospace},
		{"WidthMonospace", "Привет", "▒", 1, WidthMonospace},
		{"WidthMonospace", "你好", "▒", 1, WidthMonospace},
		{"WidthMonospace", "Hello", "H▒", 2, WidthMonospace},
		{"WidthMonospace", "Привет", "П▒", 2, WidthMonospace},
		{"WidthMonospace", "你好", "▒", 2, WidthMonospace},
		{"WidthMonospace", "Hello", "He▒", 3, WidthMonospace},
		{"WidthMonospace", "Привет", "Пр▒", 3, WidthMonospace},
		{"WidthMonospace", "你好", "你▒", 3, WidthMonospace},
		{"WidthMonospace", "Hello", "Hel▒", 4, WidthMonospace},
		{"WidthMonospace", "Привет", "При▒", 4, WidthMonospace},
		{"WidthMonospace", "你好", "你好", 4, WidthMonospace},
	}

	t.Run("panic when maxWidth<chopMarkWidt", func(t *testing.T) {
		defer func() {
			if v := recover(); v != nil {

				want := "panic: [cmn: can't fit `chopMark` with width 3 when `maxWidth` is 1]"
				got := fmt.Sprintf("panic: [%v]", v)

				tst.DiffError(t, want, got)
			}
		}()
		ChopMarkLineFunc("#####", 1, "▒", WidthBytes)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := ChopMarkLineFunc(tt.value, tt.maxWidth, "▒", tt.fnWidth)
			tst.DiffError(t, tt.want, got)
		})
	}

}

func TestChopMarkLine(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		want     string
		chopMark string
		maxWidth int
		fnWidth  WidthFunc
	}{
		{"maxWidth==0", "#####", "", "▒", 0, WidthRunes},
		{"maxWidth<0", "#####", "", "▒", -1, WidthRunes},
		{"With chopMark", "#####", "###▒", "▒", 4, WidthRunes},
		{"Without chopMark", "#####", "####", "", 4, WidthRunes},
		{"WidthRunes", "Hello", "▒", "▒", 1, WidthRunes},
		{"WidthRunes", "Привет", "▒", "▒", 1, WidthRunes},
		{"WidthRunes", "你好", "▒", "▒", 1, WidthRunes},
		{"WidthRunes", "Hello", "H▒", "▒", 2, WidthRunes},
		{"WidthRunes", "Привет", "П▒", "▒", 2, WidthRunes},
		{"WidthRunes", "你好", "你好", "▒", 2, WidthRunes},
		{"WidthRunes", "Hello", "He▒", "▒", 3, WidthRunes},
		{"WidthRunes", "Привет", "Пр▒", "▒", 3, WidthRunes},
		{"WidthRunes", "你好", "你好", "▒", 3, WidthRunes},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChopMarkLine(tt.value, tt.maxWidth, tt.chopMark)
			tst.DiffError(t, tt.want, got)
		})
	}

}

func TestChopMarkLineSelector(t *testing.T) {
	t.Run("returns `ChopMarkLine` function when Widther is `RunesWidther`", func(t *testing.T) {
		var fnWant FnChopMarkLine = ChopMarkLine

		fnGot := ChopMarkLineSelector(RunesWidther{})

		want := fmt.Sprintf("%#v", fnWant)
		got := fmt.Sprintf("%#v", fnGot)

		tst.DiffError(t, want, got)

	})
}
