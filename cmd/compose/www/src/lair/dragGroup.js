
// implementation specific to em-table.
class DragGroup {
  constructor(list, dropper) {
    this.list= list;
    this.dropper= dropper;
  }
  dragOver(over,dt) {
    const mylist= this.list;
    const start= this.dropper.start
    if (start) {
      // dont allow parents to be dropped into their children.
      // this is lair specific; we would need to check "is parent" more generically.
      var overStart;
      if (start.list === mylist) {
          overStart= (over.idx === start.idx) ||
                    (mylist.inline && (over.idx > start.idx));
      } else {
        // bad cases: a, b, c, d
        // 1. same (inline) group and idx is same (or larger)
        // 2. the item we are over has the parent of the item being moved.
        // FIX: dragging a row ( block source ) into the midst of an item.
        const overItem= this.list.items[over.idx];
        const startItem= start.list.items[start.idx];
        overStart= overItem && overItem.parent=== startItem;
      }
      if (!overStart) {
        this.dropper.setTarget(mylist, over);
        dt.dropEffect= "copy";
      }
    }
  }
  drag() {
    this.dropper.updateTarget(this.list);
  }
  dragLeave(leave, dt) {
    this.dropper.leaving= this.list;
  }
  dragEnd() {
    this.dropper.reset(true);
  }
  dragStart(start, dt) {
    this.dropper.setSource(this.list, start);
    const tgt= this._getDragImage(start, dt);
    Dropper.setDragData(dt, tgt, this._serializeItem(start));
  }
  drop(drop, dt) {
    const start= this.dropper.start;
    if (start) {
      this.list.transferTo(drop.idx, start.list, start.idx);
    }
    // clear b/c we dont always get dragEnd.
    this.dropper.reset(true);
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
