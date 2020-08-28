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
  }
  // start is a Draggable
  setStart(start, dt) {
    console.assert(start instanceof Draggable);
    this.start= start;
    this.target= start;
    this._leaving= false;

    // set drag visuals
    if (start.getDragImage) {
      const dragImage= start.getDragImage();
      Dropper.setDragData(dt, dragImage);
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
    console.log("dropper set start", start);
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
  static setDragData(dt, el, imgClasses= ["em-drag-image", "em-drag-mark"]) {
    const existed= !!el.parentElement;
    if (!existed) {
        document.body.append(el);
    }
    // set fx
    dt.effectAllowed= 'all';
    // set the drag image
    el.classList.add(...imgClasses);
    dt.setDragImage(el,-10,-10); // fix? maybe should be click relative?
    setTimeout(()=>{
      if (existed) {
        el.classList.remove(...imgClasses);
      } else {
        el.remove();
      }
    });
  }
};
