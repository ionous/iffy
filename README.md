# iffy
iffy is a reimplmentation of the sashimi interactive fiction engine with some lessons learned. 
It's a work in progress.

Current work involves re-implementing the sashimi compiler ( more accurately called an assembler. it doesn't compile code, instead it builds a working model of a game world from a description of that world. )

The basic flow is visual story editor (aka composer) -> story file -> ephemera database -> story database -> game.

Rough versions of the composer, file export, ephemera database, and game exist. Remaining steps include:
* assemble story database from ephemera ( esp. patterns and event handling )
* adapt existing game code to use the story database
* rework hosting environment to run games
* port story libraries so that games can actually do interesting things ( move from room to room, open doors,  etc. )

# building iffy
note: iffy uses [sqlite3](https://www.sqlite.org/index.html), the [best](https://en.wikipedia.org/wiki/Highlander_(film)) go-sqlite driver [requires](https://github.com/mattn/go-sqlite3/issues/467) cgo, so on windows you'll probably have to install gcc. i'm using: https://jmeubank.github.io/tdm-gcc/
