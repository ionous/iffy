# iffy
iffy is a reimplmentation of the sashimi interactive fiction engine with some lessons learned. 

The flow of story creation is:
1. Use the web-based "Composer" to write stories and supporting scripts.
2. Let the tools generate "ephemera" from the story files
  ( other sources, for instance art files, can generate ephemera too. )
3. Use the tools to "assemble" a game database from the ephemera.
4. The "Story Engine" reads and writes to the gamedb during play.
5. A "Game Client" ( command line, browser based like ["Alice and the Galactic Traveler"](https://evermany.itch.io/alice), or someday Unity, etc. ) sends commands to ( and listens for events from ) the story engine to progress play.

Rough versions of the iffy Composer, ephemera and gamedb exist. 

Current goals include:
* porting parent-child containment ( ex. rocks in a box, or people in a room. )
* printing the contents of a room. ( ex. `> look.` )
* creating new event system: user actions, state changes, custom notifications
* handling the player object 
* getting a single room story running with a proper game loop

Ongoing work includes:
* improving the composer more easily write stories and story tests
* expanding the abilities of the story engine to support unit tests and game play.

I'm going to try to commit to keeping track of progress on the [iffy wiki](https://github.com/ionous/iffy/wiki). Nothing is probably usable by other people yet. It is a work in progress.

# building iffy
note: iffy uses [sqlite3](https://www.sqlite.org/index.html), the [best](https://en.wikipedia.org/wiki/Highlander_(film)) go-sqlite driver [requires](https://github.com/mattn/go-sqlite3/issues/467) cgo, so on windows you'll probably have to install gcc. i'm using: https://jmeubank.github.io/tdm-gcc/
