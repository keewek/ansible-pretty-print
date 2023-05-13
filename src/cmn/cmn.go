// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package cmn

import (
	"fmt"
	"os"
	"runtime/debug"
)

func WillUseVar(v ...any) {
	// Do nothing, just silence UnusedVar error
	fmt.Fprintln(os.Stderr, "!!!")
	fmt.Fprintln(os.Stderr, "!!! Silencing UnusedVar errors !!!")
	fmt.Fprintln(os.Stderr, "!!!")
}

func VcsRevision() (rev string, ok bool) {

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value, true
			}
		}
	}

	return rev, ok
}
