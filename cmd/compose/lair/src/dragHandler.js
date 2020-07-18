// event handler
class DragHandler {
  // dragStart, drop, dragOver, drag, dragLeave, drop, dragEnd
  constructor(callbacks) {
    this.on= callbacks;
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
    const start= this.finder.get(evt.target, true);
    if (start) {
      this.on.dragStart(start, evt.dataTransfer);
      evt.stopPropagation();
      this.log(evt)
    }
  }
  // caught by bubbling in the group that receives the item
  onDrop(evt) {
    this.log(evt);
    const dt= evt.dataTransfer;
    // not sure why, but drop effect is often 'none' here on chrome;
    // ( even though drag end will be copy )
    /*if (dt && dt.dropEffect=== "copy")*/ {
      const drop= this.finder.get(evt.target);
      if (drop) {
        this.on.drop(drop, evt.dataTransfer);
      }
    }
    evt.stopPropagation();
    evt.preventDefault();
  }
  // the drag event targets the same element as drag start
  // and happens periodically as you move the cursor around.
  // ie. it only happens in the originating group
  onDragUpdate(evt) {
    // this.log(evt);
    this.on.drag();
    evt.stopPropagation();
  }
  // this gets triggered on the drag start element, and bubbles here.
  // if the drag start element has been removed, this will never fire.
  onDragEnd(evt) {
    this.log(evt);
    this.on.dragEnd()
    evt.stopPropagation();
    evt.preventDefault();
  }
  onDragLeave(evt) {
    this.log(evt);
    const leave= this.finder.get(evt.target);
    this.on.dragLeave(leave, evt.dataTransfer);
    evt.stopPropagation();
    evt.preventDefault();
  }
  // called on dragenter, dragover
  onDragEnterOver(evt) {
    this.log(evt);
    const over= this.finder.get(evt.target);
    if (over) {
      this.on.dragOver(over, evt.dataTransfer);
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
      "idx:", tgt.idx, "edge:", tgt.edge,
      "fx:", fx);
  }
};
