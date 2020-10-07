package object

// reserved fields
const Name = "name" // name of an object as declared by the user

// internal targets
const Prefix = '$'         // leading character used for all internal targets
const Id = "$id"           // unique identifier for an object, includes its home domain
const Exists = "$exists"   // whether a name refers to a declared game object
const Kind = "$kind"       // type of a game object
const Kinds = "$kinds"     // hierarchy of a game object ( a path )
const Counter = "$counter" // sequence counter
const Aspect = "$aspect"   // name of aspect for noun.trait
const Variables = "$var"   // named values, controlled by scope, not associated with any particular object
