package print

import (
	"strings"

	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt/writer"
)

// Parens filters writer.Output, parenthesizing a stream of writes. Close adds the closing paren.
func Parens(out writer.Output) *Filter {
	return Brackets(out, '(', ')')
}

func Brackets(out writer.Output, open, close rune) *Filter {
	return &Filter{
		First: func(c Chunk) (int, error) {
			n, _ := out.WriteRune('(')
			x, _ := c.WriteTo(out)
			return n + x, nil
		},
		Rest: func(c Chunk) (int, error) {
			return c.WriteTo(out)
		},
		Last: func() error {
			_, e := out.WriteRune(')')
			return e
		},
	}
}

// Capitalize filters writer.Output, capitalizing the first string.
func Capitalize(out writer.Output) *Filter {
	return &Filter{
		First: func(c Chunk) (int, error) {
			cap := lang.Capitalize(c.String())
			return out.WriteString(cap)
		},
		Rest: func(c Chunk) (int, error) {
			return c.WriteTo(out)
		},
	}
}

// TitleCase filters writer.Output, capitalizing every write.
func TitleCase(out writer.Output) *Filter {
	return &Filter{
		Rest: func(c Chunk) (int, error) {
			cap := lang.Capitalize(c.String())
			return out.WriteString(cap)
		},
	}
}

// Lowercase filters writer.Output, lowering every string.
func Lowercase(out writer.Output) *Filter {
	return &Filter{
		Rest: func(c Chunk) (int, error) {
			cap := strings.ToLower(c.String())
			return out.WriteString(cap)
		},
	}
}

// Slash filters writer.Output, separating writes with a slash.
func Slash(out writer.Output) *Filter {
	return &Filter{
		First: func(c Chunk) (ret int, err error) {
			return c.WriteTo(out)
		},
		Rest: func(c Chunk) (ret int, err error) {
			x, _ := out.WriteString(" /")
			n, _ := c.WriteTo(out)
			return n + x, nil
		},
	}
}
