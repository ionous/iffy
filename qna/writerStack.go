package qna

import (
	"io"
	"os"
)

type WriterStack struct {
	stack []io.Writer
}

func (k *WriterStack) PushWriter(w io.Writer) {
	k.stack = append(k.stack, w)
}

func (k *WriterStack) PopWriter() {
	if cnt := len(k.stack); cnt == 0 {
		panic("WriterStack: popping an empty stack")
	} else {
		k.stack = k.stack[0 : cnt-1]
	}
}

func (k *WriterStack) Writer() (ret io.Writer) {
	if cnt := len(k.stack); cnt > 0 {
		ret = k.stack[cnt-1]
	} else {
		ret = os.Stdout
	}
	return
}
