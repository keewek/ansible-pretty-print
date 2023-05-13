// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package processor

import (
	"bufio"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/keewek/ansible-pretty-print/src/cmn"
)

func Test_processPlay(t *testing.T) {
	t.Run("Returns 'Play' struct", func(t *testing.T) {

		tests := []struct {
			input string
			want  *Play
		}{{
			input: "play #1 (vps1): Test	TAGS: []",
			want:  &Play{"play #1 (vps1): Test", "[]"},
		}, {
			input: "play #1 (vps1): Test    TAGS:",
			want:  &Play{"play #1 (vps1): Test", ""},
		}, {
			input: "     play #1 (vps1): TestTAGS: [tag1,   tag2]     ",
			want:  &Play{"play #1 (vps1): Test", "[tag1,   tag2]"},
		}, {
			input: "TAGS:",
			want:  &Play{},
		}}

		for _, tt := range tests {
			t.Run("", func(t *testing.T) {
				got, err := processPlay(tt.input)
				if err != nil {
					t.Fatal(err)
				}
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("(-want +got): \n%s", diff)
				}
			})
		}
	})

	t.Run("Returns error upon unexpected format", func(t *testing.T) {

		tests := []struct {
			input string
			want  *Play
		}{{
			input: "     play #1 (vps1): Test TAG: [tag1,   tag2]     ",
			want:  nil,
		}}

		for _, tt := range tests {
			t.Run("", func(t *testing.T) {
				got, err := processPlay(tt.input)

				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("(-want +got): \n%s", diff)
				}

				if err == nil {
					t.Errorf("expected an error")
				}
			})
		}
	})

}
func Test_processTask(t *testing.T) {
	t.Run("Returns 'Task' struct", func(t *testing.T) {

		tests := []struct {
			input string
			want  *Task
		}{{
			input: "Block: Name	TAGS: [♪, ♪♪, ♪♪♪]",
			want:  &Task{"Block", "Name", "[♪, ♪♪, ♪♪♪]"},
		}, {
			input: "Block NameTAGS: [♪,♪♪,♪♪♪]",
			want:  &Task{"", "Block Name", "[♪,♪♪,♪♪♪]"},
		}, {
			input: "Block NameTAGS:",
			want:  &Task{"", "Block Name", ""},
		}, {
			input: "TAGS:",
			want:  &Task{},
		}}

		for _, tt := range tests {
			t.Run("", func(t *testing.T) {
				got, err := processTask(tt.input)
				if err != nil {
					t.Fatal(err)
				}
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("(-want +got): \n%s", diff)
				}
			})
		}
	})

	t.Run("Returns error upon unexpected format", func(t *testing.T) {

		tests := []struct {
			input string
			want  *Task
		}{{
			input: "     play #1 (vps1): Test TAG: [tag1,   tag2]     ",
			want:  nil,
		}}

		for _, tt := range tests {
			t.Run("", func(t *testing.T) {
				got, err := processTask(tt.input)

				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("(-want +got): \n%s", diff)
				}

				if err == nil {
					t.Errorf("expected an error")
				}
			})
		}
	})

}

func TestProcessLines(t *testing.T) {

	t.Run("Returns 'Result' struct", func(t *testing.T) {
		var ll cmn.LineBuilder

		ll.WriteLine("playbook: playbooks/vsp/playbook_vps.yml")
		ll.WriteLine("")
		ll.WriteLine("  play #1 (vps): Test	TAGS: []")
		ll.WriteLine("    tasks:")
		ll.WriteLine("      Block: Name	TAGS: [Tag1, Tag2]")
		ll.WriteLine("      Gather the package facts	TAGS: [apt, facts, vars]")
		ll.WriteLine("")
		ll.WriteLine("  play #2 (vps): Demo 2	TAGS: []")
		ll.WriteLine("    tasks:")
		ll.WriteLine("      Task 2.1	TAGS: []")
		ll.WriteLine("      Task 2.2	TAGS: []")

		want := &Result{
			[]*Row{
				{0, Passthru("playbook: playbooks/vsp/playbook_vps.yml")},
				{0, Passthru("")},
				{2, &Play{"play #1 (vps): Test", "[]"}},
				{0, Passthru("    tasks:")},
				{6, &Tasks{1, []*Task{
					{"Block", "Name", "[Tag1, Tag2]"},
					{"", "Gather the package facts", "[apt, facts, vars]"},
				}}},
				{0, Passthru("")},
				{2, &Play{"play #2 (vps): Demo 2", "[]"}},
				{0, Passthru("    tasks:")},
				{6, &Tasks{2, []*Task{
					{"", "Task 2.1", "[]"},
					{"", "Task 2.2", "[]"},
				}}},
			},
			&Stats{
				Widther:                      cmn.RunesWidther{},
				LongestPlayDescription:       "play #2 (vps): Demo 2",
				LongestPlayDescriptionLength: 21,
				LongestPlayTags:              "[]",
				LongestPlayTagsLength:        2,
				LongestTaskBlock:             "Block",
				LongestTaskBlockLength:       5,
				LongestTaskName:              "Gather the package facts",
				LongestTaskNameLength:        24,
				LongestTaskDescription:       "Gather the package facts",
				LongestTaskDescriptionLength: 24,
				LongestTaskTags:              "[apt, facts, vars]",
				LongestTaskTagsLength:        18,
			},
		}

		scanner := bufio.NewScanner(strings.NewReader(ll.String()))
		got, err := ProcessLines(scanner, cmn.RunesWidther{})

		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}
	})

	t.Run("Returns error upon unexpected play format", func(t *testing.T) {
		var ll cmn.LineBuilder
		var want *Result

		ll.WriteLine("playbook: playbooks/vsp/playbook_vps.yml")
		ll.WriteLine("")
		ll.WriteLine("  play #1 (vps): Test	TAG: []")
		ll.WriteLine("    tasks:")
		ll.WriteLine("      Block: Name	TAGS: [Tag1, Tag2]")
		ll.WriteLine("      Gather the package facts	TAGS: [apt, facts, vars]")

		scanner := bufio.NewScanner(strings.NewReader(ll.String()))
		got, err := ProcessLines(scanner, cmn.RunesWidther{})

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

		if err == nil {
			t.Errorf("expected an error")
		}

	})

	t.Run("Returns error upon unexpected task format", func(t *testing.T) {
		var ll cmn.LineBuilder
		var want *Result

		ll.WriteLine("playbook: playbooks/vsp/playbook_vps.yml")
		ll.WriteLine("")
		ll.WriteLine("  play #1 (vps): Test	TAGS: []")
		ll.WriteLine("    tasks:")
		ll.WriteLine("      Block: Name	TAG: [Tag1, Tag2]")
		ll.WriteLine("      Gather the package facts	TAGS: [apt, facts, vars]")

		scanner := bufio.NewScanner(strings.NewReader(ll.String()))
		got, err := ProcessLines(scanner, cmn.RunesWidther{})

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got): \n%s", diff)
		}

		if err == nil {
			t.Errorf("expected an error")
		}

	})

}
