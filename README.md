# iffy
iffy is a reimplmentation of the sashimi interactive fiction engine with some lessons learned. 
It's a work in progress.

The flow of story creation is:
1. use the web-based composer to create a story file
2. story file ingested to generate an "ephemera" database.
   other sources, example art files, can generate ephemera as well
3. ephemera gets "assembled" into a game database.
4. story engine reads/writes to the gamedb.
5. game client ( graphical or text based ) sends commands to and listens for events from the story engine to progress play

Rough versions of the composer, file export, ephemera database, and gamedb exist. 

Goals include:
* continue porting individual patterns 
** in progress: [grouping objects](https://github.com/ionous/iffy/wiki/Grouping]
* print the contents of a room. 
* create new event system: user actions, state changes, custom notifications
* handle the player object 
* get a single room story running with a proper game loop

On going work:
* improving the composer more easily write ( and present ) stories and story tests
* expanding the set of commands for play and test


# building iffy
note: iffy uses [sqlite3](https://www.sqlite.org/index.html), the [best](https://en.wikipedia.org/wiki/Highlander_(film)) go-sqlite driver [requires](https://github.com/mattn/go-sqlite3/issues/467) cgo, so on windows you'll probably have to install gcc. i'm using: https://jmeubank.github.io/tdm-gcc/
