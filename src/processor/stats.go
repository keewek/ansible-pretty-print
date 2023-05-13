// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package processor

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"unicode/utf8"

	"github.com/keewek/ansible-pretty-print/src/cmn"
)

type Stats struct {
	Widther                      cmn.Widther
	LongestPlayDescription       string
	LongestPlayDescriptionLength int
	LongestPlayTags              string
	LongestPlayTagsLength        int
	LongestTaskBlock             string
	LongestTaskBlockLength       int
	LongestTaskName              string
	LongestTaskNameLength        int
	LongestTaskDescription       string
	LongestTaskDescriptionLength int
	LongestTaskTags              string
	LongestTaskTagsLength        int
}

func (st *Stats) updatePlayDescription(value string) {
	descriptionLength := st.Widther.Width(value)

	if st.LongestPlayDescriptionLength < descriptionLength {
		st.LongestPlayDescriptionLength = descriptionLength
		st.LongestPlayDescription = value
	}
}

func (st *Stats) updatePlayTags(value string) {
	tagsLength := st.Widther.Width(value)

	if st.LongestPlayTagsLength < tagsLength {
		st.LongestPlayTagsLength = tagsLength
		st.LongestPlayTags = value
	}
}

func (st *Stats) updateTaskBlock(value string) {
	blockLength := st.Widther.Width(value)

	if st.LongestTaskBlockLength < blockLength {
		st.LongestTaskBlockLength = blockLength
		st.LongestTaskBlock = value
	}
}

func (st *Stats) updateTaskName(value string) {
	nameLength := st.Widther.Width(value)

	if st.LongestTaskNameLength < nameLength {
		st.LongestTaskNameLength = nameLength
		st.LongestTaskName = value
	}
}

func (st *Stats) updateTaskDescription(value string) {
	descriptionLength := st.Widther.Width(value)

	if st.LongestTaskDescriptionLength < descriptionLength {
		st.LongestTaskDescriptionLength = descriptionLength
		st.LongestTaskDescription = value
	}
}

func (st *Stats) updateTaskTags(value string) {
	tagsLength := st.Widther.Width(value)

	if st.LongestTaskTagsLength < tagsLength {
		st.LongestTaskTagsLength = tagsLength
		st.LongestTaskTags = value
	}
}

func (st *Stats) updateWithPlay(pl *Play) {
	st.updatePlayDescription(pl.Description())
	st.updatePlayTags(pl.Tags)
}

func (st *Stats) updateWithTask(t *Task) {
	st.updateTaskBlock(t.Block)
	st.updateTaskName(t.Name)
	st.updateTaskDescription(t.Description())
	st.updateTaskTags(t.Tags)
}

func (st *Stats) Lines() []string {
	type field struct {
		Index int
		Name  string
	}

	ft := reflect.TypeOf(*st)
	numField := ft.NumField()
	sv := reflect.ValueOf(*st)

	fields := make([]field, 0, numField)
	lines := make([]string, 0, numField)
	maxLen := 0

	for i := 0; i < numField; i++ {
		curField := ft.Field(i)
		if curField.Name == "Widther" {
			continue
		}
		fields = append(fields, field{i, curField.Name})

		if fl := utf8.RuneCountInString(curField.Name); fl > maxLen {
			maxLen = fl
		}
	}

	for _, field := range fields {
		lines = append(lines, fmt.Sprintf("%*s: %v", maxLen, field.Name, sv.Field(field.Index)))
	}

	return lines
}

func (st *Stats) PrintTo(w io.Writer) {

	// var fnPrint func(string)

	// if lw, ok := w.(cmn.LineWriter); ok {
	// 	fnPrint = func(s string) {
	// 		lw.WriteLine(s)
	// 	}
	// } else {
	// 	fnPrint = func(s string) {
	// 		fmt.Fprintln(w, s)
	// 	}
	// }

	for _, line := range st.Lines() {
		// fnPrint(line)
		fmt.Fprintln(w, line)
	}
}

func (st *Stats) Print() {
	st.PrintTo(os.Stdout)
}
