// Package clog wraps a spec block to record the process of building scripts.
package clog

import (
	"fmt"
	"github.com/ionous/iffy/spec"
	"io"
)

// Make wraps the passed block, logging all calls through the passed writer.
func Make(w io.Writer, c spec.Block) spec.Block {
	l := &logger{w: w}
	return clog{l, c, c}
}

type clog struct {
	l *logger
	c spec.Block
	p spec.Slot
}

type logger struct {
	w io.Writer
	t string // shared tabs
}

func (f clog) Begin() bool {
	f.l.Log("{")
	f.l.Indent()
	return f.c.Begin()
}
func (f clog) Param(field string) spec.Slot {
	f.l.Log("Param", field)
	p := f.c.Param(field)
	return clog{l: f.l, p: p}
}
func (f clog) End() {
	f.l.Dedent()
	f.l.Log("}")
	f.c.End()
}

// Cmd starts a new command
func (f clog) Cmd(name string, args ...interface{}) spec.Block {
	if len(args) > 0 {
		if _, ok := args[0].(spec.Block); ok {
			args = []interface{}{len(args), " cmd/s"}
		}
	}
	f.l.Log("Cmd "+name, args...)
	b := f.p.Cmd(name, args...)
	return clog{f.l, b, b}
}

func (f clog) Cmds(cmds ...spec.Block) spec.Block {
	if len(cmds) > 0 {
		f.l.Log("Cmds", len(cmds), " cmd/s")
	} else {
		f.l.Log("Cmds")
	}
	b := f.p.Cmds(cmds...)
	return clog{f.l, b, b}
}

func (f clog) Val(val interface{}) spec.Block {
	f.l.Log("Val", val)
	b := f.p.Val(val)
	return clog{f.l, b, b}
}

func (l *logger) Log(n string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Fprintln(l.w, l.t, n, fmt.Sprint(args...))
	} else {
		fmt.Fprintln(l.w, l.t, n)
	}
}

func (l *logger) Indent() {
	l.t += " "
}

func (l *logger) Dedent() {
	if cnt := len(l.t); cnt > 0 {
		l.t = l.t[:cnt-1]
	}
}
