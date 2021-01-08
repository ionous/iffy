# iffy
This is a reimplmentation of the Sashimi interactive fiction engine with some lessons learned. 

The flow of story creation is:
1. Use the web-based "Composer" to write stories and supporting scripts.
2. Use the tools to: 
    - first, generate "ephemera" from the story files ( other sources -- ie. art assets --- can generate ephemera, too. )
    - second, to "assemble" a game database from the ephemera.
3. The "Story Engine" reads and writes to the gamedb during play.
4. A "Game Client" then sends commands to ( and listens for events from ) the story engine to progress play.
    - Clients can be command line like traditional interactive fiction;
    - Custom like ["Alice and the Galactic Traveler"](https://evermany.itch.io/alice) which used Sashimi's engine;
    - or someday Unity, etc.

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
