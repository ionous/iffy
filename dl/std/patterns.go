package std

// PrintName, by default, prints the "printed name" property of the target,
// or the object name ( if the target lacks a "printed name" ),
// or the object's class name ( for unnamed objects. )
// A "printed name" can change during the course of play; object names never change.
type PrintName struct {
	Target *Kind
}

// PrintPluralName, by default, prints the "printed plural name" property of the target,
// or, if the target lacks a "printed plural name", the plural of the printed name
// from the pluralization table. It uses automated pluralization if there is no such pluralized entry.
type PrintPluralName struct {
	Target *Kind
}

type Patterns struct {
	*PrintName
	*PrintPluralName
	// *GroupTogether
	// *PrintGroup
}

type Commands struct {
	// *PrintNondescriptObjects
}
