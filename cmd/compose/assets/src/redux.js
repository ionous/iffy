
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
  constructor(vm, max=500) {
    this.vm= vm;
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
  // remove an existing item from a field
  deleteField(field) {
    const item = field.item; // remember me
    if (field.isRepeatable()) {
      const index= field.value.indexOf(item);
      this.invoke({
        apply() {
           // remove:
          field.value.splice(index, 1);
        },
        revoke() {
          // add:
          field.value.splice(index, 0, item);
        }
      });
    } else {
      this.invoke({
        apply() {
          field.deleteMe();
        },
        revoke(vm) {
          field.setValue(item, vm);
        }
      });
    }
  }
  // add a new item to the left or right of the target field
  addRepeat(field, newItem, leftSide) {
    this.invoke({
      apply(vm) {
        field.addRepeat(newItem, leftSide, vm);
      },
      revoke() {
        field.removeRepeat(leftSide);
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

