
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

let globalRedux;

// Redux handles undo/redo
class Redux {
  // vm is a subset of Vue used for triggering change tracking.
  constructor(vm, nodes, max=500) {
    this.vm= vm;
    this.nodes= nodes;
    this.applied= new Restack(max);
    this.revoked= new Restack(max);
    this.changed= 0;
    globalRedux= this;
  }
  static Run(act) {
    globalRedux.doit(act);
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
  doit(act) {
    this.applied.push(act).apply(this.vm);
    this.revoked.clear();
    ++this.changed;
  }
}

