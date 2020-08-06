// add a new item to the left (front) or right (end) of the targeted array.
//  hmmm... regarding target and existing elements

//  weve got multiple routes --
//  first, what are we tring to do here?

//  1. new ghost, which is append to a list
//  2. mutation ( the current test )
//    a. splice in front of ( left ), or behind ( right ) a specific repeated element
//    b. append at end of a repeating element ( left sib ), or at front ( right sib )


//  3. deleteAt, remove a specific targeted element

// -------------------------------------------------------------------------------------

//  1. mutating field == focused field
//    - we arent targeting any specific element of the field
//    - note: cursor doesnt support targeting a field and not an element, should it?

//  2a, mutating field == focused field
//    -  we are tareting a specific element, left means infront, right means behind

//  2b. mutating field != focused field
//    - we are targeting a specific element, left means behind, right means in front


// -------------------------
// how do we deal with optional arrays?




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
  step(side) {
    let ret= null;
    const { index } = this;
    if (index < 0) { // not an array
      ret= this.stepSibling(side);
    } else {
      const { parent, token } = this;
      const ar= parent.kids[token];
      const next= index+side;
      if (!ar || (next<0) || (next >= ar.length)) {
        ret= this.stepSibling(side); // an array, but out of bounds.
      } else {
        // otherwise move one index
        ret= new Cursor( parent, token, next );
      }
    }
    return ret;
  }
  // return a cursor that's one field left or right of current
  stepSibling(side) {
    let ret= null;
    const { token } = this;
    if (token) {
      const { parent } = this;
      const { kids } = parent;
      const { tokens } = parent.itemType.with;
      const tokenIndex= tokens.indexOf(token);
      for (let n= tokenIndex+side; n>=0 && n<tokens.length; n+= side) {
        const t= tokens[n];
        if (t && t[0] === '$') {
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
  // inject the passed node at the current target.
  spliceTarget(node, vm) {
    const { parent, token }= this;
    if (!token) { // no token means swap or slot
      parent.kid= node;
    } else {
      const { index } = this;
      if ((index>=0) && (token in parent.kids)) {
        parent.kids[token].splice(index, 0, node);
      } else {
        const value= (index<0)? node: [node];
        if (vm) {
          vm.set(parent, token, value);
        } else {
          parent[token]= value;
        }
      }
    }
  }
  // note: this will happily delete non-optional elements.
  // vm is an (optional) concession to trigger vue change watchers
  deleteMe(vm) {
    const { parent, token }= this;
    if (!token) { // no token means swap or slot
      parent.kid= null;
    } else {
      const { index } = this;
      if (index > 0) { // a non-zero index means remove a repeater
        parent.kids[token].splice(index, 1);
      } else if (vm) { // otherwise, remove the field entirely
        vm.delete(parent, token);
      } else {
        delete parent[token];
      }
    }
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
