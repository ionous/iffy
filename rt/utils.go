package rt

// ExecuteList runs a block of statements.
type ExecuteList []Execute

func (x ExecuteList) Execute(run Runtime) (err error) {
	for _, s := range x {
		if e := s.Execute(run); e != nil {
			err = e
			break
		}
	}
	return
}

func (x ExecuteList) ReverseExecute(run Runtime) (err error) {
	for i, cnt := 0, len(x); i < cnt; i++ {
		exec := x[cnt-i-1]
		if e := exec.Execute(run); e != nil {
			err = e
			break
		}
	}
	return
}
