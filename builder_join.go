// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"strings"
)

// InnerJoin sets inner join
func (b *Builder) InnerJoin(joinTable, joinCond any) *Builder {
	return b.join("INNER JOIN", joinTable, joinCond)
}

// LeftJoin sets left join SQL
func (b *Builder) LeftJoin(joinTable, joinCond any) *Builder {
	return b.join("LEFT JOIN", joinTable, joinCond)
}

// RightJoin sets right join SQL
func (b *Builder) RightJoin(joinTable, joinCond any) *Builder {
	return b.join("RIGHT JOIN", joinTable, joinCond)
}

// CrossJoin sets cross join SQL
func (b *Builder) CrossJoin(joinTable, joinCond any) *Builder {
	return b.join("CROSS JOIN", joinTable, joinCond)
}

// FullJoin sets full join SQL
func (b *Builder) FullJoin(joinTable, joinCond any) *Builder {
	return b.join("FULL JOIN", joinTable, joinCond)
}

// CrossApply sets cross apply SQL
func (b *Builder) CrossApply(joinTable any) *Builder {
	if b.dialect != MSSQL {
		return b.join("CROSS LEFT LATERAL", joinTable, nil)
	}
	return b.join("CROSS APPLY", joinTable, nil)
}

// OuterApply sets outer apply SQL
func (b *Builder) OuterApply(joinTable any) *Builder {
	if b.dialect != MSSQL {
		return b.join("LEFT JOIN LATERAL", joinTable, nil)
	}
	return b.join("OUTER APPLY", joinTable, nil)
}

func (b *Builder) join(joinType string, joinTable, joinCond any) *Builder {
	switch c := joinCond.(type) {
	case nil:
		if strings.HasSuffix(joinType, " APPLY") {
			b.joins = append(b.joins, join{joinType, joinTable, nil})
		}
	case Cond:
		b.joins = append(b.joins, join{joinType, joinTable, c})
	case string:
		b.joins = append(b.joins, join{joinType, joinTable, Expr(c)})
	}

	return b
}

type joins []join

func (joins joins) WriteTo(w Writer) error {
	for _, v := range joins {
		var joinTable = v.joinTable
		var alias string
		if aliased, ok := v.joinTable.(*Aliased); ok {
			joinTable = aliased.val
			alias = aliased.alias + " "
		}

		joinCond := v.joinCond
		if v.joinType == "LEFT JOIN LATERAL" && joinCond == nil {
			joinCond = Expr("TRUE")
		}

		switch tbl := joinTable.(type) {
		case Cond:
			if tbl.IsValid() {
				if _, err := fmt.Fprintf(w, " %s ", v.joinType); err != nil {
					return err
				}
				if err := tbl.WriteTo(w); err != nil {
					return err
				}
				if _, err := fmt.Fprintf(w, " %s", alias); err != nil {
					return err
				}
			}
		case *Builder:
			if _, err := fmt.Fprintf(w, " %s (", v.joinType); err != nil {
				return err
			}
			if err := tbl.WriteTo(w); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, ") %s", alias); err != nil {
				return err
			}
		case string:
			if _, err := fmt.Fprintf(w, " %s %s %s", v.joinType, tbl, alias); err != nil {
				return err
			}
		}

		if joinCond != nil {
			if _, err := fmt.Fprintf(w, "ON "); err != nil {
				return err
			}

			if err := joinCond.WriteTo(w); err != nil {
				return err
			}
		}
	}

	return nil
}
