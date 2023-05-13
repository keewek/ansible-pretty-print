// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package app

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/keewek/ansible-pretty-print/src/cmn"
	"github.com/keewek/ansible-pretty-print/src/cmn/flags"
	"github.com/keewek/ansible-pretty-print/src/printer"
	"github.com/keewek/ansible-pretty-print/src/processor"
	"golang.org/x/term"
)

// === start: Flags ===

const (
	kFlagIsChop    = "chop"
	kFlagIsDos     = "dos"
	kFlagIsIndent  = "indent"
	kFlagIsMono    = "mono"
	kFlagIsStats   = "stats"
	kFlagIsStdin   = "stdin"
	kFlagIsTable   = "table"
	kFlagIsVersion = "version"
	kFlagWidth     = "width"
)

var (
	flagIsChop    = flag.Bool(kFlagIsChop, false, "chop long lines")
	flagIsDos     = flag.Bool(kFlagIsDos, false, "DOS box-drawing characters")
	flagIsIndent  = flag.Bool(kFlagIsIndent, false, "indent block/role")
	flagIsMono    = flag.Bool(kFlagIsMono, false, "calculate string width as monospace width")
	flagIsStats   = flag.Bool(kFlagIsStats, false, "print stats")
	flagIsStdin   = flag.Bool(kFlagIsStdin, false, "read standard input")
	flagIsTable   = flag.Bool(kFlagIsTable, false, "table output")
	flagIsVersion = flag.Bool(kFlagIsVersion, false, "output version information")
	flagWidth     = flag.Int(kFlagWidth, 0, "custom line width")
)

// === end: Flags ===

const (
	DefaultTermWidth = 80
)

type Printer interface {
	Print(data *processor.Result)
	PrintTo(output io.Writer, data *processor.Result)
}

type TermSizeFunc func() (cols int, lines int, err error)

type Config struct {
	Filepath  string
	IsChop    bool
	IsDos     bool
	IsIndent  bool
	IsMono    bool
	IsStats   bool
	IsStdin   bool
	IsTable   bool
	IsVersion bool
	TermWidth int
	Widther   cmn.Widther
	Out       io.Writer
	OutErr    io.Writer
}

// func isTerminal() bool {
// 	if isatty.IsTerminal(os.Stdin.Fd()) {
// 		return true
// 	} else if isatty.IsCygwinTerminal(os.Stdin.Fd()) {
// 		return true
// 	}

// 	return false
// }

func termSize() (cols int, lines int, err error) {
	cols, lines, err = term.GetSize(int(os.Stdout.Fd()))

	return cols, lines, err
}

func DefaultConfig() *Config {

	c := &Config{
		TermWidth: DefaultTermWidth,
		Widther:   cmn.RunesWidther{},
		Out:       os.Stdout,
		OutErr:    os.Stderr,
	}

	c.ApplyFlags()
	c.Init(termSize)

	return c
}

func (c *Config) Init(fnTermSize TermSizeFunc) {
	if c.IsMono {
		c.Widther = cmn.MonospaceWidther{}
	} else {
		c.Widther = cmn.RunesWidther{}
	}

	if (c.IsChop || c.IsTable) && !flags.IsSet("width") {
		// Try determine terminal width
		cols, _, err := fnTermSize()
		if err != nil {
			fmt.Fprintln(c.OutErr, "")
			fmt.Fprintf(c.OutErr, "!!! Can't determine terminal width!\n")
			fmt.Fprintf(c.OutErr, "!!!    [width: %v; err: %v]\n", cols, err)
			fmt.Fprintf(c.OutErr, "!!!\n")
			fmt.Fprintf(c.OutErr, "!!! Falling back to default width of %v chars.\n", DefaultTermWidth)
			fmt.Fprintf(c.OutErr, "!!! Use --width to specify custom value.\n")
			fmt.Fprintf(c.OutErr, "!!!    e.g., --width $(tput cols)\n")
			fmt.Fprintln(c.OutErr, "")
		} else {
			c.TermWidth = cols
		}
	}
}

func (c *Config) AcquireBoxChars() cmn.BoxChars {
	if c.IsDos {
		return cmn.BoxCharsDos()
	}

	return cmn.BoxCharsAscii()
}

func (c *Config) AcquirePrinter() Printer {
	var p Printer

	if c.IsTable {
		tp := printer.NewTablePrinter()
		tp.SetWidther(c.Widther)
		tp.SetMaxLineWidth(c.TermWidth)
		if c.IsDos {
			tp.SetBoxChars(cmn.BoxCharsDos())
		}

		p = tp
	} else {
		cp := printer.NewColumnPrinter()

		cp.SetWidther(c.Widther)
		cp.SetIsChopLines(c.IsChop)
		cp.SetMaxLineWidth(c.TermWidth)
		cp.SetIsIndentBlock(c.IsIndent)

		p = cp
	}

	return p
}

func (c *Config) AcquireScanner() (scanner *bufio.Scanner, closer func(), _ error) {

	closer = func() {}

	if c.IsStdin {
		scanner = bufio.NewScanner(os.Stdin)
	} else if c.Filepath != "" {
		file, err := os.Open(c.Filepath)

		if err != nil {
			return nil, closer, fmt.Errorf("Config.AcquireScanner: %w", err)
		}

		scanner = bufio.NewScanner(file)
		closer = func() {
			file.Close()

		}
	}

	return scanner, closer, nil
}

func (c *Config) ApplyFlags() {

	if !flag.Parsed() {
		flag.Usage = usage(c.OutErr)
		flag.Parse()
	}

	if fp := flag.Arg(0); fp != "" {
		c.Filepath = fp
	}

	if flags.IsSet(kFlagIsChop) {
		c.IsChop = *flagIsChop
	}

	if flags.IsSet(kFlagIsDos) {
		c.IsDos = *flagIsDos
	}

	if flags.IsSet(kFlagIsIndent) {
		c.IsIndent = *flagIsIndent
	}

	if flags.IsSet(kFlagIsMono) {
		c.IsMono = *flagIsMono
	}

	if flags.IsSet(kFlagIsStats) {
		c.IsStats = *flagIsStats
	}

	if flags.IsSet(kFlagIsStdin) {
		c.IsStdin = *flagIsStdin
	}

	if flags.IsSet(kFlagIsTable) {
		c.IsTable = *flagIsTable
	}

	if flags.IsSet(kFlagIsVersion) {
		c.IsVersion = *flagIsVersion
	}

	if flags.IsSet(kFlagWidth) {
		c.TermWidth = *flagWidth
	}

}
