package next

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// GetVariable writes the value at 'name'
type GetVar struct {
	Name string
}

func (*GetVar) Compose() composer.Spec {
	return composer.Spec{
		Name:  "get_var",
		Spec:  "Get var $1",
		Group: "variables",
		Desc:  "Get Variable: Return the value of the named variable.",
	}
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

func (p *GetVar) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	var values []float64
	if e := run.GetVariable(p.Name, &values); e != nil {
		err = e
	} else {
		ret = stream.NewNumberList(values)
	}
	return
}

func (p *GetVar) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	var values []string
	if e := run.GetVariable(p.Name, &values); e != nil {
		err = e
	} else {
		ret = stream.NewTextList(values)
	}
	return
}
