class DragHelper {
  constructor() {
    this.reset();
  }
  setSource(group, item) {
    const src= DragHelper._group(group, item);
    this.source= src;
    this.target= src;
    console.log("dropper set source");
  }
  static _group(group, {el, idx, edge}) {
    return { group, el, idx, edge };
  }
  setTarget(group, item) {
    if (!group) {
      if (this.target) {
        console.log("dropper cleared");
        this.target= false;
      }
    } else if (this.target.group!== group ||
        this.target.idx !== item.idx ||
        this.target.edge !== item.edge) {
      console.log("dropper changed", item.idx, item.edge);
      this.target= DragHelper._group(group, item);
    }
  }
  reset(log) {
    if (log) {
      console.log("dropper reset");
    }
    this.source= false;
    this.target= false;
  }
  // generate a vue class for an item based on the current highlight settings.
  highlight(idx) {
    const at = this.target;
    // we cheat slightly and use == to hide the differences between idx strings and ints
    return at && (idx == at.idx) && {
        "em-drag-highlight": (!at.edge),
        "em-drag-border": (at.edge) || (at.idx != this.source.idx),
        "em-drag-mark": true,
    };
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
