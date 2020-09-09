
// note; target is from TargetFinder
// it includes: el, idx, edge
class DraggableNode extends Draggable {
  constructor(list, target, width = 1) {
    super();
    this.list= list;
    this.target= target;
    this.width= width;
  }
  // is the passed idx (directly) referenced by this draggable
  contains(list, idx) {
    let ret= false;
    if (this.list === list) {
      const beg= this.target.idx;
      const end= beg + this.width;
      ret = (idx >= beg) && (idx < end);
    }
    return ret;
  }
  getDragData() {
    const { list, target } = this;
    const item= list.items[target.idx];
    return {
      'text/plain': item.text,
    };
  }
  getDragImage() {
    return this.target.el.cloneNode(true);
  }
}

// one node, that may have many children.
// ex. an action/execute, a paragraph block, a pattern rule.
class DraggableBlock extends DraggableNode {
  getNode() {
    return this.list.items[this.target.idx];
  }
}

// a series of nodes all in a row.
// ex. several story "phrases".
class DraggableSiblings extends DraggableNode {
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
    return list.inline? new DraggableSiblings(list, target): new DraggableBlock(list, target);
  }
  dragOver(from, targetEl) {
    let ret, okay= false;
    const target= this.finder.findIdx(targetEl);
    if (target) {
      const { list }= this;

      if (from instanceof DraggableCommand) {
        okay= list.acceptsType(from.itemType);
        //
      } else if (from instanceof DraggableBlock) {
        // does the node match or implement our type or slot.
        const node= from.getNode();
        okay= list.acceptsBlock(node.itemType);
      } else if (from instanceof DraggableSiblings) {
        // dont allow parents to be dropped into their children.
        // this is lair specific; we would need to check "is parent" more generically.
        if (from.list === list) {
          okay= !from.contains(list, target.idx);
        } else {
          // bad cases: a, b, c, d
          // 1. same (inline) list and idx is same (or larger)
          // 2. the item we are target has the parent of the item being moved.
          // FIX: dragging a row ( block source ) into the midst of an item.
          const overItem= list.items[target.idx];
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
        let newItem= this.list.nodes.newFromType(from.type);
        const blank= this.list.nodes.newFromType(this.list.type);
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
        this.list.insertAt(target.idx, newItem);
      }
      else if (from instanceof DraggableNode) {
        this.list.transferTo(target.idx, from.list, from.target.idx, from.width);
      }
    }
  }
}
