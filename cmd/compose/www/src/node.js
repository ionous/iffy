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
    root._unroll(nodes);
    nodes.root= root;
    return nodes;
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

  _unroll(nodes) {
    const item= this.item;
    if (!item) {
      return;
    }
  // something, something, something
  // switch on item role, the role is going to have to come from spec.
    const role= this.itemType.uses;
    const spec= this.itemType.with;
    switch( role ) {
    case "opt": {
      // we only expect there to be at most one token
      const val= item.value;
      if (val) {
        for (const token in val) {
          const childItem= val[token];
          const kid= nodes.newNode(this, childItem, token);
          kid._unroll(nodes);
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
          nodes.newNode(this, null, token);
        } else {
          const arg = item.value[token];
          const param = spec.params[token];
          const kid= nodes.newNode(this, arg, token);
          if (!param.repeats) {
            kid._unroll(nodes);
          } else {
            kid.isArray= true;
            if (arg) {
              arg.forEach((i) => {
                const el= nodes.newNode(kid, i, token)
                el._unroll(nodes);
              });
            }
          }
        }
      }
    }
    break;
    case "slot": {
      if (item.value) {
        const kid= nodes.newNode(this,item.value);
        kid._unroll(nodes);
      }
    }
    break;
    case "num":{
    }
    break;
    case "str":{
    }
    break;
    case "txt":{
    }
    break;
    default:
      throw new Error("unknown type role", role);
    }
  }
};
