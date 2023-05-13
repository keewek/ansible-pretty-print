// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package eol

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_EndOfLine(t *testing.T) {
	t.Run("CRLF", func(t *testing.T) {
		want := "\r\n"
		got := CRLF

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("LF", func(t *testing.T) {
		want := "\n"
		got := LF

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("CrLf", func(t *testing.T) {
		eol := CrLf()

		want := "\r\n"
		got := eol.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

		got2 := string(eol)

		if diff := cmp.Diff(want, got2); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("Lf", func(t *testing.T) {
		eol := Lf()

		want := "\n"
		got := eol.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

		got2 := string(eol)

		if diff := cmp.Diff(want, got2); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("Default is LF", func(t *testing.T) {
		var eol EndOfLine

		want := "\n"
		got := eol.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("Zero value cast to string is empty", func(t *testing.T) {
		var eol EndOfLine

		want := ""
		got := string(eol)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})

	t.Run("Custom value", func(t *testing.T) {
		var eol EndOfLine = "#END#"

		want := "#END#"
		got := string(eol)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

		got2 := string(eol)

		if diff := cmp.Diff(want, got2); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

	})
}
