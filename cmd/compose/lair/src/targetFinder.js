
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
      const res= this.findIdx(target);
      this.lastTgt= target; // cache...
      this.lastRes= res;
    }
    return this.lastRes;
  }
  // search upwards from el for the dataset attributes
  findIdx(el) {
    var ret= false;
    for (const topEl= this.topEl; el !== topEl; el= el.parentElement) {
      const idx= TargetFinder.getData(el, "dragIdx");
      if (idx !== undefined) {
        const edge= TargetFinder.getData(el, "dragEdge");
        ret= { el, idx, edge  };
        break;
      }
    }
    return ret;
  }
  // from https://github.com/jquery/jquery/blob/master/src/data.js
  static getData(el, key) {
    const data= el.dataset[key];
    if ( data === "true" ) {
      return true;
    }
    if ( data === "false" ) {
      return false;
    }
    if ( data === "null" ) {
      return null;
    }
    // Only convert to a number if it doesn't change the string
    if ( data === +data + "" ) {
      return +data;
    }
    return data;
  }
};
