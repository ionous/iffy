package ref

// import (
// 	. "github.com/ionous/iffy/tests"
// 	"github.com/stretchr/testify/assert"
// 	r "reflect"
// 	"strings"
// 	"testing"
// )

// // TestEnumInverse verifies we can set a binary enumeration true/false
// func TestEnumInverse(t *testing.T) {
// 	assert := assert.New(t)
// 	var secretState SecretState
// 	if enum, e := MakeEnum(&secretState); assert.NoError(e) {
// 		assert.EqualValues(Secret, secretState)
// 		pv := r.ValueOf(&secretState).Elem()
// 		if e := enum.setValue(pv, int(secretState), false); assert.NoError(e) {
// 			assert.EqualValues(NotSecret, secretState)
// 			if e := enum.setValue(pv, int(secretState), false); assert.NoError(e) {
// 				assert.EqualValues(Secret, secretState)
// 			}
// 		}
// 	}
// }
