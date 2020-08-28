
// note; target is from TargetFinder
// it includes: el, idx, edge
class DraggableNode extends Draggable {
  constructor( list, target, width = 1) {
    super();
    this.list= list;
    this.target= target;
    this.width= width;
  }
  getDragData() {
    const { list, target } = this;
    const item= list.items[target.idx];
    return {
      'text/plain': item.text,
    };
  }
  getDragImage() {
    return this.target.el;
  }
}

class DraggableLine extends DraggableNode {
  constructor(list, target) {
    super(list, target, Number.MAX_VALUE);
  }
  getDragImage() {
    let ret = document.createElement("span");
    let sib= this.target.el;
    while (1) {
      const add = sib.cloneNode(true);
      ret.appendChild(add);
      sib= sib.nextSibling;
      if (!sib || TargetFinder.getData(sib, "dragIdx") === undefined) {
        break;
      }
    }
    return ret;
  }
}

// an event sink implementation specific to em-node-table.
class NodeTable  {
  constructor(list) {
    this.list= list; // always a DragList
    this.finder= false;
  }
  bind(containerEl) {
    this.finder= containerEl? new TargetFinder(containerEl): false;
  }
  dragStart(el, dt) {
    const found= this.finder.findIdx(el, true);
    return found && this.newDraggable(found);
  }
  newDraggable(target) {
    const { list } = this;
    return list.inline? new DraggableLine(list, target): new DraggableNode(list, target);
  }
  dragOver(start, targetEl) {
    var res;
    const target= this.finder.findIdx(targetEl);
    if (target) {
      // dont allow parents to be dropped into their children.
      // this is lair specific; we would need to check "is parent" more generically.
      let overStart;
      if (start.list === this) {
          overStart= (target.idx === start.target.idx) ||
                    (this.list.inline && (target.idx > start.target.idx));
      } else {
        // bad cases: a, b, c, d
        // 1. same (inline) list and idx is same (or larger)
        // 2. the item we are target has the parent of the item being moved.
        // FIX: dragging a row ( block source ) into the midst of an item.
        const overItem= this.list.items[target.idx];
        if (start.list) {
          const startItem= start.list.items[start.idx];
          overStart= overItem && overItem.parent === startItem;
        }
      }
      res= !overStart ? this.newDraggable( target ): false;
    }
    return res;
  }
  dragDrop(start, targetEl) {
    // not sure why, but dt.dropEffect is often 'none' here on chrome;
    // ( even though drag end will be copy )
    const drop= this.finder.findIdx(targetEl);
    if (drop) {
      if (start instanceof DraggableCommand) {
        console.log("add new", from.type);
      }
      else if (start instanceof DraggableNode) {
        this.list.transferTo(drop.idx, start.list, start.target.idx, start.width);
      }
    }
  }
}
