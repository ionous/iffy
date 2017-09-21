package std

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/rt"
)

// PrintName defines a pattern to say the target's name.
// The standard rules print the "printed name" property of the target,
// or the object name ( if the target lacks a "printed name" ),
// or the object's class name ( for unnamed objects. )
// A "printed name" can change during the course of play; object names never change.
// Analogous to Inform's "Printing the name of something."
// http://inform7.com/learn/man/WI_18_10.html
type PrintName struct {
	Target ident.Id `if:"cls:kind"`
}

// PrintPluralName defines a pattern to say the plural of the target's name.
// The standard rules print the target's "printed plural name",
// or, if the target lacks that property, the plural of the "print name" pattern.
// It uses the runtime's pluralization table, or if needed, automated pluralization.
// Analogous to Inform's "Printing the plural nam e of something."
// http://inform7.com/learn/man/WI_18_11.html
type PrintPluralName struct {
	Target ident.Id `if:"cls:kind"`
}

// PrintSeveral defines a pattern to say information about a generic group of objects.
// The standard rules print the group size in words, then prints the plural name of the target.
// Analogous to Inform's "Printing a number of something."
// http://inform7.com/learn/man/WI_18_12.html
type PrintSeveral struct {
	Target    ident.Id `if:"cls:kind"`
	GroupSize float64
}

// PrintObject defines a pattern to say simple information about an object.
// Inform's "Printing the locale description" uses I6 functions ( "list the contents of" and "say [a list of including contents]" ) which are not themselves activities, though -- apparently -- it sometimes breaks back out into activities: via "details of something".
// There's currently no strong rule for what is a command and what is a pattern.
// http://inform7.com/learn/man/WI_18_25.html
// http://inform7.com/learn/man/WI_18_26.html
type PrintObject struct {
	Target ident.Id `if:"cls:kind"`
}

// PrintSummary defines a pattern to say extra information for certain items ( such as containers. )
// ( open, closed, worn, being worn, providing light )
// ( providing light and open but empty )
// Similar to Inform's "Printing inventory/room description details of something.
// http://inform7.com/learn/man/WI_18_16.html
// http://inform7.com/learn/man/WI_18_17.html
type PrintSummary struct {
	Target ident.Id `if:"cls:kind"`
}

type PrintContent struct {
	Target            ident.Id `if:"cls:kind"`
	Header            string
	Articles, Tersely bool
}

type PlayerSurroundings struct{}

type Commence struct {
	Story rt.Object
}

// PrintBannerText, by default, says the story's:
//  . title, or "Welcome"
//  . headline and author, or "An interactive fiction"
//  . release number, major.minor.patch
//  . date of compilation. (FIX: not implemented)
//  . version of iffy
//  ex. Release 1 / 050630 / Iffy 1.0
type PrintBannerText struct {
	Story rt.Object
}
type ConstructStatusLine struct {
	Story rt.Object
}
type DescribeFirstRoom struct {
	Story rt.Object
}
type EndTurn struct {
	Story rt.Object
}
type SetInitialPosition struct {
	Story rt.Object
}
type StartTurn struct {
	Story rt.Object
}
type UpdateScore struct {
	Story rt.Object
}
