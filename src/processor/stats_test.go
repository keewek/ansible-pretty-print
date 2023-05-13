// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package processor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/keewek/ansible-pretty-print/src/cmn"
)

func Test_Stats(t *testing.T) {
	// s := &Stats{}

	t.Run("updatePlayDescription():", func(t *testing.T) {

		got := Stats{Widther: cmn.RunesWidther{}}
		want := Stats{
			Widther:                      cmn.RunesWidther{},
			LongestPlayDescription:       "♪♪♪",
			LongestPlayDescriptionLength: 3,
		}

		got.updatePlayDescription("♪♪♪")
		got.updatePlayDescription("♪♪")
		got.updatePlayDescription("♪")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("updatePlayTags():", func(t *testing.T) {

		got := Stats{Widther: cmn.RunesWidther{}}
		want := Stats{
			Widther:               cmn.RunesWidther{},
			LongestPlayTags:       "♪♪♪",
			LongestPlayTagsLength: 3,
		}

		got.updatePlayTags("♪♪♪")
		got.updatePlayTags("♪♪")
		got.updatePlayTags("♪")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("updateBlock():", func(t *testing.T) {

		got := Stats{Widther: cmn.RunesWidther{}}
		want := Stats{
			Widther:                cmn.RunesWidther{},
			LongestTaskBlock:       "♪♪♪",
			LongestTaskBlockLength: 3,
		}

		got.updateTaskBlock("♪♪♪")
		got.updateTaskBlock("♪♪")
		got.updateTaskBlock("♪")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("updateName():", func(t *testing.T) {

		got := Stats{Widther: cmn.RunesWidther{}}
		want := Stats{
			Widther:               cmn.RunesWidther{},
			LongestTaskName:       "♪♪♪",
			LongestTaskNameLength: 3,
		}

		got.updateTaskName("♪♪♪")
		got.updateTaskName("♪♪")
		got.updateTaskName("♪")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("updateTaskDescription():", func(t *testing.T) {

		got := Stats{Widther: cmn.RunesWidther{}}
		want := Stats{
			Widther:                      cmn.RunesWidther{},
			LongestTaskDescription:       "♪♪♪",
			LongestTaskDescriptionLength: 3,
		}

		got.updateTaskDescription("♪♪♪")
		got.updateTaskDescription("♪♪")
		got.updateTaskDescription("♪")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("updateTaskTags():", func(t *testing.T) {

		got := Stats{Widther: cmn.RunesWidther{}}
		want := Stats{
			Widther:               cmn.RunesWidther{},
			LongestTaskTags:       "♪♪♪",
			LongestTaskTagsLength: 3,
		}

		got.updateTaskTags("♪♪♪")
		got.updateTaskTags("♪♪")
		got.updateTaskTags("♪")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("updateWithPlay():", func(t *testing.T) {

		got := Stats{Widther: cmn.RunesWidther{}}
		want := Stats{
			Widther:                      cmn.RunesWidther{},
			LongestPlayDescription:       "♪♪♪ Name ♪♪♪",
			LongestPlayDescriptionLength: 12,
			LongestPlayTags:              "♪♪♪ Tags ♪♪♪",
			LongestPlayTagsLength:        12,
		}

		got.updateWithPlay(&Play{"♪♪♪ Name ♪♪♪", "♪♪♪ Tags ♪♪♪"})
		got.updateWithPlay(&Play{"♪♪ Name ♪♪", "♪♪ Tags ♪♪"})
		got.updateWithPlay(&Play{"♪ Name ♪", "♪ Tags ♪"})

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("updateWithTask():", func(t *testing.T) {

		got := Stats{Widther: cmn.RunesWidther{}}
		want := Stats{
			Widther:                      cmn.RunesWidther{},
			LongestTaskBlock:             "♪♪♪ Block ♪♪♪",
			LongestTaskBlockLength:       13,
			LongestTaskName:              "♪♪♪ Name ♪♪♪",
			LongestTaskNameLength:        12,
			LongestTaskDescription:       "♪♪♪ Block ♪♪♪: ♪♪♪ Name ♪♪♪",
			LongestTaskDescriptionLength: 27,
			LongestTaskTags:              "♪♪♪ Tags ♪♪♪",
			LongestTaskTagsLength:        12,
		}

		got.updateWithTask(&Task{"♪♪♪ Block ♪♪♪", "♪♪♪ Name ♪♪♪", "♪♪♪ Tags ♪♪♪"})
		got.updateWithTask(&Task{"♪♪ Block ♪♪", "♪♪ Name ♪♪", "♪♪ Tags ♪♪"})
		got.updateWithTask(&Task{"♪ Block ♪", "♪ Name ♪", "♪ Tags ♪"})

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

}
