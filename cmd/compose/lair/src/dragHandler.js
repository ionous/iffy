

// event handler
class DragHandler {
  constructor(group, dropper, {serializeItem, addItem, removeItem}) {
    this.group= group;
    this.dropper= dropper;
    this.serializeItem= serializeItem;
    this.addItem= addItem;
    this.removeItem= removeItem;
    this.finder= false;
    this.listeners= false;
  }
  listen(el) {
    if (this.listeners) {
      throw new Error("still listening");
    }
    this.finder= new TargetFinder(el);
    this.listeners= new EventGroup(el, this, {
      // the user starts dragging an item;
      // triggered on the item in question.
      dragstart: "onDragStart",
      // a dragged item has moved;
      // triggered on the drag start item, in the original group.
      drag:      "onDragUpdate",
      // a dragged item enters a valid drop target;
      // triggered on the target of the drop.
      dragenter: "onDragEnterOver",
      // triggered on the current drop target every few hundred milliseconds.
      dragover: "onDragEnterOver",
      // - dragexit: an element is no longer the drag operation's immediate selection ?
      //
      // a dragged item leaves a valid drop target;
      // triggered on the target of the drop.
      dragleave: "onDragLeave",
      // an item is dropped onto a valid drop target;
      // triggered on the target of the drop
      drop:      "onDrop",
      // a drag operation has finished, successfully or not.
      // triggered on the drag start item
      dragend:   "onDragEnd",
    });
  }
  silence() {
    this.listeners= this.listeners.silence();
    this.finder= false;
  }
  // the event targets the em-gutter (draggable=true)
  // the user is attempting to drag,
  // and bubbles up here to the item.
  onDragStart(evt) {
    const dt= evt.dataTransfer;
    const start= this.finder.get(evt.target, true);
    if (start) {
      this.dropper.setSource(this, start);
      DragHelper.setDragData(dt, start.el, this.serializeItem(start.idx));
      evt.stopPropagation();
      this.log(evt)
    }
  }
  // caught by bubbling
  // we've already removed items from
  onDrop(evt) {
    this.log(evt);
    const dt= evt.dataTransfer;
    /*if (dt === "copy")*/ {
      const drop= this.finder.get(evt.target);
      if (drop) {
        const {idx:dragIdx, handler:dragHandler} = this.dropper.source;
        const {idx:dropIdx, handler:dropHandler}= drop;
        const newGroup= dropHandler!== dragHandler;
        //
        const rub= dragHandler.removeItem(dragIdx, dropIdx, 1, newGroup);
        this.addItem(dragIdx, dropIdx, rub, newGroup);
      }
    }
    //
    evt.stopPropagation();
    evt.preventDefault();
  }
  // the drag event targets the same element as drag start
  // and happens periodically as you move the cursor around.
  // ie. it only happens in the originating group
  onDragUpdate(evt) {
    // this.log(evt);
    this.dropper.updateTarget(this);
    evt.stopPropagation();
  }
  onDragEnd(evt) {
    this.log(evt);
    //
    this.finder.reset(true);
    this.dropper.reset(true);
    //
    evt.stopPropagation();
    evt.preventDefault();
  }
  onDragLeave(evt) {
    this.log(evt);
    this.dropper.leaving= this;
    evt.stopPropagation();
    evt.preventDefault();
  }
  // called on dragenter, dragover
  onDragEnterOver(evt) {
    const over= this.finder.get(evt.target);
    if (over) {
      const dt= evt.dataTransfer;
      dt.dropEffect= "copy";
      this.log(evt);
      //
      this.dropper.setTarget(this, over);

      evt.stopPropagation();
      evt.preventDefault();
    }
  }
  log(evt) {
    // return;
    const el= evt.target;
    const dt= evt.dataTransfer;
    const tgt= this.finder.get(el) || {idx:"xxx", edge:false};
    const fx= (dt&&dt.dropEffect)||"???";
    console.log(evt.type, "@", el.nodeName,
      "idx:", this.group, tgt.idx, "edge:", tgt.edge,
      "fx:", fx);
  }
};
