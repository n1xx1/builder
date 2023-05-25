package builder

import (
	"fmt"
	"strings"
)

// CrossApply sets cross apply SQL
func (b *Builder) CrossApply(joinTable any) *Builder {
	return b.apply(applyCross, joinTable)
}

// OuterApply sets cross apply SQL
func (b *Builder) OuterApply(joinTable any) *Builder {
	return b.apply(applyOuter, joinTable)
}

// Apply sets apply table and conditions
func (b *Builder) apply(applyType applyType, joinTable any) *Builder {
	b.applies = append(b.applies, apply{applyType, joinTable})
	return b
}

type applies []apply

func (b *Builder) applyWriteTo(w Writer) error {
	switch strings.ToLower(strings.TrimSpace(b.dialect)) {
	case MSSQL:
		if len(b.applies) > 1 {
			return ErrNotSupportDialectType
		}
		if len(b.applies) == 0 {
			return nil
		}

		applyTable := b.applies[0].applyTable

		var alias string
		if aliased, ok := applyTable.(*Aliased); ok {
			applyTable = aliased.val
			alias = aliased.alias + " "
		}

		applyType := "CROSS"
		if b.applies[0].applyType == applyOuter {
			applyType = "OUTER"
		}

		switch tbl := applyTable.(type) {
		case *Builder:
			if _, err := fmt.Fprintf(w, " %s APPLY (", applyType); err != nil {
				return err
			}
			if err := tbl.WriteTo(w); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, ") %s", alias); err != nil {
				return err
			}
		case string:
			if _, err := fmt.Fprintf(w, " %s APPLY %s %s", applyType, tbl, alias); err != nil {
				return err
			}
		}

	case MYSQL, POSTGRES:
		for _, a := range b.applies {
			applyTable := a.applyTable

			var alias string
			if aliased, ok := applyTable.(*Aliased); ok {
				applyTable = aliased.val
				alias = aliased.alias + " "
			}

			applyType := "INNER"
			if a.applyType == applyOuter {
				applyType = "LEFT"
			}

			switch tbl := applyTable.(type) {
			case *Builder:
				if _, err := fmt.Fprintf(w, " %s JOIN LATERAL (", applyType); err != nil {
					return err
				}
				if err := tbl.WriteTo(w); err != nil {
					return err
				}
				if _, err := fmt.Fprintf(w, ") %s", alias); err != nil {
					return err
				}
			case string:
				if _, err := fmt.Fprintf(w, " %s JOIN %s %s", applyType, tbl, alias); err != nil {
					return err
				}
			}
		}

	default:
		return ErrNotSupportDialectType
	}

	return nil
}
