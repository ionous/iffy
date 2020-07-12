class DragGroup {
  constructor(name, dropper, {serializeItem, addItem, removeItem}) {
    this.name= name;
    this.dropper= dropper;
    this.serializeItem= serializeItem;
    this.addItem= addItem;
    this.removeItem= removeItem;
    this.handler= new DragHandler(this, dropper);
  }
   // generate a vue class for an item based on the current highlight settings.
  highlight(idx) {
    var ret;
    const {target:at, source:from} = this.dropper;
    if (at && from && at.group.name===this.name) {
      // the edge display needs a lot more work
      // it has to follow the same rules as the insertion does.
      // const edges= ["em-item--head","em-item--body","em-item--tail"];
      // const sign= Math.sign(from.idx-at.idx); // negative upper
      // const edge= edges[sign+1];
      ret= ((idx === at.idx) || (idx === at.edge)) && {
          "em-drag-highlight": true,
          // [edge]:true,
          "em-drag-mark": true,
      };
    }
    return ret;
  }
};

