// Copyright 2016 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
)

// Operator is the operator to use inside the Op condition
type Operator int

const (
	// OpEq is the = operator
	OpEq Operator = iota
	// OpNeq is the <> operator
	OpNeq
)

type opCond struct {
	op Operator
	a  interface{}
	b  interface{}
}

func (o Operator) String() string {
	switch o {
	case OpEq:
		return "="
	case OpNeq:
		return "<>"
	default:
		return ""
	}
}

// Op generates a comparison operator SQL
func Op(o Operator, a interface{}, b interface{}) Cond {
	return &opCond{o, a, b}
}

func writeTerm(v interface{}, w Writer) error {
	switch v.(type) {
	case Cond:
		if _, err := fmt.Fprint(w, "("); err != nil {
			return err
		}
		if err := v.(Cond).WriteTo(w); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, ")"); err != nil {
			return err
		}
	case *Builder:
		if _, err := fmt.Fprint(w, "("); err != nil {
			return err
		}
		if err := v.(*Builder).WriteTo(w); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, ")"); err != nil {
			return err
		}
	default:
		if _, err := fmt.Fprint(w, "?"); err != nil {
			return err
		}
		w.Append(v)
	}
	return nil
}

// WriteTo writes SQL to Writer
func (o *opCond) WriteTo(w Writer) error {
	if err := writeTerm(o.a, w); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, o.op.String()); err != nil {
		return err
	}
	if err := writeTerm(o.b, w); err != nil {
		return err
	}
	return nil
}

// And implements And with other conditions
func (o *opCond) And(conds ...Cond) Cond {
	return And(o, And(conds...))
}

// Or implements Or with other conditions
func (o *opCond) Or(conds ...Cond) Cond {
	return Or(o, Or(conds...))
}

// IsValid tests if this op is valid
func (o *opCond) IsValid() bool {
	return true
}
