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
};

