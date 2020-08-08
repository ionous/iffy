let lastList= 0;

class DragList {
  constructor(items, inline=false) {
    this.name= `list-${++lastList}`;
    this.items= items;
    this.inline= !!inline;
  }
  get block() {
    return !this.inline;
  }
  // note: cant always move items in a single moment
  // ( ex. adding/removing across groups)
  move(src, dst, width=1) {
     if (src!== dst) {
      const rub= this.removeFrom(src, dst, width);
      this.addTo(src, dst, rub);
    }
  }
  removeFrom(src, dst, width=1) {
    const sign= Math.sign(src-dst);
    return (src-dst === sign) ?
          [ this.makeBlank() ] :
          this.items.splice(src, width);
  }
  addTo(src, dst, rub) {
    const tgt= this.adjustIndex(src, dst);
    this.items.splice(tgt,0,...rub);
  }
  // when removing and re-adding to the same list,
  // the index of the "add" can change b/c of the removed element.
  adjustIndex(src, dst) {
    if (src-dst===1) {
      dst= src+1;
    } else {
      dst= dst+ Math.sign(src-dst);
    }
    return dst;
  }
  makeBlank() {
    throw new Error("not implemented");
  }
  dropFrom(drop, from) {
    const {idx:dropIdx}= drop;
    const {idx:fromIdx, list:fromList} = from;
    // add and remove can ( sometimes ) cause dragend not to fire.
    // fix? while moving items is quick and easy
    // technically, we should create new items here by serialization --
    // and wait to remove items in drag end.
    //
    let width= 1;
    if (fromList.inline) {
      width= Number.MAX_VALUE;
    }
    if (this === fromList) {
      fromList.move(fromIdx, dropIdx, width);
    } else {
      let rub= fromList.items.splice(fromIdx, width);
      const at= Math.min(Math.max(0,dropIdx+1), this.items.length);

      // moving item(s) from an inline group to a block group?
      const merge= (fromList.inline && this.block);
      if (!merge) {
        // moving a block into an inline list of items.
        if (fromList.block && this.inline) {
          this.items.splice(at, 0,...rub[0].content);
        } else {
          this.items.splice(at, 0,...rub);
        }
      } else {
        let row= new Item();
        row.content= rub.map((x)=> {
          x.parent= row;
          return x;
        });
        this.items.splice(at, 0, row);
      }
    }
  }

};

(function() {
  console.log("testing drag list");
  function test(og, src, dst, expect) {
    class TestList extends DragList {
      makeBlank() { return "_" }
    }
    const items= og.split('');
    const dl= new TestList(items);
    dl.move(src,dst);
    const res= items.join("");
    if (expect !== res) {
     console.log("Error, moving", og[src], "want:", expect, "have:", res);
    }
  };

  // a
  test("abc", 0,-1, "a_bc"); // trailing head
  test("abc", 0, 0, "abc");  // <no change>
  test("abc", 0, 1, "_abc"); // leading b
  test("abc", 0, 2, "bac");  // leading c
  test("abc", 0, 3, "bca");  // leading tail

  // b
  test("abc", 1,-1, "bac");  // trailing head
  test("abc", 1, 0, "ab_c"); // trailing a
  test("abc", 1, 1, "abc");  // <no change>
  test("abc", 1, 2, "a_bc"); // leading c
  test("abc", 1, 3, "acb");  // leading tail

  // c
  test("abc", 2,-1, "cab");  // trailing head
  test("abc", 2, 0, "acb");  // trailing a
  test("abc", 2, 1, "abc_"); // trailing b
  test("abc", 2, 2, "abc");  // <no change>
  test("abc", 2, 3, "ab_c"); // leading tail

})/*()*/;
