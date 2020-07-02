package writer

import "os"

// Sink implements a container for Output.
// it just happens to help provide simple output handling for rt.Runtime implementations.
type Sink struct {
	Output Output
}

func (k *Sink) Writer() Output {
	if k.Output == nil {
		k.Output = &FileWriter{os.Stdout}
	}
	return k.Output
}

func (k *Sink) SetWriter(out Output) (ret Output) {
	ret, k.Output = k.Writer(), out
	return
}
