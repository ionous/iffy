package writer

import (
	"unicode/utf8"
)

// Chunk wraps the four possible output types so that they can be handled generically.
type Chunk struct {
	Data interface{}
}

func (c *Chunk) Reset() {
	c.Data = nil
}

func (c *Chunk) IsEmpty() (okay bool) {
	switch b := c.Data.(type) {
	case []byte:
		okay = len(b) == 0
	case string:
		okay = len(b) == 0
	}
	return
}

func (c *Chunk) DecodeRune() (ret rune, cnt int) {
	switch b := c.Data.(type) {
	case byte:
		r := rune(b)
		ret, cnt = r, utf8.RuneLen(r)
	case rune:
		ret, cnt = b, utf8.RuneLen(b)
	case []byte:
		ret, cnt = utf8.DecodeRune(b)
	case string:
		ret, cnt = utf8.DecodeRuneInString(b)
	}
	return
}

func (c *Chunk) DecodeLastRune() (ret rune, cnt int) {
	switch b := c.Data.(type) {
	case byte:
		r := rune(b)
		ret, cnt = r, utf8.RuneLen(r)
	case rune:
		ret, cnt = b, utf8.RuneLen(b)
	case []byte:
		ret, cnt = utf8.DecodeLastRune(b)
	case string:
		ret, cnt = utf8.DecodeLastRuneInString(b)
	}
	return
}

func (c *Chunk) WriteTo(w Output) (ret int, err error) {
	switch b := c.Data.(type) {
	case byte:
		if e := w.WriteByte(b); e != nil {
			err = e
		} else {
			ret = 1
		}
	case []byte:
		ret, err = w.Write(b)
	case rune:
		ret, err = w.WriteRune(b)
	case string:
		ret, err = w.WriteString(b)
	}
	return
}

func (c *Chunk) String() (ret string) {
	switch b := c.Data.(type) {
	case []byte:
		ret = string(b)
	case byte:
		ret = string(b)
	case rune:
		ret = string(b)
	case string:
		ret = b
	}
	return
}
