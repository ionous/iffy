// Package assign helps with value conversions throughout iffy.
// It suffers from the law of leaky abstractions because
// it has to make some assumptions about the types of variables and the way they are used.
// Using package runtime increases code complexity; while using package reflect feels overkill.
package assign
