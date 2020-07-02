package writer

import "os"

type FileWriter struct {
	*os.File
}

func (fp FileWriter) WriteByte(c byte) error {
	_, e := fp.File.Write([]byte{c})
	return e
}

func (fp FileWriter) WriteRune(r rune) (int, error) {
	return fp.File.WriteString(string(r))
}
