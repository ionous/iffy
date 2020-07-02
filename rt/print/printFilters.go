package print

import (
	"strings"

	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt/writer"
)

// Parens filters writer.Output, parenthesizing a stream of writes. Close adds the closing paren.
func Parens(out writer.Output) writer.OutputCloser {
	return Brackets(out, '(', ')')
}

func Brackets(out writer.Output, open, close rune) writer.OutputCloser {
	f := &Filter{
		First: func(c writer.Chunk) (int, error) {
			n, _ := out.WriteRune('(')
			x, _ := c.WriteTo(out)
			return n + x, nil
		},
		Rest: func(c writer.Chunk) (int, error) {
			return c.WriteTo(out)
		},
		Last: func() error {
			_, e := out.WriteRune(')')
			return e
		},
	}
	writer.InitChunks(f)
	return f
}

// Capitalize filters writer.Output, capitalizing the first string.
func Capitalize(out writer.Output) writer.Output {
	f := &Filter{
		First: func(c writer.Chunk) (int, error) {
			cap := lang.Capitalize(c.String())
			return out.WriteString(cap)
		},
		Rest: func(c writer.Chunk) (int, error) {
			return c.WriteTo(out)
		},
	}
	writer.InitChunks(f)
	return f
}

// TitleCase filters writer.Output, capitalizing every write.
func TitleCase(out writer.Output) writer.Output {
	f := &Filter{
		Rest: func(c writer.Chunk) (int, error) {
			cap := lang.Capitalize(c.String())
			return out.WriteString(cap)
		},
	}
	writer.InitChunks(f)
	return f
}

// Lowercase filters writer.Output, lowering every string.
func Lowercase(out writer.Output) writer.Output {
	f := &Filter{
		Rest: func(c writer.Chunk) (int, error) {
			cap := strings.ToLower(c.String())
			return out.WriteString(cap)
		},
	}
	writer.InitChunks(f)
	return f
}

// Slash filters writer.Output, separating writes with a slash.
func Slash(out writer.Output) writer.Output {
	f := &Filter{
		First: func(c writer.Chunk) (ret int, err error) {
			return c.WriteTo(out)
		},
		Rest: func(c writer.Chunk) (ret int, err error) {
			x, _ := out.WriteString(" /")
			n, _ := c.WriteTo(out)
			return n + x, nil
		},
	}
	writer.InitChunks(f)
	return f
}
