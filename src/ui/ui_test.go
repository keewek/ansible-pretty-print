// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package ui

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/keewek/ansible-pretty-print/src/cmn"
)

func TestMsgBoxFunc(t *testing.T) {
	data := []string{
		"1",
		"12",
		"Hello world",
		"Всем привет",
		"你好世界",
	}

	t.Run("BoxCharsAscii", func(t *testing.T) {
		asciiBytes, err := os.ReadFile("testdata/ascii_bytes.txt")
		if err != nil {
			t.Fatal(err)
		}

		asciiRunes, err := os.ReadFile("testdata/ascii_runes.txt")
		if err != nil {
			t.Fatal(err)
		}

		asciiMonospace, err := os.ReadFile("testdata/ascii_monospace.txt")
		if err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			name    string
			want    string
			fnWidth cmn.WidthFunc
		}{
			{"WidthBytes", string(asciiBytes), cmn.WidthBytes},
			{"WidthRunes", string(asciiRunes), cmn.WidthRunes},
			{"WidthMonospace", string(asciiMonospace), cmn.WidthMonospace},
		}

		var b cmn.LineBuilder
		box := cmn.BoxCharsAscii()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				b.Reset()
				MsgBoxFunc(&b, data, &box, tt.fnWidth)
				got := b.String()
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Logf("-want: \n%v\n", tt.want)
					t.Logf("-got: \n%v\n", got)
					t.Errorf("(-want +got): \n%s", diff)
				}
			})
		}
	})

	t.Run("BoxCharsDos", func(t *testing.T) {
		dosBytes, err := os.ReadFile("testdata/dos_bytes.txt")
		if err != nil {
			t.Fatal(err)
		}

		dosRunes, err := os.ReadFile("testdata/dos_runes.txt")
		if err != nil {
			t.Fatal(err)
		}

		dosMonospace, err := os.ReadFile("testdata/dos_monospace.txt")
		if err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			name    string
			want    string
			fnWidth cmn.WidthFunc
		}{
			{"WidthBytes", string(dosBytes), cmn.WidthBytes},
			{"WidthRunes", string(dosRunes), cmn.WidthRunes},
			{"WidthMonospace", string(dosMonospace), cmn.WidthMonospace},
		}

		var b cmn.LineBuilder
		box := cmn.BoxCharsDos()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				b.Reset()
				MsgBoxFunc(&b, data, &box, tt.fnWidth)
				got := b.String()
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Logf("-want: \n%v\n", tt.want)
					t.Logf("-got: \n%v\n", got)
					t.Errorf("(-want +got): \n%s", diff)
				}
			})
		}
	})

}
