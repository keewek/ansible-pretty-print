// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package processor

import "fmt"

type Task struct {
	Block string
	Name  string
	Tags  string
}

func (t *Task) Description() string {
	var description string
	if t.Block == "" {
		description = t.Name
	} else {
		description = t.Block + ": " + t.Name
	}

	return description
}

func (t *Task) String() string {
	return fmt.Sprintf("%s    TAGS: %s", t.Description(), t.Tags)
}

// ---

type Tasks struct {
	PlayNumber int
	Tasks      []*Task
}

func (ts *Tasks) String() string {
	return fmt.Sprintf("[Play #%d - Tasks]", ts.PlayNumber)
}

func (ts *Tasks) Add(t *Task) {
	ts.Tasks = append(ts.Tasks, t)
}
