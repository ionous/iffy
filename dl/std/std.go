package std

type Classes struct {
	*Kind
	*Room
	*Thing
}

type Patterns struct {
	*PrintName
	*PrintPluralName
	*GroupTogether
	*PrintGroup
}

type Commands struct {
	*PrintNondescriptObjects
}
