
class Restack {
  // keep a capped number of items in the stack
  constructor(max) {
    this.top= -1;
    this.count= 0;
    this.max= max;
    this.list= [];
  }
  clear() {
    this.top= -1;
    this.count= 0;
  }
  push(el) {
    const { top, count, max, list } = this;
    const next = (top+1) % max;
    list[next]= el;
    this.top= next;
    if (count<= max) {
      this.count= count+1;
    }
    return el;
  }
  pop(nothrow) {
    var ret;
    const { top, count, max, list } = this;
    if (count) {
      ret= list[top];
      this.top= (top === 0)? max - 1: top - 1;
      this.count= count-1;
    } else if (!nothrow) {
      throw new Error("nothing to undo")
    }
    return ret;
  }
}

// Redux handles undo/redo
class Redux {
  // vm is a subset of Vue used for triggering change tracking.
  constructor(vm, nodes, max=500) {
    this.vm= vm;
    this.nodes= nodes;
    this.applied= new Restack(max);
    this.revoked= new Restack(max);
    this.changed= 0;
  }
  // throws if the undo stack is empty
  undo(nothrow) {
    let okay= false;
    const act= this.applied.pop(nothrow);
    if (act) {
      this.revoked.push(act).revoke(this.vm);
      --this.changed; // negative is okay.
      okay= true;
    }
    return okay;
  }
  // throws if the redo stack is empty
  redo(nothrow) {
    let okay= false;
    const act= this.revoked.pop(nothrow);
    if (act) {
      this.applied.push(act).apply(this.vm);
      ++this.changed;
      okay= true;
    }
    return okay;
  }
  // { function(vm) apply, revoke; }
  invoke(act) {
    this.applied.push(act).apply(this.vm);
    this.revoked.clear();
    ++this.changed;
  }
  // we can generically create optional members of runs
  // cursor c, must target a member of a run
  newAt(at, leftSide= false) {
    if (!("kids" in at.parent)) {
      throw new Error("cursor should target the field of a run");
    }
    if (at.isRepeatable()) {
      this._newElem(at, leftSide);
    } else {
      this._newField(at);
    }
  }
  // add a new item to a field
  _newField(at) {
    if (at.isRepeatable()) {
      throw new Error(`newField should target a non-repeatable field ${JSON.stringify(at)}`);
    }
    const { parent, token, param }= at;
    const newField= this.nodes.newFromType(param.type);
    this.invoke({
      apply(vm) {
        const { kids } = parent;
        vm.set(kids, token, newField);
        newField.parent= parent;
      },
      revoke(vm) {
        const { kids } = parent;
        vm.delete(kids, token);
        newField.parent= null;
      }
    });
  }
  _newElem(at, leftSide=false) {
    if (!at.isRepeatable()) {
      throw new Error(`newElem should target a repeatable field ${JSON.stringify(at)}`);
    }
    const { parent, token, param, index }= at;
    const newElem= this.nodes.newFromType(param.type);
    this.invoke({
      apply(vm) {
        // if the field doesnt exist, add the new node via a new array.
        const { kids } = parent;
        const field= kids[token];
        if (!field) {
          vm.set(kids, token, [newElem]);
        } else if (index<0) {
          field.push(newElem); // no specific element targeted, append.
        } else {
          const i= leftSide? index: index+1;
          field.splice(i, 0, newElem);
        }
        newElem.parent= parent;
      },
      revoke(vm) {
        const { kids } = parent;
        const field= kids[token];
        if (field.length <= 1) {
          vm.delete(kids, token);
        } else {
          // re-determine the index to avoid left/right side issues.
          const rub= field.indexOf(newElem);
          field.splice(rub, 1);
        }
        newElem.parent= null;
      }
    });
  }
  // remove an existing child targeted by the passed cursor
  // note: this will happily delete non-optional elements.
  deleteAt(at) {
    const { parent, token, index }= at;
    const oldKid= at.target;
    const oldChoice= parent.choice;

    this.invoke({
      apply(vm) {
        if (!token) { // no token means swap or slot
          parent.kid= null;
          if (oldChoice!== undefined) {
            parent.choice= null;
          }
        } else {
          const { kids } = parent;
          const field= kids[token];
          if (field) {
            // delete the field, or remove a single element?
            if ((index >= 0) && (field.length > 1)) {
              field.splice(index, 1);
            } else {
              vm.delete(kids, token);
            }
          }
        }
      },
      revoke(vm) {
        if (!token) { // no token means swap or slot
          parent.kid= oldKid;
          if (oldChoice!== undefined) {
            parent.choice= oldChoice;
          }
        } else {
          const { kids } = parent;
          if ((index >= 0) && (token in kids)) {
            const field= kids[token];
            field.splice(index, 0, oldKid);
          } else {
            const value= (index<0)? oldKid: [oldKid];
            vm.set(kids, token, value);
          }
        }
      }
    });
  }
  newSlot(parent, typeName) {
    const oldSlot= parent.slot;
    const newSlot= this.nodes.newFromType(typeName);
      this.invoke({
      apply() {
        parent.slot= newSlot;
        newSlot.parent= parent;
      },
      revoke() {
        parent.slot= oldSlot;
        newSlot.parent= null;
      }
    });
  }
  newSwap(parent, newChoice, typeName) {
    const oldKid= parent.kid;
    const oldChoice= parent.choice;
    const newSwap= this.nodes.newFromType(typeName);
    this.invoke({
      apply() {
        parent.kid= newSwap;
        parent.choice= newChoice;
        newSwap.parent= parent;
      },
      revoke() {
        parent.kid= newSwap;
        parent.choice= oldChoice;
        newSwap.parent= null;
      }
    });
  }
  // change a primitive value
  setPrim(item, newValue) {
    const oldValue= item.value;
    this.invoke({
      apply(vm) {
        item.value= newValue;
      },
      revoke() {
        item.value= oldValue;
      }
    });
  }

  // we move node's item and / or the items after it into the newItem's container.
  // leftSide (aka splitBefore) the els after and including the field.
  // rightSide (aka splitAfter ) the els after field not including the field.
  // field.value is the old container
  // split(node, newItem, leftSide) {
  //   const field= node.field;
  //   const parentField= node.parentNode.field;
  //   this.invoke({
  //     apply() {
  //       // field.item, ex. "$STORY_STATEMENT": [{ type: "story_statement", value: {} }]
  //       const { item } = node;
  //       const { value:oldItems } = field;
  //       const index= oldItems.indexOf(item) + (leftSide? 0: 1);
  //       const removed= oldItems.splice(index);
  //       // newItem.value holds the new container
  //       // we need to overwrite its placeholder item with our items.
  //       const newItems= newItem.value[field.token];
  //       newItems.splice(0, 1, ...removed);
  //       // and we need to put the newItem after our current group
  //       // regardless of which side of the item we broke on
  //       parentField.addRepeat( newItem );
  //     },
  //     revoke() {
  //       const { value:oldItems } = field;
  //       const newItems= newItem.value[field.token];
  //       const removed= newItems.splice(0); // remove all the items we added.
  //       oldItems.splice(oldItems.length, 0, ...removed);
  //       parentField.removeRepeat();
  //     }
  //   });
  // }
}

