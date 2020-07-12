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
      const sign=Math.sign(this.source.idx-item.idx);
      console.log("dropper changed", item.idx, sign, item.edge);
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
    var ret;
    const at = this.target;
    const from= this.source;
    if (at && from) {
      // the edge display needs a lot more work
      // it has to follow the same rules as the insertion does.
      // const edges= ["em-item--head","em-item--body","em-item--tail"];
      // const sign= Math.sign(from.idx-at.idx); // negative upper
      // const edge= edges[sign+1];
      ret= ((idx === at.idx) || (idx === at.edge)) && {
          "em-drag-highlight": true,
          // [edge]:true,
          "em-drag-mark": true,
      };
    }
    return ret;
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
