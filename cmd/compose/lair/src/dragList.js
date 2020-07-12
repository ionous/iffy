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
    const diff= src-dst;
    if (diff) {
      // movement within 1 spot adds a space to the opposite side
      if ((diff === 1) || (diff === -1)) {
        const blank= this.makeBlank();
        rub= [ blank ];
      } else{
        rub= this.items.splice(src,width);
      }
    }
    return rub;
  }
  addTo(src, dst, rub) {
    const diff= src-dst;
    if (diff) {
      // movement within 1 spot adds a space to the opposite side
      if (diff === 1) {
        dst= src+1;
      } else if (diff === -1) {
        dst= src;
      } else if (diff <0) {
        // dst>src; then removing src shifted down dst by one slot
        --dst;
      }
      this.items.splice(Math.max(0,dst),0,...rub);
    }
  }
};


(function() {
  function test(list, s, e, expect) {
    const items= list.split(',');
    const dl= new DragList(items, ()=>"_");
    dl.adjust(s,e);
    const res= items.join(",")
    if (expect !== res) {
     console.log("Error, want:", expect, "have:", res);
    }
  };
  const abc= "a,b,c";
  // console.log("moving a");
  test(abc, 0, -1, "a,_,b,c");
  test(abc, 0, 0, abc);
  test(abc, 0, 1, "_,a,b,c");
  test(abc, 0, 2, "b,a,c");
  test(abc, 0, 3, "b,c,a");

  // console.log("moving b");
  test(abc, 1, -1, "b,a,c"); // a,b,c
  test(abc, 1, 0, "a,b,_,c");
  test(abc, 1, 1, abc);
  test(abc, 1, 2, "a,_,b,c");
  test(abc, 1, 3, "a,c,b");

  // console.log("moving c");
  test(abc, 2, -1, "c,a,b");
  test(abc, 2, 0, "c,a,b");
  test(abc, 2, 1, "a,b,c,_");
  test(abc, 2, 2, abc);
  test(abc, 2, 3, "a,b,_,c");

})();
