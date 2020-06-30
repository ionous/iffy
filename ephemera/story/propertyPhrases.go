package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

// make.opt("property_phrase", "{primitive_phrase} or {aspect_phrase}");
func imp_property_phrase(k *imp.Porter, kind ephemera.Named, r reader.Map) (err error) {
	return reader.Option(r, "property_phrase", reader.ReadMaps{
		"$PRIMITIVE_PHRASE": func(m reader.Map) (err error) {
			return imp_primitive_phrase(k, kind, m)
		},
		"$ASPECT_PHRASE": func(m reader.Map) (err error) {
			return imp_aspect_phrase(k, kind, m)
		},
	})
}

// fix -- really this should allow object types, eval types and primitive types, etc.
// make.run("primitive_phrase", "{primitive_type} called {property}");
func imp_primitive_phrase(k *imp.Porter, kind ephemera.Named, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "primitive_phrase"); e != nil {
		err = e
	} else if prop, e := imp_property(k, m.MapOf("$PROPERTY")); e != nil {
		err = e
	} else if prim, e := imp_primitive_prop(k, m.MapOf("$PRIMITIVE_TYPE")); e != nil {
		err = e
	} else {
		k.NewField(kind, prop, prim)
	}
	return
}

// make.run("aspect_phrase", "{aspects} {?optional_property}");
func imp_aspect_phrase(k *imp.Porter, kind ephemera.Named, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "aspect_phrase"); e != nil {
		err = e
	} else if aspect, e := imp_aspect(k, m.MapOf("$ASPECT")); e != nil {
		err = e
	} else {
		// inform gives these the name "<noun> condition"
		// we could only do that with an after the fact reduction, and with some additional mdl data.
		// ( ex. in case the same aspect is assigned twice, or twice at difference depths )
		// for now the name of the field is the name of the aspect
		if v := m.MapOf("$OPTIONAL_PROPERTY"); len(v) != 0 {
			err = errutil.New("optional property names aren't supported for aspects")
		} else {
			k.NewField(kind, aspect, tables.PRIM_ASPECT)
		}
	}
	return
}

// make.run("optional_property", "called {property}");
func imp_optional_property(k *imp.Porter, r reader.Map) (ret ephemera.Named, err error) {
	if m, e := reader.Unpack(r, "optional_property"); e != nil {
		err = e
	} else if n, e := imp_property(k, m.MapOf("$PROPERTY")); e != nil {
		err = e
	} else {
		ret = n
	}
	return
}
