package reflector

import (
	. "github.com/ionous/iffy/tests"
	"github.com/stretchr/testify/assert"
	r "reflect"
	"strings"
	"testing"
)

// TestEnumChoices verifies a stringerified enum generates good choices.
func TestEnumChoices(t *testing.T) {
	assert := assert.New(t)
	if enum, e := MakeEnum((*TriState)(nil)); assert.NoError(e, "enum should generate") {
		var reduce []string
		for _, c := range enum {
			reduce = append(reduce, c.Name())
		}
		assert.EqualValues("NoYesMaybe", strings.Join(reduce, ""))
		assert.EqualValues(0, enum.ChoiceToIndex("No"))
		assert.EqualValues(-1, enum.ChoiceToIndex("Never"))
		assert.EqualValues("Maybe", enum.IndexToChoice(2).Name())
		assert.EqualValues("$maybe", enum.IndexToChoice(2).Id())
		assert.Empty(enum.IndexToChoice(7).IsValid())
		assert.Empty(enum.IndexToChoice(-1).IsValid())
	}
}

// TestEnumEmpty verifies an empty enum is an error.
func TestEnumEmpty(t *testing.T) {
	assert := assert.New(t)
	_, e := MakeEnum((*EmptyState)(nil))
	assert.Error(e, "enum shouldnt generate")
}

// TestEnumTooLong verifies an enum with too many choices is an error.
func TestEnumTooLong(t *testing.T) {
	assert := assert.New(t)
	_, e := MakeEnum((*TooLongState)(nil))
	assert.Error(e, "enum shouldnt generate")
}

// TestEnumInverse verifies we can set a binary enumeration true/false
func TestEnumInverse(t *testing.T) {
	assert := assert.New(t)
	var secretState SecretState
	if enum, e := MakeEnum(&secretState); assert.NoError(e) {
		assert.EqualValues(Secret, secretState)
		pv := r.ValueOf(&secretState).Elem()
		if e := enum.setValue(pv, int(secretState), false); assert.NoError(e) {
			assert.EqualValues(NotSecret, secretState)
			if e := enum.setValue(pv, int(secretState), false); assert.NoError(e) {
				assert.EqualValues(Secret, secretState)
			}
		}
	}
}
