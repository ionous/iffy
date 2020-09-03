package story

import (
	"github.com/ionous/iffy/ephemera/reader"
)

// opt("program_hook", "{activity} or {result:program_return}");
func imp_program_hook(k *Importer, r reader.Map) (ret programHook, err error) {
	err = reader.Option(r, "program_hook", reader.ReadMaps{
		"$ACTIVITY": func(m reader.Map) (err error) {
			// unexpected type: type wanted "execute" got "activity" at
			// ahhh... this isnt a "slot" its a pointer to a specific commnd
			var slot executeSlot
			if e := k.DecodeAny(m, &slot.cmd); e != nil {
				err = e
			} else {
				ret = &slot
			}
			return
		},
		"$RESULT": func(m reader.Map) (err error) {
			ret, err = imp_program_return(k, m)
			return
		},
	})
	return
}

// run("program_return", "return {result:program_result}");
// note: this slat exists for composer formatting reasons only...
func imp_program_return(k *Importer, r reader.Map) (ret programHook, err error) {
	if m, e := reader.Unpack(r, "program_return"); e != nil {
		err = e
	} else {
		ret, err = imp_program_result(k, m.MapOf("$RESULT"))
	}
	return
}

// opt("program_result", "a {simple value%primitive:primitive_func} or an {object:object_func}");
func imp_program_result(k *Importer, r reader.Map) (ret programHook, err error) {
	err = reader.Option(r, "program_result", reader.ReadMaps{
		"$PRIMITIVE": func(m reader.Map) (err error) {
			ret, err = imp_primitive_func(k, m)
			return
		},
		"$OBJECT": func(m reader.Map) (err error) {
			ret, err = imp_object_func(k, m)
			return
		},
	})
	return
}

// opt("primitive_func", "{a number%number_eval}, {some text%text_eval}, {a true/false value%bool_eval}")
func imp_primitive_func(k *Importer, r reader.Map) (ret programHook, err error) {
	err = reader.Option(r, "primitive_func", reader.ReadMaps{
		"$NUMBER_EVAL": func(m reader.Map) (err error) {
			var slot numberSlot
			if e := k.DecodeSlot(m, slot.SlotType(), &slot.cmd); e != nil {
				err = e
			} else {
				ret = &slot
			}
			return
		},
		"$TEXT_EVAL": func(m reader.Map) (err error) {
			var slot textSlot
			if e := k.DecodeSlot(m, slot.SlotType(), &slot.cmd); e != nil {
				err = e
			} else {
				ret = &slot
			}
			return
		},
		"$BOOL_EVAL": func(m reader.Map) (err error) {
			var slot boolSlot
			if e := k.DecodeSlot(m, slot.SlotType(), &slot.cmd); e != nil {
				err = e
			} else {
				ret = &slot
			}
			return
		},
	})
	return
}

// run("object_func", "an object named {name%text_eval}");
func imp_object_func(k *Importer, r reader.Map) (ret programHook, err error) {
	if m, e := reader.Unpack(r, "object_func"); e != nil {
		err = e
	} else {
		// FIX: we should wrap the text with a runtime "test the object matches" command
		// and, -- for simple text values -- add ephemera for "name, type"
		var slot textSlot
		if e := k.DecodeSlot(m, slot.SlotType(), &slot.cmd); e != nil {
			err = e
		} else {
			ret = &slot
		}
	}
	return
}
