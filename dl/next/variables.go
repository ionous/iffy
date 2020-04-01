package next

import "github.com/ionous/iffy/rt"

// GetVariable writes the value at 'name'
type GetVar struct {
	Name string
}

func (p *GetVar) GetBool(run rt.Runtime) (ret bool, err error) {
	err = run.GetVariable(p.Name, &ret)
	return
}

func (p *GetVar) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = run.GetVariable(p.Name, &ret)
	return
}

func (p *GetVar) GetText(run rt.Runtime) (ret string, err error) {
	err = run.GetVariable(p.Name, &ret)
	return
}

// FIX:
// func (p *GetVar) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
// 	var values []float64
// 	if e := run.GetVariable(p.Name, &values); e != nil {
// 		err = e
// 	} else {
// 		ret = qna.NewNumberList(values)
// 	}
// 	return
// }

// func (p *GetVar) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
// 	var values []string
// 	if e := run.GetVariable(p.Name, &values); e != nil {
// 		err = e
// 	} else {
// 		ret = qna.NewTextList(values)
// 	}
// 	return
// }
