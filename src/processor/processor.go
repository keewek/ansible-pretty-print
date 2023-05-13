// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package processor

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/keewek/ansible-pretty-print/src/cmn"
)

type Row struct {
	Indent int          // TODO: ATM `Indent` is a deadweight.
	Data   fmt.Stringer //       Drop it or implement indent detection
}

type Result struct {
	Rows  []*Row
	Stats *Stats
}

func processPlay(line string) (*Play, error) {
	pair := strings.Split(strings.TrimSpace(line), "TAGS:")

	if len(pair) == 2 {
		name := strings.TrimSpace(pair[0])
		tags := strings.TrimSpace(pair[1])

		return &Play{
			Name: name,
			Tags: tags,
		}, nil
	}

	return nil, errors.New("processor.processPlay: unexpected play format")
}

func processTask(line string) (*Task, error) {
	pair := strings.Split(strings.TrimSpace(line), "TAGS:")

	if len(pair) == 2 {
		block := ""
		name := strings.TrimSpace(pair[0])
		tags := strings.TrimSpace(pair[1])

		pair = strings.SplitN(strings.TrimSpace(name), ":", 2)

		if len(pair) == 2 {
			block = strings.TrimSpace(pair[0])
			name = strings.TrimSpace(pair[1])
		}

		return &Task{
			Block: block,
			Name:  name,
			Tags:  tags,
		}, nil
	}

	return nil, errors.New("processor.processTask: unexpected task format")

}

// func calcIndent(line string) int {
// 	size := len(line)
// 	indent := 0

// 	for i := 0; i < size; i++ {

// 		if line[i] != ' ' {
// 			break
// 		}

// 		indent++
// 	}

// 	return indent
// }

func ProcessLines(scanner *bufio.Scanner, widther cmn.Widther) (*Result, error) {
	var tasks *Tasks

	rows := make([]*Row, 0, 2)
	stats := &Stats{Widther: widther}

	playsCount := 0
	isProcessTasks := false

	for scanner.Scan() {
		line := scanner.Text()
		// indent := calcIndent(line)

		if strings.HasPrefix(line, "  play") {
			playsCount++

			play, err := processPlay(line)
			if err != nil {
				return nil, err
			}

			rows = append(rows, &Row{Indent: 2, Data: play})

			stats.updateWithPlay(play)
			continue

		} else if strings.HasPrefix(line, "    tasks") {
			isProcessTasks = true
			tasks = &Tasks{PlayNumber: playsCount}

			rows = append(rows, &Row{Data: Passthru(line)})
			rows = append(rows, &Row{6, tasks})
			continue

		} else if isProcessTasks {

			if strings.HasPrefix(line, "      ") {
				task, err := processTask(line)
				if err != nil {
					return nil, err
				}

				tasks.Add(task)

				stats.updateWithTask(task)
				continue

			} else {
				isProcessTasks = false
			}
		}

		rows = append(rows, &Row{Data: Passthru(line)})

	}

	result := &Result{rows, stats}

	return result, scanner.Err()
}
