// event handler
class DragHandler {
  constructor(group, dropper) {
    this.group= group;
    this.dropper= dropper;
    this.finder= false;
    this.listeners= false;
    //
    this.inline= false;
  }
  listen(el, inline=false) {
    if (this.listeners) {
      throw new Error("still listening");
    }
    this.inline= inline;
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
      let tgt= start.el;
      // create a temporary set of elements for an image
      // the blur drag source style is left to the .highlight
      if (this.inline) {
        tgt = document.createElement("span");
        let sib= start.el;
        while (1) {
          const add = sib.cloneNode(true);
          tgt.appendChild(add);
          sib= sib.nextSibling;
          if (!sib || TargetFinder.getData(sib, "dragIdx") === undefined) {
            break;
          }
        }
      }
      this.dropper.setSource(this.group, start);
      Dropper.setDragData(dt, tgt, this.group.serializeItem(start.idx));
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
        const {idx:dropIdx}= drop;
        const {idx:dragIdx, group:dragGroup} = this.dropper.source;
        const newGroup= this.group!== dragGroup;
        // add and remove can ( sometimes ) cause dragend not to fire.
        // fix? while moving items is quick and easy
        // technically, we should create new items here by serialization --
        // and wait to remove items in drag end.
        //
        let width= 1;
        if (dragGroup.handler.inline) {
          width= Number.MAX_VALUE;
        }
        const rub= dragGroup.removeItem(dragIdx, dropIdx, width, newGroup);
        this.group.addItem(dragIdx, dropIdx, rub, newGroup);
        // clear b/c we dont always get dragEnd.
        this.dropper.reset(true);
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
    this.dropper.updateTarget(this);
    evt.stopPropagation();
  }
  // this gets triggered on the drag start element, and bubbles here.
  // if the drag start element has been removed, this will never fire.
  onDragEnd(evt) {
    this.log(evt);
    //
    this.dropper.reset(true);
    evt.stopPropagation();
    evt.preventDefault();
  }
  onDragLeave(evt) {
    this.log(evt);
    this.dropper.leaving= this.group;
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
      this.dropper.setTarget(this.group, over);

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
      "idx:", this.group.name, tgt.idx, "edge:", tgt.edge,
      "fx:", fx);
  }
};
