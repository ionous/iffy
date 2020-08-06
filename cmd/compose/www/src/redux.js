
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
  newAt(c, leftSide= false) {
    if (!("kids" in c.parent)) {
      throw new Error("cursor should target the field of a run");
    }
    if (c.isRepeatable()) {
      this.newElem(c, leftSide);
    } else {
      throw new Error("new field?")
    }
  }
  // add a new item to the left (front) or right (end) of the targeted array.
  // hmmm... regarding target and existing elements

  newElem(c, leftSide=false) {
    const { parent, target, token, param, index: i  }= c;

    const newElem= this.nodes.newFromType(parent, param.type);
    const index=  (i>=0)? i: (leftSide? 0: parent.kids.length-1);

    this.invoke({
      apply(vm) {
        let els= parent.kids[token];
        if (els === undefined) {
          els= vm.set(parent.kids, token, []);
        }
        if (leftSide) {
          els.unshift(newElem);
        } else {
          els.push(newElem);
        }
      },
      revoke(vm) {
        if (!target) {
          vm.delete(parent, token);
        } else {
          const els= parent.kids[token];
          if (leftSide) {
            els.shift();
          } else {
            els.pop();
          }
        }
      }
    });
  }
  newSlot(parent, typeName) {
    const oldSlot= parent.slot;
    const newSlot= this.nodes.newFromType(parent, typeName);
      this.invoke({
      apply() {
        parent.slot= newSlot;
      },
      revoke() {
        parent.slot= oldSlot;
      }
    });
  }
  newSwap(parent, newChoice, typeName) {
    const oldKid= parent.kid;
    const oldChoice= parent.choice;
    const newSwap= this.nodes.newFromType(parent, typeName);
      this.invoke({
      apply() {
        parent.kid= newSwap;
        parent.choice= newChoice;
      },
      revoke() {
        parent.kid= newSwap;
        parent.choice= oldChoice;
      }
    });
  }
  // remove an existing item from a field
  deleteAt(curse) {
    const { target } = curse;
    this.invoke({
      apply(vm) {
        curse.deleteMe(vm);
      },
      revoke(vm) {
        curse.spliceTarget(target, vm);
      }
    });
  }

  // change the value of an existing item
  setChild(item, newValue) {
    // resuse primitive, that's fine for now.
    return this.setPrim(item, newValue);
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
  // add a new item to a field
  addField(field, newItem) {
    if (!field.isOptional()) {
      throw new Error(`unexpected mutation ${JSON.stringify(field)}`);
    }
    if (field.item != null) {
      throw new Error(`attempting to add existing field ${JSON.stringify(field)}`);
    }
    this.invoke({
      apply(vm) {
        field.setValue(newItem, vm);
      },
      revoke() {
        field.deleteMe();
      }
    });
  }
  // we move node's item and / or the items after it into the newItem's container.
  // leftSide (aka splitBefore) the els after and including the field.
  // rightSide (aka splitAfter ) the els after field not including the field.
  // field.value is the old container
  split(node, newItem, leftSide) {
    const field= node.field;
    const parentField= node.parentNode.field;
    this.invoke({
      apply() {
        // field.item, ex. "$STORY_STATEMENT": [{ type: "story_statement", value: {} }]
        const { item } = node;
        const { value:oldItems } = field;
        const index= oldItems.indexOf(item) + (leftSide? 0: 1);
        const removed= oldItems.splice(index);
        // newItem.value holds the new container
        // we need to overwrite its placeholder item with our items.
        const newItems= newItem.value[field.token];
        newItems.splice(0, 1, ...removed);
        // and we need to put the newItem after our current group
        // regardless of which side of the item we broke on
        parentField.addRepeat( newItem );
      },
      revoke() {
        const { value:oldItems } = field;
        const newItems= newItem.value[field.token];
        const removed= newItems.splice(0); // remove all the items we added.
        oldItems.splice(oldItems.length, 0, ...removed);
        parentField.removeRepeat();
      }
    });
  }
}

