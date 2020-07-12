class DragHelper {
  constructor() {
    this.reset();
  }
  reset(log) {
    if (log) {
      console.log("dropper reset");
    }
    this.source= false;
    this.target= false;
    this.leaving= false;
  }
  setSource(handler, item) {
    const src= DragHelper.setGroup(handler, item);
    this.source= src;
    this.target= src;
    console.log("dropper set source");
  }
  setTarget(handler, item) {
   if (this.target.handler!== handler ||
        this.target.idx !== item.idx ||
        this.target.edge !== item.edge)
    {
      const sign=Math.sign(this.source.idx-item.idx);
      console.log("dropper changed", handler.group, item.idx, sign, item.edge);
      this.target= DragHelper.setGroup(handler, item);
    }

    this.leaving= false;
  }
  updateTarget() {
    if (this.leaving === this.target.handler) {
      console.log("dropper target cleared");
      this.target= false;
    }
  }
  // generate a vue class for an item based on the current highlight settings.
  highlight(group, idx) {
    var ret;
    const at = this.target;
    const from= this.source;
    if (at && from && at.handler.group===group) {
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
  // add the handler to the passed parameters
  static setGroup(handler, {el, idx, edge}) {
    return { handler, el, idx, edge };
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
