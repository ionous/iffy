/* ---------------------------------------------------
test data for web
 --------------------------------------------------- */
insert into mdl_kind(kind, path) values
	('things', ''),
	('vehicles', 'things'),
	('cars', 'vehicles,things'),
	('people', 'things');

insert into mdl_field(kind, field, type) values
	('things', 'brief', 'text'),
	('vehicles', 'flightiness', 'aspect'),
	('cars', 'num wheels', 'digi');

insert into mdl_default(kind, field, value) values
	('cars', 'flightiness', 'flightless'),
	('cars', 'num wheels', 4);

insert into mdl_aspect(aspect, trait) values
	('flightiness', 'flightless'),
	('flightiness', 'glide worthy'),
	('flightiness', 'flight worthy');

insert into mdl_rel(relation, kind, cardinality, otherKind) values
	('containing', 'vehicles', 'one_any', 'people');

insert into mdl_noun(noun, kind) values
	('dune buggy', 'cars'),
	('enterprise', 'vehicles'),
	('riker', 'people'),
	('picard', 'people');

insert into mdl_name(noun, name, rank) values
	('dune buggy', 'dune buggy', 0),
	('dune buggy', 'dune', 1),
	('dune buggy', 'buggy', 1),
	('enterprise', 'enterprise', 0),
	('riker', 'riker', 0),
	('picard', 'picard', 0);

insert into mdl_start(noun, field, value) values
	( 'dune buggy', 'num wheels', 3);

insert into mdl_pair(noun, relation, otherNoun) values
	( 'dune buggy', 'containing', 'picard'),
	( 'enterprise', 'containing', 'riker');

insert into mdl_spec(type, name, spec) values
	('kind',  'things', 'From inform: ''Represents anything interactive in the world. People, pieces of scenery, furniture, doors and mislaid umbrellas might all be examples, and so might more surprising things like the sound of birdsong or a shaft of sunlight.'''),
	('aspect', 'flightiness', 'The flight worthiness of vehicles, an example of an aspect with several traits.'),
	('trait', 'glide worthy', 'Better at landing than taking off.'),
	('relation', 'containing', 'The outside of insides.'),
	('field', 'cars.num wheels', 'Not all cars are created equal, or even even.');

