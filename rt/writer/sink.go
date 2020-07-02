package writer

// Sink implements a container for Output.
// it just happens to help provide simple output handling for rt.Runtime implementations.
type Sink struct {
	Output Output
}

func (k *Sink) Writer() Output {
	return k.Output
}

func (k *Sink) SetWriter(out Output) (ret Output) {
	ret, k.Output = k.Output, out
	return
}
