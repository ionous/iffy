// uses serialization to detect whether some node in a tree includes a type which uses a newline
// we dont want "inline" elements following elements that have blocks/newlines.
class BlockSearch {
  constructor(...blockTypes) {
    const ts= blockTypes.join("|");
    this.regexp= new RegExp(`"type":"(?:${ts})"`);
  }
  hasBlock(node) {
    return this.hasText(node.serialize(false));
  }
  hasText(txt) {
    return txt && this.regexp.test(txt);
  }
}

// base class for the runtime story model.
class Node {
  constructor(parent, itemType, itemId) {
    this.parent = parent;  // Node
    this.id= itemId;
    this.itemType= itemType; // Type
  }
  // return typeName
  get type() {
    return this.itemType && this.itemType.name;
  }
  // return argument info for runs, swap options, string choices.
  getParam(token) {
    const spec= this.itemType.with;
    return spec && spec.params[token];
  }
  toJSON() {
    return {
      id: this.id,
      type: this.type,
    };
  }
  serialize(pretty=true) {
    const args= pretty?[0,2]:[];
    return JSON.stringify(this, ...args);
  }
  unroll(nodes, itemValue) {
    throw new Error("unroll unhandled");
  }
  static LabelFromParam(param, def= "") {
    let ret= def;
    if (param && param.value !== null) {
      // for recapitulation ( where the param value is null and the user can type anything )
      // use the empty string as the label.
      ret= param.label || (""+param);
    }
    return ret;
  }
};

// A set of conceptually related nodes.
class RunNode extends Node {
  constructor(parent, itemType, itemId) {
    super(parent, itemType, itemId);
    // sparse map of token to node
    // ( or, for repeating elements, token to array of nodes.
    this.kids= {};
  }
  toJSON() {
    return {
      id: this.id,
      type: this.type,
      value: this.kids,
    };
  }
  getKid(field) {
    return this.kids[field];
  }
  // returns the old kid ( if any )
  putField(field, newKid, vm=null) {
    const { kids } = this;
    const param= this.getParam(field);
    console.assert( !!param, `missing field ${field}`);
    // console.assert( !newKid || allTypes.areCompatible(newKid.type, param.type), `incompatible field ${field} ${newKid.type}`);
    // the types should be exactly the same, *slots* permit inherited types.
    console.assert( !newKid || newKid.type === param.type, `incompatible field ${field} ${newKid.type}`);
    //
    const oldKid= kids[field];
    if (oldKid) {
      oldKid.parent= null;
    }
    if (vm) {
      vm.$set(kids, field, newKid);
    } else {
      Vue.set(kids, field, newKid);
    }
    if (newKid) {
      newKid.parent= this;
    }
    return oldKid;
  }

  // similar to splice.
  // the first parameter is the name of the kid
  // if start is <0, delete count and addition happen from the end
  splices(field, start, deleteCount, ...newItems) {
    const param= this.getParam(field);
    console.assert( param && param.repeats, `missing or non-repeating field ${field}`);
    //
    const kids= this.getKid(field);
    if (start < 0) {
      start= Math.max(0, kids.length - deleteCount);
    }
    const removed=  kids.splice(start, deleteCount, ...newItems).map(el=>{
      el.parent= null;
      return el;
    })
    // set after clearing in case we have removed items that we are also adding.
    newItems.forEach(el=>{
      el.parent= this;
    });
    return removed;
  }
  // visit each parameter and argument in turn
  //callback(currentValue [, index [, array]]
  forEach(callback) {
    const spec= this.itemType.with;
    for (const token of spec.tokens) {
      const param= this.getParam(token);
      callback({
        token,
        param,
        kid: this.kids[token],
      });
    }
  };
  unroll(nodes, itemValue) {
    this.forEach(({idx, token, param})=>{
      if (param) {
        const arg = itemValue[token];
        if (arg !== undefined) {
          if (!param.repeats) {
            const kid= nodes.newFromItem(this, arg);
            this.kids[token]= kid;
          } else if (arg) {
            const kids= arg.map((el) => nodes.newFromItem(this, el));
            this.kids[token]= kids;
          }
        }
      }
    });
  };
};

// Swaps between a small set of options.
class SwapNode extends Node {
  constructor(parent, itemType, itemId) {
    super(parent, itemType, itemId);
    this.choice= null; // token
    this.kid= null;
  }
  toJSON() {
    const { choice, kid } = this;
    const value= choice && {
        [choice]: kid
    };
    return {
      id: this.id,
      type: this.type,
      value,
    };
  }
  unroll(nodes, itemValue) {
    if (itemValue) {
      const spec= this.itemType.with;
      for (let t=0; t< spec.tokens.length; ++t) {
        const token= spec.tokens[t];
        if (token in itemValue) {
          const kid= nodes.newFromItem(this, itemValue[token]);
          this.kid= kid;
          this.choice= token;
          break;
        }
      }
    }
  }
  putSwap(newChoice, newKid) {
    const oldKid= this.kid;
    if (oldKid) { // clear first in case newKid==oldKid
      oldKid.parent= null;
    }
    this.kid= newKid;
    this.choice= newChoice;
    if (newKid) {
      newKid.parent= this;
    }
  }
  setSwap(newChoice, newKid) {
    const node= this;
    const oldKid= node.kid;
    const oldChoice= node.choice;
    Redux.Run({
      apply() {
        node.putSwap(newChoice, newKid);
      },
      revoke() {
        node.putSwap(oldChoice, oldKid);
      }
    });
  }
};

// Pick one node from a (potentially large) set of types.
class SlotNode extends Node {
  constructor(parent, itemType, itemId) {
    super(parent, itemType, itemId);
    this.kid= null;
  }
  toJSON() {
    return {
      id: this.id,
      type: this.type,
      value: this.kid,
    };
  }
  unroll(nodes, itemValue) {
    if (itemValue.value != null) {
      const kid= nodes.newFromItem(this, itemValue);
      this.kid= kid;
    }
  }
  // fill out the passed slot with a newly created node of typeName
  putSlot(newKid) {
    const oldKid= this.kid;
    if (oldKid) { // clear first in case newKid==oldKid
      oldKid.parent= null;
    }
    this.kid= newKid;
    if (newKid) {
      newKid.parent= this;
    }
  }
  // undoable put
  setSlot(newKid) {
    const node= this;
    const oldKid= node.kid;
    Redux.Run({
      apply() {
        node.putSlot(newKid);
      },
      revoke() {
        node.putSlot(oldKid);
      }
    });
  }
};

// A leaf node representing a concrete value.
class PrimNode extends Node {
  constructor(parent, itemType, itemId) {
    super(parent, itemType, itemId);
    this.value= null;
  }
  toJSON() {
    return {
      id: this.id,
      type: this.type,
      value: this.value,
    };
  }
  unroll(nodes, itemValue) {
    this.value= itemValue;
  }
  // change a primitive value
  setPrim(newValue) {
    const node= this;
    const oldValue= node.value;
    Redux.Run({
      apply(vm) {
        node.value= newValue;
      },
      revoke() {
        node.value= oldValue;
      }
    });
  }
}

// currently, ids are serialized. that's not the long term plan.
class Ids {
  constructor(namespace= false) {
    this.counter=0;
    this.base= namespace !== false? namespace: ("id-" + Date.now().toString(16) + "-");
  }
  nextId() {
    return this.base + this.counter++;
  }
}

// node wraps items to provide a complete tree
// item is the direct serialized to disk data:
// { id string; type string; value any; } object;
class Nodes {
  constructor(pool, idNamespace=false) {
    this.pool= pool;
    this.root= null; //
    this.ids= new Ids(idNamespace);
  }
  unroll(item) {
    // newFromItem "unrolls" the item data.
    const root= this.newFromItem(null, item);
    this.root= root;
    return this;
  }
  newFromType(typeName) {
    const item= Types.createItem(typeName);
    return this.newFromItem(null, item);
  }
  newFromItem(parent, item) {
    const kid= this.newNode(parent, item.type, item.id);
    if (item.value) {
      kid.unroll(this, item.value);
    }
    return kid;
  }
  newNode(parent, typeName, itemId) {
    const newNode= {
      opt: (...args) => new SwapNode(...args),
      run: (...args)  => new RunNode(...args),
      slot: (...args) => new SlotNode(...args),
      num: (...args)  => new PrimNode(...args),
      str: (...args)  => new PrimNode(...args),
      txt: (...args)  => new PrimNode(...args),
    };
    const itemType= typeName && Types.get(typeName);
    if (typeName && !itemType && typeName[0] !== '$') {
      throw new Error(`missing type ${typeName}`);
    }
    const role= itemType.uses;
    const kid= newNode[role](
      parent,
      itemType,
      itemId || this.ids.nextId(),
    );
    if (this.pool) {
      this.pool[kid.id]= kid;
    }
    return kid;
  }
};
