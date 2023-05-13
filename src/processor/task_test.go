// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package processor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Task(t *testing.T) {
	task := &Task{Block: "Block", Name: "Name", Tags: "Tags"}

	t.Run("Implements 'Stringer' interface", func(t *testing.T) {
		got := task.String()
		want := "Block: Name    TAGS: Tags"

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("Description(): returns 'Name' field when 'Block' field is empty", func(t *testing.T) {
		task := &Task{Block: "", Name: "Name", Tags: "Tags"}

		got := task.Description()
		want := "Name"

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("Description(): returns a join of 'Block' and 'Name' fields when 'Block' field is not empty", func(t *testing.T) {
		task := &Task{Block: "Block", Name: "Name", Tags: "Tags"}

		got := task.Description()
		want := "Block: Name"

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

}

func Test_Tasks(t *testing.T) {

	items := []*Task{
		{Block: "Block_1", Name: "Name_1", Tags: "Tags_1"},
		{Block: "Block_2", Name: "Name_2", Tags: "Tags_2"},
		{Block: "Block_3", Name: "Name_3", Tags: "Tags_3"},
	}

	task04 := &Task{Block: "Block_4", Name: "Name_4", Tags: "Tags_4"}

	tasks := &Tasks{PlayNumber: 42, Tasks: items}

	t.Run("Implements 'Stringer' interface", func(t *testing.T) {
		got := tasks.String()
		want := "[Play #42 - Tasks]"

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Play.String() mismatch (-want +got): \n%s", diff)
		}
	})

	t.Run("Add(): adds task", func(t *testing.T) {
		want := make([]*Task, len(items))

		copy(want, items)
		want = append(want, task04)

		tasks.Add(task04)

		got := tasks.Tasks

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Play.String() mismatch (-want +got): \n%s", diff)
		}
	})

}
