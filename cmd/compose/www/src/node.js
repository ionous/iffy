// node wraps items to provide a complete tree
// item is the direct serialized to disk data:
// { id string; type string; value any; } object;
//
// note: item does not necessarily equal parentNode.item[token] b/c of arrays.
// fix? replace Node with an indexedDB based get item by id.


class Nodes {
  constructor(item) {
    this.all= {};
    this.nodeCounter=0; // used for generating keys
    this.root= null; //
  }
  newNode(parentNode, item, token) {
    const kid= new Node(parentNode, item, token);
    const key= `node-${++this.nodeCounter}`;
    this.all[key]= kid;
    kid.key= key;
    return kid;
  }
  static Unroll(item) {
    const nodes= new Nodes();
    const root= nodes.newNode(null, item);
    nodes._unroll(root, item);
    nodes.root= root;
    return nodes;
  }

  _unroll(node, item) {
    if (!item) {
      return;
    }
    const role= node.itemType.uses;
    const spec= node.itemType.with;
    switch( role ) {
    case "opt": {
      // we only expect there to be at most one token
      const val= item.value;
      if (val) {
        for (const token in val) {
          const childItem= val[token];
          const kid= this.newNode(node, childItem, token);
          this._unroll(kid, childItem);
          break;
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
          this.newNode(node, null, token);
        } else {
          const arg = item.value[token];
          const param = spec.params[token];
          const kid= this.newNode(node, arg, token);
          if (!param.repeats) {
            this._unroll(kid, arg);
          } else {
            kid.isArray= true;
            if (arg) {
              arg.forEach((i) => {
                const el= this.newNode(kid, i, token);
                this._unroll(el, i);
              });
            }
          }
        }
      }
    }
    break;
    case "slot": {
      const slot= item.value;
      if (slot) {
        const kid= this.newNode(node, slot);
        this._unroll(kid, slot);
      }
    }
    break;
    case "num":
    case "str":
    case "txt":
      this.newNode(node,item.value);
    break;
    default:
      throw new Error("unknown type role", role);
    }
  }

};

class Node {
  constructor(parentNode, item, token) {
    // REFACTOR
    let og= parentNode;
    if (parentNode && parentNode.isArray) {
      parentNode= parentNode.parentNode;
    }

    if (parentNode && parentNode.itemType.uses === 'run' && (typeof token !== 'string')) {
      throw new Error(`unexpected token '${token}'`)
    }
    this.item = item;  // Item
    this.token= token;
    this.parentNode = parentNode;  // Node
    this.kids= [];
    //
    if (!parentNode) {
      this.field = false;
    } else {
      this.field= new ItemField( parentNode.item, token );
      og.kids.push(this);
    }
  }
  toJSON() {
    return {
      key: this.key,
      field: this.field,
      kids: this.kids.map(k=> k.key),
    };
  }
  get plainText() {
    return (this.token && !this.token.startsWith("$")) ? this.token: false;
  }
  get itemType() {
    const typeName= this.item.type;
    const type= Types.get(typeName);
    if (!type) {
      throw new Error(`missing type ${typeName}`);
    }
    return type;
  }
  get firstChild() {
    const { kids } = this;
    return kids.length && kids[0];
  }
};
