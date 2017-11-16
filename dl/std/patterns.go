package std

import (
	"github.com/ionous/iffy/ident"
)

type Commence struct {
	Story ident.Id `if:"cls:story"`
}

// CommonObjects should act as a source for those objects which a player "can also see" in a given location.
// see also: DescribeLocation.
type CommonObjects struct {
	Location ident.Id `if:"cls:room"`
}

// ConstructStatusLine should set the values of the story's left and right status text.
type ConstructStatusLine struct {
	Story ident.Id `if:"cls:story"`
}

type DescribeFirstRoom struct {
	Story ident.Id `if:"cls:story"`
}

// DescribeObject should print a brief description about the targeted object.
// see also: DescribeLocation.
type DescribeObject struct {
	Object ident.Id `if:"cls:kind"`
}

type EndTurn struct {
	Story ident.Id `if:"cls:story"`
}

// IsCeiling should return true if the targeted object stops visibility.
type IsCeiling struct {
	Object ident.Id `if:"cls:kind"`
}

// IsNotableEnclosure should return true if the targeted object is important enough to merit a description when it acts as a parent of the player.
// see also: DescribeLocation.
type IsNotableEnclosure struct {
	Object ident.Id `if:"cls:kind"`
}

// IsNotableScenery should return true if the targeted object is important enough to merit a description when defined as scenery in a location.
// see also: DescribeLocation.
type IsNotableScenery struct {
	Object ident.Id `if:"cls:kind"`
}

// IsUnremarkable should return true to avoid being treated as a notable object.
// see also: DescribeLocation.
type IsUnremarkable struct {
	Object ident.Id `if:"cls:kind"`
}

// NotableObjects should return a list of objects in the targeted location which merit a brief description.
// see also: DescribeObject, DescribeLocation, UnremarkableObjects
type NotableObjects struct {
	Location ident.Id `if:"cls:kind"`
}

// PrintCommonObjects should print a simple sentence regarding the associated objects.
type PrintCommonObjects struct {
	Objects []ident.Id
}

// PrintLocation should print the complete position of the player.
// For instance: "The lab (in the box) (on the desk)".
// see also: PlayerSurroundings.
type PrintLocation struct {
	Location ident.Id `if:"cls:kind"`
}

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

// PrintObject defines a pattern to say simple information about an object.
// Inform's "Printing the locale description" uses I6 functions ( "list the contents of" and "say [a list of including contents]" ) which are not themselves activities, though -- apparently -- it sometimes breaks back out into activities: via "details of something".
// There's currently no strong rule for what is a command and what is a pattern.
// http://inform7.com/learn/man/WI_18_25.html
// http://inform7.com/learn/man/WI_18_26.html
type PrintObject struct {
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

// PlayerSurroundings should return the name of the player's current location ( as text. )
type PlayerSurroundings struct{}

// PrintBannerText, by default, says the story's:
//  . title, or "Welcome"
//  . headline and author, or "An interactive fiction"
//  . release number, major.minor.patch
//  . date of compilation. (FIX: not implemented)
//  . version of iffy
//  ex. Release 1 / 050630 / Iffy 1.0
type PrintBannerText struct {
	Story ident.Id
}
type SetInitialPosition struct {
	Story ident.Id
}
type StartTurn struct {
	Story ident.Id
}
type UpdateScore struct {
	Story ident.Id
}

// NotableObjects should return a list of objects in the targeted location which merit a casaul mention.
// see also: DescribeObject, DescribeLocation, NotableObjects.
type UnremarkableObjects struct {
	Location ident.Id
}

// VisibleParents should return a list of all parents of the targeted object, stopping if one of the parents acts as a "ceiling". Examples of ceilings include: rooms and closed boxes.
type VisibleParents struct {
	Object ident.Id
}
