class Dropper {
  constructor() {
    this.reset();
  }
  newGroup(name, ops) {
    return new DragGroup(name, this, ops);
  }
  reset(log) {
    if (log) {
      console.log("dropper reset");
    }
    this.source= false;
    this.target= false;
    this.leaving= false;
  }
  setSource(group, item) {
    const src= Dropper.setGroup(group, item);
    this.source= src;
    this.target= src;
    console.log("dropper set source");
  }
  setTarget(group, item) {
   if (this.target.group!== group ||
        this.target.idx !== item.idx ||
        this.target.edge !== item.edge)
   {
      const sign=Math.sign(this.source.idx-item.idx);
      console.log("dropper changed", group.name, item.idx, sign, item.edge);
      this.target= Dropper.setGroup(group, item);
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
  static setGroup(group, {el, idx, edge}) {
    return { group, el, idx, edge };
  }
  static setDragData(dt, el, data, imgClasses= ["em-drag-image", "em-drag-mark"]) {
    // set fx
    dt.effectAllowed= 'all';
    // set the drag image
    el.classList.add(...imgClasses);
    dt.setDragImage(el,10,10); // fix? maybe should be click relative?
    setTimeout(()=>{
      el.classList.remove(...imgClasses);
    });
    // set drag content
    for (const k in data) {
      const v= data[k];
      dt.setData(k, v);
    }
  }
};
