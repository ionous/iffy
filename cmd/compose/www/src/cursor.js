class Cursor  {
  // points to a potential node from a parent node.
  constructor(parent, token, index=-1) {
    this.parent= parent; // Node
    this.token= token;   // arg name; null for kid in swap/slot. swap choice not recorded.
    this.index= index;   // index if targeting the element of an array
  }
  toJSON() {
    const { token, parent, target } = this;
    return {
      parent: parent ? parent.id: null,
      token: token,
      item: target ? target.id: null,
    }
  }
  get target() {
    const { parent, token, index }= this;
    const el= token? parent.kids[token]: parent.kid;
    return (index>=0) ? el[index]: el;
  }
  get param() {
    let ret= null;
    const { token } = this;
    if (token) {
      const spec= this.parent.itemType.with;
      ret= spec.params[token];
    }
    return ret;
  }
  isPlainText() {
    const { token } = this;
    return token && token.startsWith('$');
  }
  isOptional() {
    const { param } = this;
    return param && param.optional;
  }
  isRepeatable() {
    const { param } = this;
    return param && param.repeats;
  }
  isDeletable() {
    let ret= false;
    // a null token indicates a parent swap or slot
    // and the node there is always deletable.
    const { token }= this;
    if (!token) {
      ret= true;
    } else {
      // otherwise we have a parent run
      // ( there's never a parent prim )
      const { param }= this;
      if (param.optional) {
        ret = true
      } else if (param.repeats) {
        const ar= this.parent.kids[token];
        if (ar) {
          ret = ar.length > 1;
        }
      }
    }
    return ret;
  }
  // return a cursor one element to the left or right of current
  step(side, plainText=false) {
    let ret= null;
    const { index } = this;
    if (index < 0) { // not an array
      ret= this.stepSibling(side, plainText);
    } else {
      const { parent, token } = this;
      const ar= parent.kids[token];
      const next= index+side;
      if (!ar || (next<0) || (next >= ar.length)) {
        ret= this.stepSibling(side, plainText); // an array, but out of bounds.
      } else {
        // otherwise move one index
        ret= new Cursor( parent, token, next );
      }
    }
    return ret;
  }
  // return a cursor that's one field left or right of current
  stepSibling(side, plainText) {
    let ret= null;
    const { token } = this;
    if (token) {
      const { parent } = this;
      const { kids } = parent;
      const { tokens } = parent.itemType.with;
      const tokenIndex= tokens.indexOf(token);
      for (let n= tokenIndex+side; n>=0 && n<tokens.length; n+= side) {
        const t= tokens[n];
        if (plainText || t.startsWith('$')) {
          var i;
          const el= kids[t];
          if (!Array.isArray(el)) {
            i= -1;
          } else if (side<0) { // B[ 1, 2, 3 ], A. left of A is B's tail
            i= el.length-1;
          } else {
            i= 0;   // A, B[ 1, 2, 3 ]. right of A is B's head.
          }
          ret= new Cursor(parent, t, i);
          break;
        }
      }
    }
    return ret;
  }
  static At( node ) {
    let ret= null;
    const { parent }= node;
    if (parent) {
      const swapOrSlot= 'kid' in parent;
      if (!swapOrSlot) {
        // if the parent doesnt have "kid", then it must have "kids"
        // and, node is going to be in there somewhere...
        ret= Cursor.FromRun(parent, node);
      } else if (parent.kid === node) {
        ret= new Cursor(parent, null);
      }
    }
    return ret;
  }
  // returns new Cursor ( or nothing )
  static FromRun( parent, node ) {
    let ret= null; // class Cursor
    if (parent) {
      const { kids }= parent;
      if (kids) {
        for (const t in kids) {
          const v= kids[t]; // v is a node, or nodes.
          if (v === node) {
            ret= new Cursor(parent, t);
            break;
          } else if (Array.isArray(v)) {
            const i= v.indexOf(node);
            if (i >= 0) {
              ret= new Cursor(parent, t, i);
              break;
            }
          }
        }
      }
    }
    return ret;
  }
}
