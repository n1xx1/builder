// Copyright 2016 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
)

type condExists struct {
	b *Builder
}

// Exists defines equals exists
func Exists(b *Builder) Cond {
	return &condExists{b}
}

// WriteTo writes SQL to Writer
func (e *condExists) WriteTo(w Writer) error {
	if e.b.optype != selectType {
		return ErrUnexpectedExistsQuery
	}

	if _, err := fmt.Fprint(w, "EXISTS ("); err != nil {
		return err
	}
	if err := e.b.WriteTo(w); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, ")"); err != nil {
		return err
	}
	return nil
}

// And implements And with other conditions
func (e *condExists) And(conds ...Cond) Cond {
	return And(e, And(conds...))
}

// Or implements Or with other conditions
func (e *condExists) Or(conds ...Cond) Cond {
	return Or(e, Or(conds...))
}

// IsValid tests if this Eq is valid
func (e *condExists) IsValid() bool {
	return e.b != nil
}
