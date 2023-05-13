// SPDX-FileCopyrightText: 2023 Alexander Bugrov <abugrov+dev@gmail.com>
//
// SPDX-License-Identifier: MIT

package processor

import "fmt"

type Play struct {
	Name string
	Tags string
}

func (pl *Play) Description() string {
	return pl.Name
}

func (pl *Play) String() string {
	return fmt.Sprintf("%s    TAGS: %s", pl.Name, pl.Tags)
}
