package story

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
)

func imp_test_statement(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "test_statement"); e != nil {
		err = e
	} else if n, e := imp_test_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else {
		err = reader.Slot(m.MapOf("$TEST"), "testing", reader.ReadMaps{
			"test_output": func(m reader.Map) error {
				return imp_test_output(k, n, m)
			},
		})
	}
	return
}

func imp_test_scene(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "test_scene"); e != nil {
		err = e
	} else if n, e := imp_test_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else {
		lastScene := k.Current.Domain
		k.Current.Domain = n
		err = reader.Repeats(m.SliceOf("$STORY_STATEMENT"), k.Bind(imp_story_statement))
		k.Current.Domain = lastScene
	}
	return
}

func imp_test_rule(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "test_rule"); e != nil {
		err = e
	} else if testName, e := imp_test_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else if hook, e := imp_program_hook(k, m.MapOf("$HOOK")); e != nil {
		err = e
	} else if prog, e := k.NewProg(hook.SlotType(), hook.CmdPtr()); e != nil {
		err = e
	} else {
		k.NewTestProgram(testName, prog)
	}
	return
}

// return expectation
func imp_test_output(k *Importer, testName ephemera.Named, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "test_output"); e != nil {
		err = e
	} else if expect, e := imp_lines(k, m.MapOf("$LINES")); e != nil {
		err = e
	} else {
		k.NewTestExpectation(testName, "execute", expect)
	}
	return
}

func imp_lines(k *Importer, r reader.Map) (ret string, err error) {
	return reader.String(r, "lines")
}
