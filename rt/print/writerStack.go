package print

import (
	"io"
	"os"
)

type Stack struct {
	stack []io.Writer
}

func (k *Stack) Write(p []byte) (n int, err error) {
	return k.Writer().Write(p)
}

func (k *Stack) PushWriter(w io.Writer) {
	k.stack = append(k.stack, w)
}

func (k *Stack) PopWriter() {
	if cnt := len(k.stack); cnt == 0 {
		panic("Stack: popping an empty stack")
	} else {
		k.stack = k.stack[0 : cnt-1]
	}
}

func (k *Stack) Writer() (ret io.Writer) {
	if cnt := len(k.stack); cnt > 0 {
		ret = k.stack[cnt-1]
	} else {
		ret = os.Stdout
	}
	return
}
