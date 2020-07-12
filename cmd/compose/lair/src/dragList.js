class DragList {
  constructor(items, makeBlank) {
    this.items= items;
    this.makeBlank= makeBlank;
  }
  // here for testing, illustration
  // can't really be in a single function b/c it can happen across groups.
  adjust(src, dst, width=1) {
    const rub= this.removeFrom(src, dst, width);
    this.addTo(src, dst, rub);
  }
  removeFrom(src, dst, width=1) {
    var rub;
    const d= src-dst;
    if (d >0) {
      rub= this._remove(src,dst, 1, width);
    } else if (d<0) {
      rub= this._remove(src,dst,-1, width);
    }
    return rub;
  }
  addTo(src, dst, rub) {
    const d= src-dst;
    if (d >0) {
      this._add(src,dst, 1, rub);
    } else if (d<0) {
      this._add(src,dst,-1, rub);
    }
  }
  _remove(src,dst,sign, width) {
    return (src-dst === sign) ?
          [ this.makeBlank() ] :
          this.items.splice(src, width);
  }
  _add(src,dst,sign, rub) {
    // if (src-dst!==sign) {
    //   dst= dst+sign;
    // } else if (sign>0) {
    //   dst= src+sign;
    // } else {
    //   dst= dst+sign;
    // }
    if (src-dst===1) {
      dst= src+1;
    } else {
      dst= dst+sign;
    }
    this.items.splice(dst,0,...rub);
  }
};


(function() {
  console.log("testing drag list");
  function test(og, s, e, expect) {
    const items= og.split('');
    const dl= new DragList(items, ()=>"_");
    dl.adjust(s,e);
    const res= items.join("");
    if (expect !== res) {
     console.log("Error, moving", og[s], "want:", expect, "have:", res);
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
