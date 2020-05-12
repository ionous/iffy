package object

const Exists = "$exists"   // whether a name refers to a declared game object
const Kind = "$kind"       // type of a game object
const Kinds = "$kinds"     // hierarchy of a game object
const Counter = "$counter" // sequence counter

// originally these were just "pattern"
// it helps to have a type hint to GetField
// it might be better to pass the out pointer so GetField can do that.
const BoolRule = "$bool_rule"
const NumberRule = "$number_rule"
const TextRule = "$text_rule"
const ExecuteRule = "$execute_rule"
const NumListRule = "$num_list_rule"
const TextListRule = "$text_list_rule"
