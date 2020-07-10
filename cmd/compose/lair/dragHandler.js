// event handler
class DragHandler {
  constructor(focus, {serializeItem, removeItem}) {
    this.focus= focus;
    this.serializeItem= serializeItem;
    this.removeItem= removeItem;
    this.startIdx= false;
    this.lastEl= false;
    this.lastIdx= false;
    this.lastEdge= false;
    this.leaving= false;
  }
  // the event targets the em-gutter (draggable=true)
  // the user is attempting to drag,
  // and bubbles up here to the item.
  onDragStart(evt) {
    const dt= evt.dataTransfer;
    const el= DragHelper.findEl(evt.target);
    const idx= DragHelper.findIdx(el);
    this.startIdx= idx;
    DragHelper.setDragData(el, dt, this.serializeItem(idx));
    evt.stopPropagation();
    this.log(evt);
  }
  // the drag event targets the same element as drag start
  // and happens periodically as you move the cursor around.
  onDragUpdate(evt) {
    // this.log(evt);
    if (this.leaving) {
      this.focus.setIdx(false,false);
    }
  }
  onDragEnd(evt) {
    this.log(evt);
    const dt= evt.dataTransfer;
    if ((dt === "move") && (this.lastIdx!== false)) {
      this.removeItem(this.lastIdx);
    }
    this.lastEl= false;
    this.lastIdx= false;
    this.focus.setIdx(false,false);
    //
    evt.stopPropagation();
    evt.preventDefault();
  }
  onDragLeave(evt) {
    this.log(evt);
    this.leaving= true;
    evt.stopPropagation();
  }
  // called on dragenter, dragover
  onDragItem(evt) {
    var idx, edge;
    if (this.lastEl === evt.target) {
      idx= this.lastIdx;
      edge= this.lastEdge;
    } else {
      idx= DragHelper.findIdx(evt.target, false);
      if (idx<0) {
        idx= DragHelper.findIdx(evt.target, idx, "dragPrev");
        edge= 1;
      } else if (idx!==this.startIdx) {
        edge=-1;
      } else {
        edge= 0;
      }
    }
    //
    if (idx !== false) {
      this.focus.setIdx(idx, edge);
      //
      this.leaving= false;
      this.lastIdx= idx;
      this.lastEdge= edge;
      this.lastEl= evt.target;
      //
      const dt= evt.dataTransfer;
      dt.dropEffect= "copy";
    }
    evt.stopPropagation();
    evt.preventDefault();
    this.log(evt);
  }
  log(evt) {
    // return;
    const el= evt.target;
    const dt= evt.dataTransfer;
    const idx= DragHelper.findIdx(el, "xxx");
    console.log(evt.type, idx, el.nodeName,
      "fx:", (dt&&dt.dropEffect)||"???");
  }
};

