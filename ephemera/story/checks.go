package story

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
)

func imp_test_statement(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "test_statement"); e != nil {
		err = e
	} else if n, e := imp_test_name(k, m.MapOf("$TEST_NAME")); e != nil {
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
	} else if n, e := imp_test_name(k, m.MapOf("$TEST_NAME")); e != nil {
		err = e
	} else {
		pop := k.SetCurrentTest(n)
		err = reader.Repeats(m.SliceOf("$STORY_STATEMENT"), k.Bind(imp_story_statement))
		pop()
	}
	return
}

func imp_test_rule(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "test_rule"); e != nil {
		err = e
	} else if testName, e := imp_test_name(k, m.MapOf("$TEST_NAME")); e != nil {
		err = e
	} else if hook, e := imp_program_hook(k, m.MapOf("$HOOK")); e != nil {
		err = e
	} else if prog, e := k.NewGob(hook.SlotType(), hook.CmdPtr()); e != nil {
		err = e
	} else {
		k.NewTestProgram(testName, prog)
	}
	return
}

func imp_test_name(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	err = reader.Option(r, "test_name", reader.ReadMaps{
		"$CURRENT_TEST": func(m reader.Map) (err error) {
			// we dont have to parse current test, just its existence is enough
			ret = k.StoryEnv.Recent.Test
			return
		},
		"$NAMED_TEST": func(m reader.Map) (err error) {
			ret, err = imp_named_test(k, m)
			return
		},
	})
	return
}

func imp_named_test(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	if m, e := reader.Unpack(r, "named_test"); e != nil {
		err = e
	} else {
		ret, err = imp_test_text(k, m.MapOf("$NAME"))
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
