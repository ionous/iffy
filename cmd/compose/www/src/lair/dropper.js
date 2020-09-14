// a single root Dropper exists per app.
// vue components can watch it for changes.
// DragHandler instances write to it.
class Dropper {
  constructor(hack) {
    this.parent= hack; // for shift
    this.reset();
  }
  reset(log) {
    if (log && (this.start || this.target || this._leaving)) {
      console.log("dropper reset");
    }
    this.start= false;   // a Draggable
    this.target= false;  // a Draggable
    this._leaving= false; // target
    this.dragging= false;
  }
  // start is a Draggable
  setStart(start, dt, imgClasses= ["em-drag-image"]) {
    console.assert(start instanceof Draggable);

    // set drag visuals
    if (start.getDragImage) {
      const el= start.getDragImage();
      el.classList.add(...imgClasses);
      const add= !el.parentElement;
      document.body.append(el); // needed to display
      dt.setDragImage(el,-10,-10); // fix? maybe should be click relative?
      setTimeout(() => {
          el.remove();
      });
    }

    // set drag content
    if (start.getDragData) {
      const dragData= start.getDragData();
      for (const k in dragData) {
        const v= dragData[k] || "<missing data>";
        dt.setData(k, v);
      }
    }
    //
    dt.effectAllowed= 'all';
    console.log("dropper set start", start);

    // note: can't modify styles until a frame after drag starts:
    // otherwise, chrome will (sometimes) cancel the drag.
    const delayed= this;
    delayed.dragging= null; // pending
    setTimeout(()=> {
      if (delayed.dragging === null) {
        delayed.parent.shift= false; // HACK to turn off shift display while dragging.
        delayed.dragging= start;
        delayed.start= start;
        delayed.target= start;
      }
    });
  }
  // called from DragHandler.onDragEnterOver ( @dragenter, @dragover )
  // target is a Draggable
  setTarget(target) {
    console.assert(target instanceof Draggable);
    this.target= target;
    this._leaving= false;
  }
  // called from DragHandler.onDragUpdate (@drag)
  updateTarget() {
    if (this._leaving === this.target) {
      console.log("dropper target cleared");
      this.target= false;
    }
  }
  // called from DragHandler.onDragLeave (@leave)
  leaving(target) {
    this._leaving= target;
  }
};
