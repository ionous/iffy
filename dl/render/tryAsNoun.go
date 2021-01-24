package render

// TryAsNoun provides a backdoor for templates.
// currently they can't analyze the context they're being called in, so:
// if the command is being used to get an object value
// and no variable can be found in the current context of the requested name,
// see if the requested name is an object instead.
type TryAsNoun int

const (
	TryAsVariable TryAsNoun = 1 << iota
	TryAsObject
	TryAsBoth = (TryAsVariable | TryAsObject)
)

// the default setting is to try as a variable.
// however, if we've explicitly set "TryAsObject" we also need to explicitly set "TryAsVariable" if we want that.
func (flags TryAsNoun) tryVariable() bool {
	return (flags == 0) || (flags&TryAsVariable != 0)
}

func (flags TryAsNoun) tryObject() bool {
	return flags&TryAsObject != 0
}
