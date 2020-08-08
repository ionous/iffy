let lastItem=0;

const itemTypes= {
    "lower" : '+', // plus
    "middle": '!',
    "upper" : '\u201C', // left quote
    // '\u2630' hamburder heaven
};
const typeNames= Object.keys(itemTypes);

class Item {
  constructor(parent, content) {
    this.id= `id-${lastItem++}`;
    this.parent= parent;
    this.content= content;
    this._input= false; // cached input,output itemTypes
    this._output= false;
  }
  // the output of a string item is the left side of its string;
  // it connects to the input of an array item ( aka a row )
  get outputType() {
    let t= this._output;
    if (!t) {
      t= this.getType(0);
      this._output=t;
    }
    return t;
  }
  // the desired input of a string item is the right side of its string;
  // its the constraint for the output of its rightward sibling.
  get inputType() {
    let t= this._input;
    if (!t) {
      t= this.getType(-1);
      this._input=t;
    }
    return t;
  }
  getType(at) {
    var ret;
    // determine the first or last string of this item.
    let c= this.content;
    // content is either an array of items, or already a string.
    if (c && Array.isArray(c)) {
      const firstOrLastItems= c.slice(at);
      c= (firstOrLastItems.length)? firstOrLastItems[0].content: "";
    }
    // have a string ( or couldnt find a string )
    if (c) {
      ret= Item.TypeForString(c, at);
    } else {
      const i= Math.floor(Math.random() * typeNames.length);
      ret= typeNames[i];
    }
    return ret;
  }
  // at: 0 is front of string; -1 is end of string.
  static TypeForString(str, at=0) {
    var ret;
    // multiply by 2 to move off of the trailing full-stop.
    const l= str.slice(at*2)[0].toLowerCase();
    if (l < "h") {
      ret = "lower";
    } else if (l < "p") {
      ret = "middle";
    } else {
      ret = "upper";
    }
    return ret;
  }
}
