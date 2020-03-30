package rt

// Block runs a block of statements.
type Block []Execute

func (x Block) Execute(run Runtime) (err error) {
	for _, s := range x {
		if e := s.Execute(run); e != nil {
			err = e
			break
		}
	}
	return
}

func (x Block) ReverseExecute(run Runtime) (err error) {
	for i, cnt := 0, len(x); i < cnt; i++ {
		exec := x[cnt-i-1]
		if e := exec.Execute(run); e != nil {
			err = e
			break
		}
	}
	return
}
