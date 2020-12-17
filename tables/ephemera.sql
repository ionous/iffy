/* alternative name for a noun */
create table eph_alias( idNamedAlias int, idNamedActual int );
/* collection of related states ( called traits ) */
create table eph_aspect( idNamedAspect int );
/* likelihood that a trait applies to a particular kind */
create table eph_certainty( certainty text, idNamedTrait int, idNamedKind text );
/* test programs and the results they are expected to produce */
create table eph_check( idNamedTest int, idProg int );
/* initial values for the properties of nouns of belonging to the specified kind */
create table eph_default( idNamedKind int, idNamedProp int, value blob );
/* test programs and the results they are expected to produce */
create table eph_expect( idNamedTest int, testType text, expect text );
/* property name and type associated with a kind of object */
create table eph_field( idNamedKind int, idNamedField int, primType text, primAff text );
/* collection of related nouns, plural named kind, singular named parent
 * ex. cats are a kind of animal.
 */
create table eph_kind( idNamedKind int, idNamedParent int );
/* user specified appellation and the location that specification came from 
 * domain references another name
 */
create table eph_named( name text, og text, category text, domain int, idSource int, offset text );
/* a named object in the game world and its kind (singular) */
create table eph_noun( idNamedNoun int, idNamedKind int );
/* declarations and references to pattern parameter and pattern return types.
if idNamedPattern is idNamedParam it indicates a return type.
idNamedParam can be in the kind category, or one of the predefined eval types .
idProg holds program initialization for locals ( and possibly default parameter values )
for now idProg < 0 is a pattern or pattern parameter reference;
the "category" of idNamedParam also indicates reference vs. declaration
*/
create table eph_pattern( idNamedPattern int, idNamedParam int, idNamedType int, affinity text, idProg int );
/* rule for the collective name of a singular word */
create table eph_plural( idNamedPlural int, idNamedSingluar int );
/* type is the name of the command container for de-serialization of the prog
 * idSource is the origin of the ephemera just like eph_named.
 * fix? unclear why there's no offset, and if here is domain, source, offset concerns
 * why that isnt a table of its own for declarations to share.
 */
create table eph_prog( idSource int, progType text, prog blob );
/* connection between two kinds of object */
create table eph_relation( idNamedRelation int, idNamedKind int, idNamedOtherKind int, cardinality text check (cardinality in ('one_one','one_any','any_one','any_any')));
/* connection between two object instances */	
create table eph_relative( idNamedHead int, idNamedStem int, idNamedDependent int );
/* function handler for a pattern */
create table eph_rule( idNamedPattern int, idProg int );/* uri, file name or other identification for the origin of the various ephemera recorded in the db. 
while its not particularly useful to have a one column primitive data type column
someday, this might contain source modification times or other useful info.
*/
create table eph_source( src text );
/* only one trait from a given aspect can be true for a noun at a time. */	
create table eph_trait( idNamedTrait int, idNamedAspect int, rank int );
/* initial value for a noun's field, trait, or aspect 
 * fix? why does there need to be both value and default?
 * the names ( idNamedNoun vs idNamedKind ) should be separation enough.
 */
create table eph_value( idNamedNoun int, idNamedProp int, value blob );
/* word indicating a particular relationship between nouns */
create table eph_verb( idNamedStem int, idNamedRelation int, verb text );
