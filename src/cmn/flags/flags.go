// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package flags

import (
	"flag"
)

var activeFlags map[string]struct{} = make(map[string]struct{})

// func init() {
// 	ReScan()
// }

func ReScan() {
	DisableAll()

	flag.Visit(func(f *flag.Flag) {
		activeFlags[f.Name] = struct{}{}
	})
}

func IsSet(name string) bool {
	if len(activeFlags) == 0 {
		ReScan()
	}
	_, ok := activeFlags[name]

	return ok
}

func Enable(name ...string) {

	for _, name := range name {
		activeFlags[name] = struct{}{}
	}
}

func EnableAll() {

	flag.VisitAll(func(f *flag.Flag) {
		activeFlags[f.Name] = struct{}{}
	})
}

func Disable(name ...string) {
	for _, name := range name {
		delete(activeFlags, name)
	}
}

func DisableAll() {
	if len(activeFlags) > 0 {
		activeFlags = make(map[string]struct{})
	}
}
