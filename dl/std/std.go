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

type Patterns struct {
	*group.GroupTogether
	*group.PrintGroup
	*PrintName
	*PrintPluralName
	*PrintSeveral
	*PrintObject
	*PrintSummary
	*PrintContent
	*Commence
	*PlayerSurroundings
	*PrintBannerText
	*ConstructStatusLine
	// *DescribeFirstRoom
	// *EndTurn
	// *PrintBanner
	// *SetInitialPosition
	// *StartTurn
	// *UpdateScore
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
	*Player
	*locate.LocationOf
}
