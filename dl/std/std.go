package std

import (
	"github.com/ionous/iffy/dl/std/group"
)

type Classes struct {
	*Kind
	*Room
	*Thing
	*Actor
	*Container
}

type Patterns struct {
	*group.GroupTogether
	*group.PrintGroup
	*PrintName
	*PrintPluralName
	*PrintSeveral
}

type Commands struct {
	// Runtime
	*PrintNondescriptObjects
	*PrintObjects
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
// Analogous to Inform's "Printing the name of something."
// http://inform7.com/learn/man/WI_18_10.html
type PrintName struct {
	Target *Kind
}

// PrintPluralName executes a pattern to print the plural of the target's name.
// The standard rules print the target's "printed plural name",
// or, if the target lacks that property, the plural of the "print name" pattern.
// It uses the runtime's pluralization table, or if needed, automated pluralization.
// Analogous to Inform's "Printing the plural name of something."
// http://inform7.com/learn/man/WI_18_11.html
type PrintPluralName struct {
	Target *Kind
}

// PrintSeveral executes a pattern to print information about a generic group of objects.
// The standard rules print the group size in words, then prints the plural name of the target.
// Analogous to Inform's "Printing a number of something."
// http://inform7.com/learn/man/WI_18_12.html
type PrintSeveral struct {
	Target    *Kind
	GroupSize float64
}
