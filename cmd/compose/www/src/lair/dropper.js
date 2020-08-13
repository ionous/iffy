class Dropper {
  constructor() {
    this.reset();
  }
  reset(log) {
    if (log && (this.start || this.target || this.leaving)) {
      console.log("dropper reset");
    }
    this.start= false;
    this.target= false;
    this.leaving= false;
  }
  get dragging() {
    return !!this.start;
  }
  // found includes { el, idx, edge }
  setSource(group, found) {
    const src= Dropper.record(group, found);
    this.start= src;
    this.target= src;
    console.log("dropper set start", found);
  }
  setTarget(group, found) {
   if (this.target.group!== group ||
        this.target.idx !== found.idx ||
        this.target.edge !== found.edge)
   {
      this.target= Dropper.record(group, found);
    }
    this.leaving= false;
  }
  updateTarget() {
    if (this.leaving === this.target.group) {
      console.log("dropper target cleared");
      this.target= false;
    }
  }
  // add the group to the passed parameter set
  static record(group, {el, idx, edge}) {
    return { group, el, idx, edge };
  }
  static setDragData(dt, el, data, imgClasses= ["em-drag-image", "em-drag-mark"]) {
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
    // set drag content
    for (const k in data) {
      const v= data[k];
      dt.setData(k, v);
    }
  }
};
