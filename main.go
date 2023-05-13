// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

// ansible-pretty-print is CLI tool that pretty-prints an output of
// 'ansible-playbook --list-tasks' command
package main

import (
	"os"

	"github.com/keewek/ansible-pretty-print/src/app"
)

func main() {
	r := app.Run(app.DefaultConfig())
	os.Exit(r)
}
