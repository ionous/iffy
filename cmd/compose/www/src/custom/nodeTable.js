
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

//
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
    this.list= list; // always a NodeList
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
  dragOver(from, targetEl) {
    let ret, okay= false;
    const target= this.finder.findIdx(targetEl);
    if (target) {
      if (from instanceof DraggableCommand) {
        const type= allTypes.all[from.type];
        const spec= type && type.with;
        if (spec && spec.slots) {
          if (this.list.type === "paragraph") {
            if ((from.type=== "paragraph") || (spec.slots.indexOf("story_statement")>=0)) {
              okay= true;
            }
          } else {
            if (spec.slots.indexOf(this.list.type)>=0) {
               okay= true;
            }
          }
        }
      }
      else if (from instanceof DraggableNode) {
        // dont allow parents to be dropped into their children.
        // this is lair specific; we would need to check "is parent" more generically.
        if (from.list === this) {
            const isAtStart= (target.idx === from.target.idx) ||
                      (this.list.inline && (target.idx > from.target.idx));
            if (!isAtStart) {
              okay= true;
            }
        } else {
          // bad cases: a, b, c, d
          // 1. same (inline) list and idx is same (or larger)
          // 2. the item we are target has the parent of the item being moved.
          // FIX: dragging a row ( block source ) into the midst of an item.
          const overItem= this.list.items[target.idx];
          if (from.list) {
            const fromItem= from.list.items[from.idx];
            const overStart= overItem && overItem.parent === fromItem;
            if (!overStart) {
              okay= true;
            }
          }
        }
      }
    }
    return okay && this.newDraggable(target);
  }
  dragDrop(from, targetEl) {
    // not sure why, but dt.dropEffect is often 'none' here on chrome;
    // ( even though drag end will be copy )
    const target= this.finder.findIdx(targetEl);
    if (target) {
      if (from instanceof DraggableCommand) {
        let newItem= this.list.redux.nodes.newFromType(from.type);
        const blank= this.list.makeBlank();
        if (this.list.type !== "paragraph") {
          blank.kid= newItem;
          newItem.parent= blank;
          newItem= blank;
        } else {
          const parent= blank.kids["$STORY_STATEMENT"][0];
          newItem.parent= parent;
          parent.kid= newItem;
          newItem= blank;
        }
        this.list.addBlank(target.idx, newItem);
      }
      else if (from instanceof DraggableNode) {
        this.list.transferTo(target.idx, from.list, from.target.idx, from.width);
      }
    }
  }
}
