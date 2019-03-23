package builder

import "fmt"

// When Case when
type When struct {
	Cond  Cond
	Value string
}

type columnCase struct {
	parts  []When
	defVal string
	alias  *string
}

// Case ...
func Case(defVal string, whens ...When) Column {
	var parts = make([]When, 0, len(whens))
	for _, val := range whens {
		parts = append(parts, val)
	}
	return columnCase{
		parts:  parts,
		defVal: defVal,
	}
}

// CaseAs ...
func CaseAs(alias string, defVal string, whens ...When) Column {
	var parts = make([]When, 0, len(whens))
	for _, val := range whens {
		parts = append(parts, val)
	}
	aliasS := new(string)
	*aliasS = alias
	return columnCase{
		parts:  parts,
		defVal: defVal,
		alias:  aliasS,
	}
}

func (v columnCase) WriteTo(w Writer) error {
	fmt.Fprint(w, "(case ")

	for _, x := range v.parts {
		_, err := fmt.Fprint(w, "when ")
		if err != nil {
			return err
		}
		err = x.Cond.WriteTo(w)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, " then %s ", x.Value)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w, "else %s end)", v.defVal)
	if err != nil {
		return err
	}

	if v.alias != nil {
		_, err := fmt.Fprintf(w, " as %s", *v.alias)
		if err != nil {
			return err
		}
	}

	return nil
}
