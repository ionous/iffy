package parser

type Commands struct {
	*Action
	*AllOf
	*AnyOf
	*Focus
	*Multi
	*Noun
	*Target
	*Word
}
