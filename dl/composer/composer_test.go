package composer

import "testing"

type Name struct{}
type PtrThing struct{}
type Direct string

func (*Name) Compose() Spec {
	return Spec{
		Name: "named",
	}
}

func (*PtrThing) Compose() Spec {
	return Spec{
		// Name: "named",
	}
}

func (*Direct) Compose() Spec {
	return Spec{
		// Name: "named",
	}
}

func TestNames(t *testing.T) {
	if want, have := "named", SpecName((*Name)(nil)); want != have {
		t.Errorf("have %q, want %q", have, want)
	} else if want, have := "ptr_thing", SpecName((*PtrThing)(nil)); want != have {
		t.Errorf("have %q, want %q", have, want)
	} else if want, have := "direct", SpecName((*Direct)(nil)); want != have {
		t.Errorf("have %q, want %q", have, want)
	}
}
