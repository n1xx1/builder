package builder

import "fmt"

// Column select column
type Column interface {
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

type columnAs struct {
	val   string
	alias string
}

// As create alias for select
func As(val string, alias string) Column {
	return columnAs{val, alias}
}

func (v columnAs) WriteTo(w Writer) error {
	_, err := fmt.Fprintf(w, "%s as %s", v.val, v.alias)
	if err != nil {
		return err
	}
	return nil
}
