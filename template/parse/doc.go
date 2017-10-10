// Package parse is a port of the BSD-licensed golang text/template/parse package
// with all the necessary changes to support iffy's smaller, custom template syntax.
//
// Major changes include:
//  * Interfaces and structs for states instead of functions;
//  * Arrays instead of channels;
//  * Granular files and classes; including the separation of item data into its own package;
//  * All of the template language differences needed for iffy.
//
// Package text/template's choice of state functions limits the ability to compare states, log states, and store state relevant data when necessary.
// To avoid the need for state data, the original code simulates state execution inside of states using custom for loops. It therefore requires channels to avoid blocking.
//
// It's probably better not to *require* channels for lexing;
// rather, just ensure the algorithm can support channels when needed.
//
package parse
