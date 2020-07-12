
// cache for looking up drop targets by html element.
class TargetFinder {
  // pass the upper most element to stop searching at
  // ex. a container of drop items.
  constructor(el) {
    this.topEl= el;
    this.reset();
  }
  reset(log) {
    if (log) {
      console.log("target reset");
    }
    this.lastTgt= false;
    this.lastRes= false;;
  }
  // return the idx and edge targeted by el.
  get(target, noedges) {
    if (this.lastTgt !== target) {
      let res= this.findIdx(target, false);
      if (!res && !noedges) {
        res= this.findIdx(target, true);
      }
      // cache...
      this.lastTgt= target;
      this.lastRes= res;
    }
    return this.lastRes;
  }
  // search upwards from el for the named dataset attribute
  findIdx(el, edge) {
    var ret= false;
    const key= edge? "dragEdge": "dragIdx";
    for (const topEl= this.topEl; el !== topEl; el= el.parentElement) {
      const idx= el.dataset[key];
      if (idx !== undefined) {
        ret= { el, idx, edge };
        break;
      }
    }
    return ret;
  }
};
