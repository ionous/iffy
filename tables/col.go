package tables

// Col describes a table column
type Col struct {
	Name, Type, Check string
}

// NamesOf returns
func NamesOf(cols []Col) []string {
	keys := make([]string, 0, len(cols))
	for _, c := range cols {
		if len(c.Name) > 0 {
			keys = append(keys, c.Name)
		}
	}
	return keys
}
