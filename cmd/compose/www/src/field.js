// an item in the context of its parent.
// fix: its not entirely clear why this is separate from Node
// it is used as part of "Sibling" iterators, so those might need adjust.
class ItemField {
  constructor(parentItem, token, tokenIndex) {
    this.parentItem= parentItem;
    if (!parentItem) {
      throw new Error("creating ItemField but no parent item specified");
    }
    // token is the location in parent which contains the item.
    // note: that if the parent is a "slot" there is no token.
    this.token= token;
    if (!token && this.parentType.uses === 'run') {
      console.log(JSON.stringify(parentItem,0,2));
      throw new Error("creating ItemField but no token in the parent item specified");
    }
    if (tokenIndex !== undefined) {
      this._tokenIndex= tokenIndex;
    }
  }
  toJSON() {
    return {
      parent: this.parentItem ? this.parentItem.id : -1,
      token: this.token,
      item: Array.isArray(this.value) ? "[...]": (this.value && this.value.id)
    };
  }
  get tokenIndex() {
    if (this._tokenIndex === undefined) {
      const parentType = this.parentType;
      // find the index of the field's token.
      this._tokenIndex=
           (parentType && parentType.uses === 'run') ?
            parentType.with.tokens.indexOf( this.token ):
            -1;
    }
    return this._tokenIndex;
  }

  // spec of parent item
  get parentType() {
    return Types.get(this.parentItem.type);
  }
  // object {
  //   string label
  //   string type
  //   boolean optional?
  //   boolean repeats?
  // } parameter info
  get param() {
    return this.token && this.parentType.with.params[this.token];
  }
  // value of the argument in the parent
  get value() {
    const { token, parentItem: { value }  } = this;
    return token? value[token]: value;
  }
  isEmpty() {
    // double equal for undefined.
    const { value } = this;
    return (value == null) ||
          (Array.isArray(value) && !value.length);
  isOptional() {
    const { param } = this;
    return param && param.optional;
  }
  isRepeatable() {
    const { param } = this;
    return param && param.repeats;
  }
  isDeletable() {
    const { uses } = this.parentType;
    return ( uses === 'run' ) ? this._isDeletableRun() :
           ( uses === 'opt' || uses === 'slot' )? this._isDeletableSlat() :
           false;
  }
  // future: look at whether this parameter is an item,
  // extract and store id instead.
  // ( item value is not primitive, str, num, txt )
  setValue(v, vm) {
    var ret;
    if (!this.token) {
      ret= this.parentItem.value= v;
    } else if (vm) {
      // ick. vue instances use '$set', Vue proper uses 'set'
      ret= (vm.$set || vm.set)( this.parentItem.value, this.token, v );
    } else {
      this.parentItem.value[this.token]=v;
      ret= v;
    }
    return ret;
  }
  deleteMe() {
    // yuck. opt needs to see its value nil to trigger a view change.
    // it might be nice if an empty run could be null, but its {} right now.
    if (!this.token || this.parentType.uses === 'opt') {
      // no token implies a slat or slot
      this.parentItem.value= null;
    } else {
      delete this.parentItem.value[ this.token ];
    }
  }
  addRepeat(newItem, leftSide, vm) {
    let items= this.value;
    if (items === undefined) {
      items= this.setValue([], vm);
    }
    if (leftSide) {
      items.unshift(newItem);
    } else {
      items.push(newItem);
    }
  }
  removeRepeat(leftSide) {
    const items= this.value;
    if (this.param.optional && items.length===1) {
      this.deleteMe();
    } else if (leftSide) {
      items.shift();
    } else {
      items.pop();
    }
  }
  _isDeletableRun() {
    const { param } = this; // we're in a run so it's okay to assume param exists
    return param.optional || (param.repeats && this.value.length>1);
  }
  _isDeletableSlat() {
    return !this.isEmpty();
  }
}
