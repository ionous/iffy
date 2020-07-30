// node wraps items to provide a complete tree
// item is the direct serialized to disk data:
// { id string; type string; value any; } object;
//
// note: item does not necessarily equal parent.item[token] b/c of arrays.
// fix? replace Node with an indexedDB based get item by id.


class Nodes {
  constructor(pool) {
    this.all= pool;
    this.nodeCounter=0; // used for generating keys
    this.root= null; //
  }
  newNode(parent, typeName) {
    const itemType= typeName && Types.get(typeName);
    if (typeName && !itemType) {
      throw new Error(`missing type ${typeName}`);
    }
    const kid= new Node(parent, itemType);
    const key= `node-${++this.nodeCounter}`;
    if (this.all) {
      this.all[key]= kid;
    }
    kid.key= key;
    return kid;
  }
  static Unroll(item, pool) {
    const nodes= new Nodes(pool);
    const root= nodes.newNode(null, item.type);
    nodes._unroll(root, item);
    nodes.root= root;
    return nodes;
  }
  _unroll(node, item) {
    const role= node.itemType.uses;
    const spec= node.itemType.with;
    switch( role ) {
    case "opt": {
      // we only expect there to be at most one token
      const val= item.value;
      if (val) {
        for (let t=0; t< spec.tokens.length; ++t) {
          const token= spec.tokens[t];
          if (token in val) {
            // the param and token in this case describe the option
            const childItem= val[token];
            const param= spec.params[token];
            const kid= this.newNode(node, childItem.type);
            kid.token= token;
            kid.param= param;
            this._unroll(kid, childItem);
            break;
          }
        }
      }
    }
    break;
    case "run": {
      // problems, of course.
      // 1. run doesnt walk data, it walks spec
      // 2. the kids expand repeats in place
      // 3. plain text tokens become {} too
      for (const token of spec.tokens) {
        if (!token.startsWith("$")) {
          const kid= this.newNode(node);
          kid.plainText= token;
        } else {
          const arg = item.value[token];
          const param = spec.params[token];
          if (arg !== undefined) {
            const kid= this.newNode(node, arg.type);
            kid.token= token;
            kid.param= param;
            //
            if (!param.repeats) {
              this._unroll(kid, arg);
            } else {
              kid.isArray= true;
              if (arg) {
                arg.forEach((i) => {
                  const el= this.newNode(kid, i.type);
                  el.inArray= true;
                  el.token= token;
                  el.param= param;
                  this._unroll(el, i);
                });
              }
            }
          }
        }
      }
    }
    break;
    case "slot": {
      const slot= item.value;
      if (slot) {
        const kid= this.newNode(node, slot.type);
        this._unroll(kid, slot);
      }
    }
    break;
    case "num":
    case "str":
    case "txt": {
      node.value= item.value;
    } break;
    default:
      throw new Error("unknown type role", role);
    }
  }
};

class Node {
  constructor(parent, itemType) {
    this.parent = parent;  // Node
    this.itemType= itemType;
    this.kids= [];
    //
    if (parent) {
      parent.kids.push(this);
    }
  }
  serialize() {
    throw new Error("reimplement");
    return JSON.stringify(story.item, 0, 2);
  }
  get type() {
    return this.itemType && this.itemType.name;
  }
  get firstChild() {
    const { kids } = this;
    return kids.length && kids[0];
  }
};
