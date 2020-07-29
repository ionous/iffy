// node wraps items to provide a complete tree
// item is the direct serialized to disk data:
// { id string; type string; value any; } object;
//
// note: item does not necessarily equal parent.item[token] b/c of arrays.
// fix? replace Node with an indexedDB based get item by id.

let nodeCounter=0;
class Node {
  constructor(item, parentNode, token) {
    // REFACTOR
    if (parentNode && parentNode.isArray) {
      parentNode= parentNode.parentNode;
    }

    if (parentNode && parentNode.itemType.uses === 'run' && (typeof token !== 'string')) {
      throw new Error(`unexpected token '${token}'`)
    }
    this.item = item;  // Item
    this.key= `node-${++nodeCounter}`;
    this.token= token;
    this.parentNode = parentNode;  // Node
    if (!parentNode) {
      this.field = false;
    } else {
      this.field= new ItemField( parentNode.item, token );
    }
    this.kids= [];
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
  // fix? we could pass in an index though....
  newKid(item, token) {
    const kid= new Node(item, this, token);
    this.kids.push(kid);
    return kid;
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

  // 1. replicate existing structure first, just precreate it.
  static Unroll(item) {
    const node= new Node(item);
    node._unroll();
    return node;
  }
  _unroll() {
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
          const kid= this.newKid(childItem, token);
          kid._unroll();
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
          this.newKid(null, token);
        } else {
          const arg = item.value[token];
          const param = spec.params[token];
          const kid= this.newKid(arg, token);
          if (!param.repeats) {
            kid._unroll();
          } else {
            kid.isArray= true;
            if (arg) {
              arg.forEach((i) => {
                const el= kid.newKid(i, token)
                el._unroll();
              });
            }
          }
        }
      }
    }
    break;
    case "slot": {
      if (item.value) {
        const kid= this.newKid(item.value);
        kid._unroll();
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
