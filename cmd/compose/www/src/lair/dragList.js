let _DragList_lastList= 0; // helper to generate unique names

// interface to dragGroup and dropper.
class DragList {
  constructor(items, inline=false) {
    this.name= `list-${++_DragList_lastList}`;
    this.items= items;
    this.inline= !!inline;
  }
  // sent to the target list.
  transferTo(dstIdx, fromGroup, fromIdx) {
    throw new Error("not implemented");
  }
};
