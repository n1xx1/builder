package builder

import (
	"errors"
	"fmt"
	"io"
)

type Aliased struct {
	table interface{}
	alias string
}

func As(table interface{}, alias string) *Aliased {
	return &Aliased{table, alias}
}

func (v *Aliased) IsColumnType() {}

func (v *Aliased) WriteTo(w Writer) error {
	switch v1 := v.table.(type) {
	case *Builder:
		if _, err := io.WriteString(w, "("); err != nil {
			return err
		}
		if err := v1.WriteTo(w); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, ") as %s", v.alias); err != nil {
			return err
		}
	case string:
		_, err := fmt.Fprintf(w, "%s as %s", v1, v.alias)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown table type")
	}

	return nil
}
