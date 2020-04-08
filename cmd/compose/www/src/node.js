// node wraps items to provide a complete tree
// item is the direct serialized to disk data:
// { id string; type string; value any; } object;
//
// note: item does not necessarily equal parent.item[token] b/c of arrays.
// fix? replace Node with an indexedDB based get item by id.
class Node {
  constructor(item, parentNode, token) {
    if (parentNode && parentNode.itemType.uses === 'run' && (typeof token !== 'string')) {
      throw new Error(`unexpected token '${token}'`)
    }
    this.item = item;  // Item
    this.parentNode = parentNode;  // Node
    this.field= parentNode && new ItemField( parentNode.item, token );
  }
  // fix? we could pass in an index though....
  newKid(item, token) {
    return new Node(item, this, token);
  }
  get itemType() {
    const typeName= this.item.type;
    const type= Types.get(typeName);
    if (!type) {
      throw new Error(`missing type ${typeName}`);
    }
    return type;
  }
}
