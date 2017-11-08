package text

import (
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
	r "reflect"
)

type Commander interface {
	CreateName(string) (string, error)
	CreateCommand(string) (*ops.Command, error)
	EmplaceCommand(i interface{}) (*ops.Command, error)
	CreateExpression(postfix.Expression, r.Type) (*ops.Command, error)
}

func ConvertDirectives(c Commander, dirs []template.Directive) (*ops.Command, error) {
	eng := Engine{Commander: c}
	return eng.convert(dirs)
}
