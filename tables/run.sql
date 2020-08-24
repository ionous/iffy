/**
 * for saving, restoring a player's game session.
 */
create table if not exists 
	run_domain( domain text, active int, primary key( domain )); 
