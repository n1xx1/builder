package builder

import "fmt"

// Column select column
type Column interface {
	IsColumnType()
	WriteTo(Writer) error
}

type columnString string

func (v columnString) WriteTo(w Writer) error {
	_, err := fmt.Fprintf(w, "%s", v)
	if err != nil {
		return err
	}
	return nil
}

func (v columnString) IsColumnType() {}

type columnMssqlTop struct {
	top int
	sub Column
}

func (v *columnMssqlTop) WriteTo(w Writer) error {
	if _, err := fmt.Fprintf(w, "TOP %d ", v.top); err != nil {
		return err
	}
	if err := v.sub.WriteTo(w); err != nil {
		return err
	}
	return nil
}

func (v *columnMssqlTop) IsColumnType() {}
