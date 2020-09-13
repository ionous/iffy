
class NodeList {
  // token should target an array of the passed type
  constructor( nodes, node, token, type ) {
    const items= node.getKid(token);
    if (!items) {
      throw new Error(`node ${node.id} has empty child ${token}`);
    }
    this.nodes= nodes;
    this.node= node;
    this.token= token;
    this.type= type // typeName
    this.items= items;
    this.inline= false;
  }
  get length() {
    return this.items.length;
  }
  at(i) {
    const { items } = this;
    if (i<0) {
      i= items.length+i;
    }
    console.assert(i>=0 && i<items.length, "index out of range");
    return items[i];
  }
  acceptsType(typeName) {
    const okay= allTypes.areCompatible(typeName, this.type)
    return okay;
  }
  // by default assume that list is a list of slots
  // which means we can insert slots, or slats which implement the slot.
  insertAt(idx, typeName) {
    // create the requested item
    console.assert(allTypes.areCompatible( typeName, this.type ));
    //
    let newItem= this.nodes.newFromType(typeName);
    if (newItem.type !== this.type) {
      const slot= this.nodes.newFromType(this.type);
      slot.putSlot(newItem); // asserts if the item isnt compatible.
      newItem= slot;
    }
    this.spliceInto(idx, newItem);
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
      this.moveTo(toIdx, fromIdx, width);
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
  spliceInto(at, ...newItems) {
    const { node, token } = this;
    Redux.Run({
      apply() {
        node.splices(token, at, 0, ...newItems);
      },
      revoke() {
        node.splices(token, at, newItems.length);
      },
    });
  }
  // move items within this same list
  moveTo(dst, src, width, nothrow) {
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
    // 0, src:1<-remove and dst needs to slide down, 2, 3, dst:4
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
