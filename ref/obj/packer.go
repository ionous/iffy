package obj

import (
	r "reflect"
)

type Packer interface {
	Pack(dst, src r.Value) error
}
