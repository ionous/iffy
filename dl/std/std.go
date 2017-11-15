package std

import (
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/dl/std/group"
)

type Classes struct {
	*Kind
	*Room
	*Thing
	*Actor
	*Container
	*Pawn
	*Story
}

// FIX? maybe the expected return of the pattern can be specified here via struct tags.
type Patterns struct {
	*Commence
	*CommonObjects
	*ConstructStatusLine
	*DescribeObject
	*group.GroupTogether
	*IsNotableEnclosure
	*IsNotableScenery
	*IsUnremarkable
	*NotableObjects
	*PlayerSurroundings
	*PrintBannerText
	*PrintContent
	*PrintLocation
	*PrintName
	*PrintObject
	*group.PrintGroup
	*PrintPluralName
	*PrintSeveral
	*PrintSummary
	*VisibleParents
	*Children
	*Parents

	// *DescribeFirstRoom
	// *EndTurn
	// *PrintBanner
	// *SetInitialPosition
	// *StartTurn
	// *UpdateScore
}

type Commands struct {
	// Runtime
	*DescribeLocation
	*locate.LocationOf
	*LowerAn
	*LowerThe
	*Pluralize
	*PrintNondescriptObjects
	*PrintObjects
	*UpperAn
	*UpperThe
	// Pluralizer
	*PluralRule
	*Player
}
