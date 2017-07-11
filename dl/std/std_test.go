package std_test

import (
	. "github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestStd(t *testing.T) {
	suite.Run(t, new(StdSuite))
}

type StdSuite struct {
	suite.Suite
}

func (assert *StdSuite) SetupTest() {
}

func (assert *StdSuite) TestNames() {
	classes := ref.NewClasses()
	objects := ref.NewObjects(classes)

	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*Classes)(nil))

	unique.RegisterValues(unique.PanicValues(objects),
		// a named room
		&Room{
			Kind: Kind{Name: "Nowhere"},
		},
		// an unnamed thing
		&Thing{},
		// a named thing
		&Thing{
			Kind: Kind{Name: "pen"},
		},
		// a thing with a printed name
		&Thing{
			Kind: Kind{Name: "sword", PrintedName: "plastic sword"},
		},
		//
		// &Thing{Name: "trevor", CommonProper: ProperNamed},
	)

}
