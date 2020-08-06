// base class for the runtime story model.
class Node {
  constructor(parent, itemType, itemId) {
    this.parent = parent;  // Node
    this.id= itemId;
    this.itemType= itemType; // Type
  }
  toJSON() {
    return {
      id: this.id,
      type: this.type,
    };
  }
  serialize() {
    throw new Error("reimplement");
    return JSON.stringify(story.item, 0, 2);
  }
  unroll(nodes, itemValue) {
    throw new Error("unroll unhandled");
  }
  // return typeName
  get type() {
    return this.itemType && this.itemType.name;
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
      kids: this.kids,
    };
  }
  getChildCount() {
    return this.kids.length;
  }
  getChild(token) {
    return this.kids[token];
  }
  // visit each parameter and argument in turn
  //callback(currentValue [, index [, array]]
  forEach(callback) {
    const spec= this.itemType.with;
    const kids= this.kids;
    //
    for (const token of spec.tokens) {
      callback({
        token,
        param: spec.params[token],
        kid: kids[token],
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
    return {
      id: this.id,
      type: this.type,
      choice: this.choice,
      kid: this.kid,
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
      kid: this.kid,
    };
  }
  unroll(nodes, itemValue) {
    if (itemValue.value != null) {
      const kid= nodes.newFromItem(this, itemValue);
      this.kid= kid;
    }
  }
};

// A leaf node representing a concrete value.
class PrimNode extends Node {
  constructor(parent, itemType, itemId) {
    super(parent, itemType, itemId);
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
};

// node wraps items to provide a complete tree
// item is the direct serialized to disk data:
// { id string; type string; value any; } object;
class Nodes {
  constructor(pool) {
    this.pool= pool;
    this.nodeCounter=0; // used for generating keys
    this.root= null; //
  }
  unroll(item) {
    // newFromItem "unrolls" the item data.
    const root= this.newFromItem(null, item);
    this.root= root;
    return this;
  }
  newFromType(parent, typeName) {
    const item= Types.createItem(typeName);
    const kid= this.newNode(parent, item.type, item.id);
    kid.unroll(this, item.value);
    return kid;
  }
  newFromItem(parent, item) {
    const kid= this.newNode(parent, item.type, item.id);
    kid.unroll(this, item.value);
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
      itemId || `node-${++this.nodeCounter}`,
    );
    if (this.pool) {
      this.pool[kid.id]= kid;
    }
    return kid;
  }
};
