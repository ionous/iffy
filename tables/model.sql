/**
 * tables describing the the world and its rules.
 */

/* documentation for pieces of the model: kinds, nouns, fields, etc. */
create table mdl_spec(type text, name text, spec text, primary key(type, name));
/* enumerated values used by kinds and nouns */
create table mdl_aspect(aspect text, trait text, rank int, primary key(aspect, trait));
/* a class of objects with shared characteristics */
create table mdl_kind(kind text, path text, primary key(kind));
/* properties of a kind. type is a PRIM */
create table mdl_field(kind text, field text, type text, primary key(kind, field));
/* default values for the field of a kind ( or one of its descendant kinds ) */
create table mdl_default(kind text, field text, value blob );
/* relation and constraint between two kinds of nouns */
create table mdl_rel(relation text, kind text, cardinality text, otherKind text, primary key(relation));
/* a person, place, or thing in the world. */
create table mdl_noun(noun text, kind text, primary key(noun));
/* words which refer to nouns. in cases where two words may refer to the same noun, 
   the lower rank of the association wins. */
create table mdl_name(noun text, name text, rank int);
/* relation between two specific nouns. these change over the course of a game. */
create table mdl_pair(noun text, relation text, otherNoun text);
/* initial values for various noun properties. these change over the course of a game. */
create table mdl_start(noun text, field text, value blob);
/* stored programs, a work in progress */ 
create table mdl_prog( idSource int, type text, bytes blob );
/* stored tests, a work in progress 
 * tbd: should named test be turned back into text?
 * still need to determine if "sources" should be listed in model ( for debugging )
 */ 
create table mdl_check( name text, idProg int, expect text );