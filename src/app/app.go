// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package app

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/keewek/ansible-pretty-print/src/cmn"
	"github.com/keewek/ansible-pretty-print/src/processor"
	"github.com/keewek/ansible-pretty-print/src/ui"
)

// _go:generate sh -c 'git describe --abbrev=0 > version.txt'
//
//go:embed version.txt
var ver string

// _go:generate sh -c 'git describe > describe.txt'
// _go:embed describe.txt
// var desc string

func version(output io.Writer) func() {
	return func() {

		ver = strings.TrimSpace(ver)
		// desc = strings.TrimSpace(desc)

		if !strings.HasPrefix(ver, "v") {
			ver = "v" + ver
		}

		revision := ""
		if rev, ok := cmn.VcsRevision(); ok {
			revision = fmt.Sprintf(" (revision: %v)", rev)
			// revision = fmt.Sprintf(" (%s; %s)", desc, rev)
		}

		fmt.Fprintf(output, "ansible-pretty-print %s%s\n", ver, revision)
		fmt.Fprintln(output, "Copyright: 2023 Alexander Bugrov")
		fmt.Fprintln(output, "License: MIT <https://spdx.org/licenses/MIT.html>")
		fmt.Fprintln(output)
	}
}

func usage(output io.Writer) func() {
	return func() {
		fmt.Fprintf(output, "Usage: %v [OPTION]... [FILE]\n", os.Args[0])
		fmt.Fprintln(output, "Pretty-print Ansible's --list-tasks output")
		fmt.Fprintln(output)
		flag.PrintDefaults()
	}
}

func Run(c *Config) int {

	ui.Box = c.AcquireBoxChars()
	ui.FnWidth = c.Widther.Width

	if c.IsVersion {
		version(c.OutErr)()
		return 0
	}

	scanner, closer, err := c.AcquireScanner()
	if err != nil {
		fmt.Fprintf(c.OutErr, "app.Run: %v\n", err)
		return 1
	}

	defer closer()

	if scanner == nil {
		usage(c.OutErr)()
		return 0
	}

	result, err := processor.ProcessLines(scanner, c.Widther)
	if err != nil {
		fmt.Fprintf(c.OutErr, "app.Run: %v\n", err)
		return 1
	}

	if c.IsStats {
		ui.MsgBoxTo(c.Out, result.Stats.Lines())
	}

	p := c.AcquirePrinter()
	p.PrintTo(c.Out, result)

	return 0
}
