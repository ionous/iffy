// event handler
class DragHandler {
  // dropper is a global instance of Dropper;
  // the event sink is an object with the following optional methods:
  //   - bind(container)
  //   - dragStart(el, dt) -> returns Draggable
  //   - dragOver(start, el, evt)
  //   - dragDrop(start, el, evt)
  //   - dragUpdate()
  //   - dragLeave()
  //   - dragEnd()
  constructor(dropper, sink) {
    this.dropper= dropper;
    this.sink= sink;
    this.listeners= false;
  }
  listen(el) {
    const { listeners, sink } = this;
    if (listeners) {
      throw new Error("still listening");
    }
    this.listeners= new EventGroup(el, this, {
      // the user starts dragging an item;
      // triggered on the item in question.
      dragstart: "onDragStart",
      // a dragged item has moved;
      // triggered on the drag start item, in the original sink.
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
    if (sink.bind) {
      sink.bind(el);
    }
    return this; // for chaining
  }
  silence() {
    const { listeners, sink } = this;
    if (sink.bind) {
      sink.bind(null);
    }
    listeners.silence();
    this.listeners= null;
  }
  // the event targets the em-gutter (draggable=true)
  // the user is attempting to drag,
  // and bubbles up here to the item.
  onDragStart(evt) {
    const { dropper, sink } = this;
    if (sink.dragStart) {
      const start= sink.dragStart(evt.target, evt.dataTransfer);
      if (start) {
        dropper.setStart(start, evt.dataTransfer);
        evt.stopPropagation();
        this.log(evt)
      }
    }
  }
  // caught by bubbling in the sink that receives the item
  onDrop(evt) {
    const { dropper, sink } = this;
    this.log(evt);
    if (dropper.start && sink.dragDrop) {
      // not sure why, but dt.dropEffect is often 'none' here on chrome;
      // ( even though drag end will be copy )
      sink.dragDrop(dropper.start, evt.target, evt.dataTransfer);
    }
    dropper.reset(true); // clear here b/c we dont always get dragEnd.
    evt.stopPropagation();
    evt.preventDefault();
  }
  // the drag event targets the same element as drag start
  // and happens periodically as you move the cursor around.
  // ie. it only happens in the originating sink
  onDragUpdate(evt) {
    // this.log(evt);
    const { dropper, sink } = this;
    if (sink.dragUpdate) {
      sink.dragUpdate(evt);
    }
    dropper.updateTarget();
    evt.stopPropagation();
  }
  // this gets triggered on the drag start element, and bubbles here.
  // if the drag start element has been removed, this will never fire.
  onDragEnd(evt) {
    this.log(evt);
    const { dropper, sink } = this;
    if (sink.dragEnd) {
      sink.dragEnd()
    }
    dropper.reset(true);
    evt.stopPropagation();
    evt.preventDefault();
  }
  onDragLeave(evt) {
    this.log(evt);
    const { dropper, sink, target } = this;
    if (sink.dragLeave) {
      sink.dragLeave();
    }
    dropper.leaving(target);
    evt.stopPropagation();
    evt.preventDefault();
  }
  // called on dragenter, dragover
  onDragEnterOver(evt) {
    this.log(evt);
    const { dropper, sink } = this;
    const { start } = dropper;
    if (sink.dragOver && start)  {
      const { target, dataTransfer:dt } = evt;
      const over= sink.dragOver(start, target, evt);
      if (over !== undefined) {
        if (over) {
          dropper.setTarget(over);
          dt.dropEffect= "copy";
        }
        evt.stopPropagation();
        evt.preventDefault();
      }
    }
  }
  log(evt) {
    return;
    const el= evt.target;
    const dt= evt.dataTransfer;
    // const tgt= this.finder.findIdx(el) || {idx:"xxx", edge:false};
    const fx= (dt&&dt.dropEffect)||"???";
    console.log(evt.type, "@", el.nodeName,
      // "idx:", tgt.target.idx, "edge:", tgt.target.edge,
      "fx:", fx);
  }
};
