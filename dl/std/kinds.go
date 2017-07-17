package std

//go:generate stringer -type=SingularPlural
type SingularPlural int

const (
	SingularNamed SingularPlural = iota
	PluralNamed
)

//go:generate stringer -type=CommonProper
type CommonProper int

const (
	CommonNamed CommonProper = iota
	ProperNamed
)

// Kind represents people, places, and other nounish elements of the game world.
type Kind struct {
	Name              string `if:"id"` // how the author refers to an instance of Kind. Cannot change.
	PrintedName       string // overrides the author's name, can change if necessary.
	PrintedPluralName string // when there are multiple objects grouped together in the event the default pluralization doesnt work.
	IndefiniteArticle string // ex. some, a clutch.
	SingularPlural
	CommonProper
	// AmbiguouslyPlural
	// ListGroupKey string
}

// Room "Represents geographical locations, both indoor and outdoor, which are not necessarily areas in a building. A player in one room is mostly unable to sense, or interact with, anything in a different room. Rooms are arranged in a map."
type Room struct {
	Kind
	Dark        bool
	Visited     bool
	Description string
}

// Thing "Represents anything interactive in the model world that is not a room. People, pieces of scenery, furniture, doors and mislaid umbrellas might all be examples, and so might more surprising things like the sound of birdsong or a shaft of sunlight."
type Thing struct {
	Kind
	Description string
	Brief       string // known as "initial appearance"

	// this is part of the room display:
	// unmarked for listing not marked for listing,
	// mentioned not unmentioned

	Scenery bool // unmentioned in the room description
	Handled bool // controls use of the initial apperance

	// Usually unlit not lit,
	// inedible not edible,
	// portable not fixed in place,
	// matching key (object).
	// Usually not wearable, pushable between rooms
	//
	// described not undescribed: Used to exclude things from room descriptions ( ex. the player "yourself".) Assumes some other text already reveals or implies the item' existance. If the item isn't intended to move, it's better to make it "scenery".
	// Note: objects become "described" when carried or worn by the player.
	// Note: nothing on top of an "undescribed" supporter will be visible in a room description.

}