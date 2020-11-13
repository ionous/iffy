package test

type GroupPartition struct {
	// the name of the object this record describes
	// ex. tile X, or tile W ( from a scrabble set )
	Name string
	// objects with the same label are considered to be in the same group.
	// ex. "scrabble tiles"
	Label string
	// by default groups are numerable:
	// the group is prefixed by the number of items in the group.
	// ex. five scrabble tiles.
	Innumerable bool
	// whether and how to print objects.
	ObjectGrouping ObjectGrouping
}

// ObjectGrouping defines how objects in groups should display.
type ObjectGrouping int

//go:generate stringer -type=ObjectGrouping
const (
	// indicates we dont want the individual objects in the group
	// ex. the scrabble tiles, the usual utensils, several things.
	WithoutObjects ObjectGrouping = iota
	// indicates the individual objects shouldnt use articles
	// ex. tiles X and W from a Scrabble set.
	WithoutArticles
	// indicates individual objects should have articles.
	// ex. the X and the W tiles from a Scrabble set.
	WithArticles
)
