let lastGroup= 0;
class DragGroup {
  constructor(dropper, {serializeItem, addItem, removeItem}) {
    this.name= `group-${++lastGroup}`;
    this.dropper= dropper;
    this.serializeItem= serializeItem;
    this.addItem= addItem;
    this.removeItem= removeItem;
    this.handler= new DragHandler(this, dropper);
  }
   // generate a vue class for an item based on the current highlight settings.
  highlight(idx) {
    let highlight= false;
    let edge= false;
    const {target:at, source:from} = this.dropper;
    if (at && from && at.group ===this) {
      edge= idx === at.edge;
      highlight=(idx === at.idx) || edge;
    }
    return {
      "em-row": true,
      "em-drag-mark": highlight,
      "em-drag-highlight": highlight,
      "em-drag-head": edge && (at.idx < 0),
      "em-drag-tail": edge && (at.idx > 0),
    }
  }
};

