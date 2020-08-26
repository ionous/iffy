// event handler
class DragHandler {
  // group is an object with the following optional methods:
  // bind, dragStart, dragOver, drop, drag, dragLeave, dragEnd.
  constructor(dropper, group) {
    this.group= group;
    this.dropper= dropper;
    this.listeners= false;
  }
  listen(el) {
    const { listeners, group } = this;
    if (listeners) {
      throw new Error("still listening");
    }
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
    if (group.bind) {
      group.bind(el);
    }
  }
  silence() {
    const { listeners, group } = this;
    if (group.bind) {
      group.bind(null);
    }
    listeners.silence();
    this.listeners= null;
  }
  // the event targets the em-gutter (draggable=true)
  // the user is attempting to drag,
  // and bubbles up here to the item.
  onDragStart(evt) {
    const { dropper, group } = this;
    if (group.dragStart) {
      const start= group.dragStart(evt.target, evt.dataTransfer);
      if (start!==undefined) {
        if (start) {
          dropper.setSource(group, start);
        }
        evt.stopPropagation();
        this.log(evt)
      }
    }
  }
  // caught by bubbling in the group that receives the item
  onDrop(evt) {
    const { dropper, group } = this;
    this.log(evt);
    if (dropper.start && group.drop) {
      group.drop(dropper.start, evt.target, evt.dataTransfer);
    }
    dropper.reset(true); // clear here b/c we dont always get dragEnd.
    evt.stopPropagation();
    evt.preventDefault();
  }
  // the drag event targets the same element as drag start
  // and happens periodically as you move the cursor around.
  // ie. it only happens in the originating group
  onDragUpdate(evt) {
    // this.log(evt);
    const { dropper, group } = this;
    if (group.drag) {
      group.drag();
    }
    dropper.updateTarget();
    evt.stopPropagation();
  }
  // this gets triggered on the drag start element, and bubbles here.
  // if the drag start element has been removed, this will never fire.
  onDragEnd(evt) {
    this.log(evt);
    const { dropper, group } = this;
    if (group.dragEnd) {
      group.dragEnd()
    }
    dropper.reset(true);
    evt.stopPropagation();
    evt.preventDefault();
  }
  onDragLeave(evt) {
    this.log(evt);
    const { dropper, group } = this;
    if (group.dragLeave) {
      group.dragLeave(evt.dataTransfer);
    }
    dropper.leaving= group;
    evt.stopPropagation();
    evt.preventDefault();
  }
  // called on dragenter, dragover
  onDragEnterOver(evt) {
    this.log(evt);
    const { dropper, group } = this;
    const { start } = dropper;
    if (group.dragOver && start)  {
      const { target, dataTransfer:dt } = evt;
      const res= group.dragOver(start, target, evt);
      if (res !== undefined) {
        if (res) {
          dropper.setTarget(group, res);
          dt.dropEffect= "copy";
        }
        evt.stopPropagation();
        evt.preventDefault();
      }
    }
  }
  log(evt) {
    // return;
    const el= evt.target;
    const dt= evt.dataTransfer;
    // const tgt= this.finder.findIdx(el) || {idx:"xxx", edge:false};
    const fx= (dt&&dt.dropEffect)||"???";
    console.log(evt.type, "@", el.nodeName,
      // "idx:", tgt.idx, "edge:", tgt.edge,
      "fx:", fx);
  }
};
