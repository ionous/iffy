
// implementation specific to em-table.
class DragGroup {
  constructor(list) {
    this.list= list;
    this.finder= false;
  }
  bind(el) {
    this.finder= el? new TargetFinder(el): false;
  }
  dragOver(start,target) {
    var res;
    const over= this.finder.get(target);
    if (over) {
      // dont allow parents to be dropped into their children.
      // this is lair specific; we would need to check "is parent" more generically.
      let overStart;
      if (start.group === this) {
          overStart= (over.idx === start.idx) ||
                    (this.inline && (over.idx > start.idx));
      } else {
        // bad cases: a, b, c, d
        // 1. same (inline) group and idx is same (or larger)
        // 2. the item we are over has the parent of the item being moved.
        // FIX: dragging a row ( block source ) into the midst of an item.
        const overItem= this.list.items[over.idx];
        if (start.group.list) {
          const startItem= start.group.list.items[start.idx];
          overStart= overItem && overItem.parent === startItem;
        }
      }
      res= !overStart ? over: false;
    }
    return res;
  }
  dragStart(el, dt) {
    let okay= false;
    const start= this.finder.get(el, true);
    if (start) {
      const tgt= this._getDragImage(start, dt);
      Dropper.setDragData(dt, tgt, this._serializeItem(start));
      okay= start;
    }
    return okay;
  }
  drop(from, el) {
    // not sure why, but dt.dropEffect is often 'none' here on chrome;
    // ( even though drag end will be copy )
    const drop= this.finder.get(el);
    if (drop) {
      this.list.transferTo(drop.idx, from.group.list, from.idx);
    }
  }
  // fix: this should be more ....
  _serializeItem(start) {
    const item= this.list.items[start.idx];
    return {
      'text/plain': item.text,
    };
  }
  _getDragImage(start, dt) {
    let tgt= start.el;
    // create a temporary set of elements for an image
    // the blur drag start style is left to the .highlight
    if (this.list.inline) {
      tgt = document.createElement("span");
      let sib= start.el;
      while (1) {
        const add = sib.cloneNode(true);
        tgt.appendChild(add);
        sib= sib.nextSibling;
        if (!sib || TargetFinder.getData(sib, "dragIdx") === undefined) {
          break;
        }
      }
    }
    return tgt;
  }
};
