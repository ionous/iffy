package std

type Classes struct {
	*Kind
	*Room
	*Thing
	*Actor
}

type Patterns struct {
	*PrintName
	*PrintPluralName
	*GroupTogether
	*PrintGroup
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
