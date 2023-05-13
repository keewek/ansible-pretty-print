// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package processor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Play(t *testing.T) {
	p := &Play{Name: "Name", Tags: "Tags"}

	t.Run("Implements 'Stringer' interface", func(t *testing.T) {
		got := p.String()
		want := "Name    TAGS: Tags"

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("Description(): returns 'Name' field", func(t *testing.T) {
		got := p.Description()
		want := "Name"

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

}
