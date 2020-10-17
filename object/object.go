package object

// internal targets for GetField
const Prefix = '$' // leading character used for all internal targets
//
const Aspect = "$aspect"   // name of aspect for noun.trait
const Counter = "$counter" // sequence counter
const Exists = "$exists"   // whether a name refers to a declared game object
const Id = "$id"           // unique identifier for an object, includes its home domain
const Kind = "$kind"       // type of a game object
const Kinds = "$kinds"     // hierarchy of a game object ( a path )
const Name = "$name"       // name of an object as declared by the user
const Variables = "$var"   // named values, controlled by scope, not associated with any particular object
