
// look up drop targets by html element.
// relies on 'data-data-idx' and 'data-drag-edge' elements
// grouped together under a common container el.
class TargetFinder {
  // pass the upper most element to stop searching at
  // ex. a container of drop items.
  constructor(el) {
    this.topEl= el;
    this.lastTgt= false;
    this.lastRes= false;
  }
  // return the idx and edge targeted by el. { el, idx, edge }
  findIdx(target, reset) {
    if (reset || (this.lastTgt !== target)) {
      const res= this._findIdx(target);
      this.lastTgt= target; // cache...
      this.lastRes= res;
    }
    return this.lastRes;
  }
  // search upwards from el for the dataset attributes
  _findIdx(el) {
    let ret= false;
    for (const topEl= this.topEl; el !== topEl; el= el.parentElement) {
      const idx= TargetFinder.getData(el, "dragIdx");
      if (idx !== undefined) {
        const edge= TargetFinder.getData(el, "dragEdge");
        ret= { el, idx, edge };
        break;
      }
    }
    return ret;
  }
  // from https://github.com/jquery/jquery/blob/master/src/data.js
  static getData(el, key) {
    var ret;
    const ds= el.dataset;
    if (ds) {
      const data= el.dataset[key];
      switch (data) {
      case "true":
        ret= true;
        break;
      case "false":
        ret= false;
        break;
      case "null":
        ret= null;
        break;
      default:
        // Only convert to a number if it doesn't change the string
        const num= +data;
        if ( data === num + "" ) {
          ret= num;
        } else {
          ret= data;
        }
      }; // end switch
    }
    return ret;
  }
};
