/**
 * for saving, restoring a player's game session.
 */
create table if not exists 
	run_domain( domain text, active int, primary key( domain )); 

create table if not exists 
	run_pair( noun text, relation text, otherNoun text, active int default 1, unique( noun, relation, otherNoun ) ); 
