// Copyright 2016 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
)

type callCond struct {
	fn   string
	args []any
}

// Call generates a function call SQL, requires builder.Cond(s), *builder.Builder or a value as parameters
func Call(fn string, args ...any) Cond {
	return &callCond{fn, args}
}

// WriteTo writes SQL to Writer
func (c *callCond) WriteTo(w Writer) error {
	if _, err := fmt.Fprintf(w, "%s(", c.fn); err != nil {
		return err
	}
	for i, arg := range c.args {
		if i > 0 {
			if _, err := fmt.Fprint(w, ","); err != nil {
				return err
			}
		}
		if err := writeTerm(arg, w); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(w, ")"); err != nil {
		return err
	}
	return nil
}

// And implements And with other conditions
func (c *callCond) And(conds ...Cond) Cond {
	return And(c, And(conds...))
}

// Or implements Or with other conditions
func (c *callCond) Or(conds ...Cond) Cond {
	return Or(c, Or(conds...))
}

// IsValid tests if this condition is valid
func (c *callCond) IsValid() bool {
	return true
}
