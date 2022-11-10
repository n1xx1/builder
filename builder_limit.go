// Copyright 2018 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"strings"
)

func (b *Builder) limitWriteTo(w Writer) error {
	if strings.TrimSpace(b.dialect) == "" {
		return ErrDialectNotSetUp
	}

	if b.limitation != nil {
		limit := b.limitation
		if limit.offset < 0 || limit.limitN <= 0 {
			return ErrInvalidLimitation
		}
		// unset limit condition to prevent final.WriteTo from recursing forever
		b.limitation = nil
		defer func() {
			b.limitation = limit
		}()

		switch strings.ToLower(strings.TrimSpace(b.dialect)) {
		case ORACLE:
			if len(b.selects) == 0 {
				b.selects = append(b.selects, columnString("*"))
			}

			var final *Builder
			selects := b.selects
			b.selects = append(selects, columnString("ROWNUM RN"))

			var wb *Builder
			if b.optype == setOpType {
				wb = Dialect(b.dialect).Select("at.*", "ROWNUM RN").
					From(b, "at")
			} else {
				wb = b
			}

			selectsInterface := make([]interface{}, len(selects))
			for i, v := range selects {
				selectsInterface[i] = v
			}

			if limit.offset == 0 {
				final = Dialect(b.dialect).Select(selectsInterface...).From(wb, "at").
					Where(Lte{"at.RN": limit.limitN})
			} else {
				sub := Dialect(b.dialect).Select("*").
					From(b, "at").Where(Lte{"at.RN": limit.offset + limit.limitN})

				final = Dialect(b.dialect).Select(selectsInterface...).From(sub, "att").
					Where(Gt{"att.RN": limit.offset})
			}

			return final.WriteTo(w)
		case SQLITE, MYSQL, POSTGRES:
			// if type UNION, we need to write previous content back to current writer
			if b.optype == setOpType {
				if err := b.WriteTo(w); err != nil {
					return err
				}
			}

			if limit.offset == 0 {
				if _, err := fmt.Fprint(w, " LIMIT ", limit.limitN); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(w, " LIMIT %v OFFSET %v", limit.limitN, limit.offset); err != nil {
					return err
				}
			}
		case MSSQL:
			if false {
				if len(b.selects) == 0 {
					b.selects = append(b.selects, columnString("*"))
				}

				var final *Builder
				selects := b.selects
				b.selects[0] = &columnMssqlTop{limit.limitN + limit.offset, b.selects[0]}
				b.selects = append(b.selects, columnString("ROW_NUMBER() OVER (ORDER BY (SELECT 1)) AS RN"))

				var wb *Builder
				if b.optype == setOpType {
					wb = Dialect(b.dialect).Select("*", "ROW_NUMBER() OVER (ORDER BY (SELECT 1)) AS RN").
						From(b, "at")
				} else {
					wb = b
				}

				if limit.offset == 0 {
					final = Dialect(b.dialect).From(wb, "at")
					final.selects = selects // manually patch becuase []Column is not convertible to []interface{}
				} else {
					final = Dialect(b.dialect).From(wb, "at").Where(Gt{"at.RN": limit.offset})
					final.selects = selects // manually patch becuase []Column is not convertible to []interface{}
				}

				return final.WriteTo(w)
			}
			if _, err := fmt.Fprintf(w, " OFFSET %v ROWS FETCH NEXT %v ROWS ONLY", limit.offset, limit.limitN); err != nil {
				return err
			}
		default:
			return ErrNotSupportType
		}
	}

	return nil
}
