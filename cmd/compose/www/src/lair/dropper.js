// a single root Dropper exists per app.
// vue components can watch it for changes.
// DragHandler instances write to it.
class Dropper {
  constructor() {
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
  setStart(start, dt, imgClasses= ["em-drag-image", "em-drag-mark"]) {
    console.assert(start instanceof Draggable);
    this.start= start;
    this.target= start;
    this._leaving= false;

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
        const v= dragData[k];
        dt.setData(k, v);
      }
    }
    //
    dt.effectAllowed= 'all';
    console.log("dropper set start", start);

    const delayed= this;
    delayed.dragging= null; // pending
    setTimeout(()=> {
      if (delayed.dragging === null) {
        delayed.dragging= start;
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
