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
		// erase limit condition
		b.limitation = nil
		ow := w.(*BytesWriter)

		switch strings.ToLower(strings.TrimSpace(b.dialect)) {
		case ORACLE:
			if len(b.selects) == 0 {
				b.selects = append(b.selects, columnString("*"))
			}

			var final *Builder
			selects := b.selects
			b.selects = append(selects, columnString("ROWNUM RN"))

			var wb *Builder
			if b.optype == unionType {
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

			return final.WriteTo(ow)
		case SQLITE, MYSQL, POSTGRES:
			// if type UNION, we need to write previous content back to current writer
			if b.optype == unionType {
				if err := b.WriteTo(ow); err != nil {
					return err
				}
			}

			if limit.offset == 0 {
				fmt.Fprint(ow, " LIMIT ", limit.limitN)
			} else {
				fmt.Fprintf(ow, " LIMIT %v OFFSET %v", limit.limitN, limit.offset)
			}
		case MSSQL:
			fmt.Fprintf(ow, " OFFSET %v ROWS FETCH NEXT %v ROWS ONLY", limit.offset, limit.limitN)
		default:
			return ErrNotSupportType
		}
	}

	return nil
}
