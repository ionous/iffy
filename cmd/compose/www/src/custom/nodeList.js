
class NodeList {
  // token should target an array of the passed type
  constructor( nodes, node, token, type ) {
    this.nodes= nodes;
    this.node= node;
    this.token= token;
    this.type= type
    this.items= node.getKid(token);
    this.inline= false;
  }
  get length() {
    return this.items.length;
  }
  // returns number of elements added
  addTo(at, exe) {
    const { node, items } = this;
    exe.parent= node;
    items.splice(at, 0, exe);

  }
  // returns the element or elements removed
  // when we drag, we re/move a single execute ( a line ) at once.
  // returns a single statement
  removeFrom(at) {
    var one;
    const rub= this.items.splice(at, 1);
    if (rub.length) {
      one= rub[0];
      one.parent= null;
     }
     return one;
  }
  // at:index, from:Draggable
  transferTo(toIdx, fromList, fromIdx, width=1) {
    const toList= this;
    //
    if (toList === fromList) {
      this.move(toIdx, fromIdx, width);
    } else {
      Redux.Run({
        added: 0, // inelegant to say the least.
        apply() {
          const paraEls= fromList.removeFrom(fromIdx);
          this.added= toList.addTo( toIdx, paraEls );
        },
        revoke() {
          const paraEls= toList.removeFrom( toIdx, this.added );
          fromList.addTo( fromIdx, paraEls );
        },
      });
    }
  }
  insertAt(at, newItem) {
    const { node, items } = this;
    if (at<0) {
      at= items.length;
    }
    Redux.Run({
      apply() {
        newItem.parent= node;
        items.splice(at,0,newItem);
      },
      revoke() {
        newItem.parent= null;
        items.splice(at,1);
      },
    });
    return newItem;
  }
  // move items within this same list
  move(src, dst, width, nothrow) {
    const { items } = this;
    if ((!width) || (width<0)) {
      const e= new Error("invalid width");
      if (nothrow) { console.error(e); return e; }
      throw e;
    }
    if ((dst > src) && (dst < src+width)) {
      const e= new Error("invalid dest");
      if (nothrow) { return e; }
      throw e;
    }
    if (src+width> items.length) {
      width= items.length-src;
    }
    if (dst > src) {
      dst -= width;
    }
    Redux.Run({
      apply() {
        const rub= items.splice(src, width);
        items.splice(dst, 0, ...rub);
      },
      revoke() {
        const rev= items.splice(dst, width);
        items.splice(src, 0, ...rev);
      }
    });
  }
}
