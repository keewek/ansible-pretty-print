// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package tst

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// ReprValue returns a Go-syntax representation of the value
//
// Equivalent to fmt.Sprintf("%#v", value)
func ReprValue(value any) string {
	return fmt.Sprintf("%#v", value)
}

// ReprType returns a Go-syntax representation of the type of the value
//
// Equivalent to fmt.Sprintf("%T, value)
func ReprType(value any) string {
	return fmt.Sprintf("%T", value)
}

func DiffError(t *testing.T, want, got any) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got): \n%s", diff)
	}
}

func DiffErrorBefore(t *testing.T, want, got any, fn func()) {
	if diff := cmp.Diff(want, got); diff != "" {
		fn()
		t.Errorf("(-want +got): \n%s", diff)
	}
}
