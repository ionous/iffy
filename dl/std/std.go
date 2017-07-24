package std

import (
	"github.com/ionous/iffy/dl/std/group"
)

type Classes struct {
	*Kind
	*Room
	*Thing
	*Actor
}

type Patterns struct {
	*group.GroupTogether
	*group.PrintGroup
	*PrintName
	*PrintPluralName
}

type Commands struct {
	// Runtime
	*PrintNondescriptObjects
	*UpperThe
	*LowerThe
	*UpperAn
	*LowerAn
	*Pluralize
	// Pluralizer
	*PluralRule
}

// PrintName executes a pattern to print the target's name.
// The standard rules print the "printed name" property of the target,
// or the object name ( if the target lacks a "printed name" ),
// or the object's class name ( for unnamed objects. )
// A "printed name" can change during the course of play; object names never change.
type PrintName struct {
	Target *Kind
}

// PrintPluralName executes a pattern to print the plural of the target's name.
// The standard rules print the target's "printed plural name",
// or, if the target lacks that property, the plural of the "print name" pattern.
// It uses the runtime's pluralization table, or if needed, automated pluralization.
type PrintPluralName struct {
	Target *Kind
}
