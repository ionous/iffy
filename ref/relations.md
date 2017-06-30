relations.md

# tasks

* set relation, table line, sorted injection.
* views: sync in relation to support views
* property tags to check view; exclusion for temporary/id-less objects.
* error checking for view tags
* changes to the view trigger changes in relations

# relation 

set relation is the main thing.
a slice with elements, clearing a relation swaps the line with the last one and shrinks the slice.

no special query syntax right now, just a walk through all relations. each an objct 

the big question is how set relation works. you have to generically walk all objects and evaluate whether this line is for you.

if you use reflect its going to be darn slow. you could, of course, use a sorted index (on each side?) -- but that seems like a bunch of work. something you should delgate to something like sqlite or bolt.

you could use package sort, keep the table sorted on primary (first) and secondary (second) strings. inject as needed. 

the question is, what are you sorting.
i think you could store the strings outside, extract them using reflection -- it doubles your memory, so what. your entry becomes tableline { first, second, pointer to relation instance. }
the question is whether you have two indicies or just one.

one Gremlin has many Rocks
which, oddly, means one gremlin has many entries, and each rock has only one entry.

primary:
    0: claire, boomba, ptr
    1: claire, rocky, ptr
    2: grace,  pume, ptr
    3: hiro,   loofah,  ptr
    4: marja,  petra, ptr
seconday:
    0: boomba, claire, ptr
    1: loofah, hiro, ptr
    2: petra,  marja
    3: pume,   grace
    4: rocky , claire, ptr

-or-, of course, you could have the table lines memory, but the indicies are just indicies into the ttable lines

i want to remove claire/rocky 
first i find in primary 0, then in find in secondary 4
do we have a list of table lines? probably not, because why bother.
traversal can happen through the primary index.

A.  gremlinRocks.Relate{nil, rocky)
//     now rocky is owned by no one

B. gremlinRocks.Relate{claire, nil)
//     now claire owns nothing.

C/ gremlinRocks.Relate{claire, rocky)
//     now claire owns rocky, and no one else does.

step 1 is realizing that we dont have two parameter calls:
A.  gr:= New("gremlin rocks", +/- map[])
    gr.SetValue("rock", "rocky")
    runtime.Relate(gr)
    // now rocky is owned by no one.

get the relation,    
dtermine which index to use.

still not clear to me all the details, and the relation of primary/secondary to one/multi

A. use secondary index, find rocky, 4.
    we know that we only have one owner [ we are "many" ] so we dont expect to find more entires. remove that entry.

B .

C . 

tell me about unique 
claire is not unique
rocky is unique.

// In a one-to-many relation, the one side appears many times, the many side appears once. ex. parent-son, parent-daughter. The many side is considered "unique".

if i add a claire-rocky relation,
its okay that claire exists multiple times.
but, rocky has to be removed.

isnt that odd?
in a one-to-one relation, both are considered unique and they are the one side.



# Lines 

type Index []Line

this basically stores copies of lines, which might be vaguely better lookup wise ( the strings are still afar ) 

it could also be *Line for slightly better memory size ( 1ptr vs 3 )

does this start to involve the garbage collector tho for another new object.

int would be the same size as Line, int16, or int8 would be smaller ( a byte width extaction to grow? )

with int, wed have to resort on remove -- or copy down. and maybe there'd be metrics for which is faster. ( probably always copy. )

for locality issues you could bring the three strings into a byte stream so they are all touching. but how much that matters on a pc anyways...?

couple of other thoughts:

 * you could store strings in a byte array, merged -- with an index to the second, and the data. you could have your own comparision against length to account for the first string and second string -- then compare raw against the whole thing to get a good sorted result.
 * you could store crcs instead -- using registration to stave off conflicting instance ids -- and then just sort by int. the lexical ordering might be good -- an interface or alternative implementation might provide thediffernce.


# views 

on the calling side, views ( object pointers/arrays of ) will check their tags to see what they are a view of.

a map in the relation will contain a sync flag, cleared whenever a change to the relation occurs. [ map[string]bool ]

initially views will be readonly, then add changes in having the view talk to the relation.

error checking: the view flag will be checked at startup time , we will keep a list of relations to resolve. [ very compiler like, but whatever. ] and check that they all align.

temporary objects cannot can not participate in views because theres no way to clean up their flags. temporaries are easily findable as those objects lacking ids.

possibly sync will include the following states:
    
* unknown, may or may not have information on the object in question.
* object not in relation
* object view needs sync
* object view up-to-date

itd be fine for now to have just true/false: not in relation and synced is just a nil pointer in the view, right? unknown is basically the same as not synced. 

TODO: evaluate if sync would actually work for all three relation types.
