class DragHelper {
  constructor() {
     this.currIdx= false;
     this.currEdge= false;
  }
  setIdx(idx, edge=0) {
    if (this.currIdx !== idx || this.currEdge !== edge) {
      console.log("changing focus", idx, edge);
      this.currIdx= idx;
      this.currEdge= edge;
    }
  }
    // FIX: move this to drag context.
  // generate a vue class for an item based on the current highlight settings.
  highlight(idx, selector="em-item") {
    let ret= false;
    // console.log("highlight", idx, d.currIdx);
    if (idx === this.currIdx) {
      const curr= idx< 0? -1: this.currEdge;
      const edge= curr< 0 ? "--head":
                  curr> 0 ? "--tail":
                  "--body";
      ret= {
        [selector+edge]: true,
        "em-drag-mark": true,
      };
    }
    return ret;
  }
  static findEl(el, nothrow, name="dragIdx") {
    for (; el; el= el.parentElement) {
      if (el.dataset[name] !== undefined) {
        break;
      }
    }
    if (!el && nothrow === undefined) {
      throw new Error("missing drag index");
    }
    return el;
  }
  static findIdx(el, defaultVal, name="dragIdx") {
    let ret= defaultVal;
    const idxEl= DragHelper.findEl(el, defaultVal, name);
    const val= idxEl && idxEl.dataset[name];
    if (val !== undefined) {
      ret= parseInt(val);
    }
    return ret;
  }
  static setDragData(el, dt, data, imgClasses= ["em-drag-image", "em-drag-mark"]) {
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
