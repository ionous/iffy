
// note; target is from TargetFinder
// it includes: el, idx, edge
class DraggableNode extends Draggable {
  constructor(list, target, width = 1) {
    super();
    this.list= list;
    this.target= target;
    this.width= width;
  }
  getNode(ofs=0) {
    console.assert(ofs>=0 && ofs< this.width, "offset out of range");
    return this.list.items[this.target.idx+ofs];
  }
  getType() {
    const node= this.getNode();
    return node.type;
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
      'text/plain': item.type,
    };
  }
  getDragImage() {
    return this.target.el.cloneNode(true);
  }
}

// one node, that may have many children.
// ex. an action/execute statement, a paragraph block, a pattern rule.
class DraggableBlock extends DraggableNode {
}

// a series of nodes all in a row.
// ex. several story "phrases".
class DraggableSiblings extends DraggableNode {
  constructor(list, target) {
    super(list, target, list.length-target.idx);
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
class NodeTableEvents  {
  constructor(list, copier) {
    this.list= list; // always a NodeList
    this.finder= false;
    this.copier= copier; // FIX! probably should be handled by the dropper and use drag effect
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
    let ret;
    const target= this.finder.findIdx(targetEl);
    if (target && this.list.acceptsType(from.getType())) {
      // fix? we return a "draggable" for use by .highlight and .hovering checks
      ret= this.newDraggable(target)
    }
    return ret;
  }
  dragDrop(from, targetEl) {
    const copying= this.copier.active;
    this.copier.active= false; // ugh.
    // not sure why, but dt.dropEffect is often 'none' here on chrome;
    // ( even though drag end will be copy )
    const target= this.finder.findIdx(targetEl);
    if (target) {
      const { list } = this;

      if (from instanceof DraggableCommand) {
        list.insertAt(target.idx, from.type); // add draggable command
      }
      else if (from instanceof DraggableNode) {
        if (!copying) {
          list.transferTo(target.idx, from.list, from.target.idx, from.width);
        } else {
          let newList=[];
          for (let at=from.target.idx; newList.length< from.width; ++at) {
            const src= from.list.at(at).serialize(false);
            const noids= src.replace(/"id":"[^"]*"[,]?/g, '');
            const newData= JSON.parse(noids);
            const newNodes= list.nodes.newFromItem(null, newData);
            newList.push(newNodes);
          }
          list.spliceInto(target.idx, ...newList);
        }
      }
    }
  }
}
