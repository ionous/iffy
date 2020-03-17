class Sibling {
  constructor(item, field, tokenIndex) {
    this.item= item;
    this.field= field;
  }

  // does field have sibling elements on its side (1,-1)
  static HasAdjacentEls(item, field, side) {
    const { value:items } = field;
    const i= items.indexOf(item);
    return side>0 ? (i+side) < items.length
                  : (i+side) >=0;
  }
  // does field have sibling elements on its side (1,-1)
  hasAdjacentEls(side) {
    return Sibling.HasAdjacentEls(this.item, this.field, side);
  }
  // return a new iterator; doesnt modify this.
  step(side) {
    var ret;
    const i= this.nextToken(side);
    if (i >= 0) {  // could be plain text if not a param
      const { field }= this;
      const nextToken= field.parentType.with.tokens[i];
      const nextItem= field.parentItem.value[ nextToken ];
      const sib= new ItemField( field.parentItem, nextToken );
      ret= new Sibling(nextItem, sib, i);
    }
    return ret;
  }
  nextToken(side) {
    var ret= -1;
    const { field }= this;
    if (field.tokenIndex >= 0) {
      const {tokens, params}= field.parentType.with;
      for (let i= field.tokenIndex+side; (i >=0 && i< tokens.length); i+=side) {
        // note: some tokens are plain text
        const nextToken= tokens[i];
        const nextParam= params[ nextToken ];
        if (nextParam) {  // could be plain text if not a param
          ret= i;
          break;
        }
      }
    }
    return ret;
  }
}
