package object

const Prefix = '$' // leading character used for all internal targets

// internal targets for GetField
const Aspect = "$aspect"   // name of aspect for noun.trait
const Counter = "$counter" // sequence counter
const Value = "$value"     // returns the object rt.Value
const Variables = "$var"   // named values, controlled by scope, not associated with any particular object

// internal fields for object
const Name = "$name"   // name of an object as declared by the user
const Kind = "$kind"   // type of a game object
const Kinds = "$kinds" // hierarchy of a game object ( a path )
