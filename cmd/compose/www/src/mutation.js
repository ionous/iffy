// helper to generate maps of slash commands
class CommandMap {
  // return an ordered map of label -> state index
  static NewFromState( state ) {
    let out= {};
    const { removes } = state;
    if (removes) {
      CommandMap.add(out, "remove", 0, [removes]);
    }
    CommandMap.add(out, "append", +1, state.right);
    CommandMap.add(out, "insert", -1, state.left);
    return out;
  }
  // list = [Cursor]
  static add(out, prefix, side, list) {
    let dupes= {};
    for (let i=0; i< list.length; ++i) {
      const c= list[i];
      const { target, param } = c;

      let txt= param? (param.label || param.type): target.type;
      const counter= dupes[txt] || 0;
      dupes[txt]= counter+1;
      if (counter) {
        txt= `${txt} (${counter})`;
      }
      let msg= `/${prefix} ${txt}`;
      out[msg]= side*(i+1);
    }
  }
}

class Mutation {
  // fix? i dont like these extras/afters objects
  constructor(nodes, state, extras=null, after=null) {
    this.nodes= nodes;
    this.state= state;
    this.commandMap= CommandMap.NewFromState(state);
    if (extras) {
      this.commandMap= Object.assign(extras, this.commandMap, after)
    }
  }
  // <0:left, >0:right, 0:remove.
  mutate(which) {
    if (typeof which === 'function') {
      which(this)
    } else {
      const { state } = this;
      if (!which) {
        this.deleteAt(state.removes);
      } else {
        const curse= (which<0) ? state.left[-which-1]: state.right[which-1];
        // FIX: probably have to handle terminal injections: starting node == mutation
        this.newAt(curse, which<0);
      }
    }
  }

  // remove an existing child targeted by the passed cursor
  // note: this will happily delete non-optional elements.
  deleteAt(at) {
    const { parent, token, index }= at;
    const oldKid= at.target;
    const oldChoice= parent.choice;

    Redux.Run({
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
    Redux.Run({
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
    Redux.Run({
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
}

// places in the tree where mutations can happen
class MutationState {
  constructor(node) {
    this.focus= node;
    this.left = [];     // array of Cursor(s) indicating insertion points
    this.right= [];     // array of Cursor(s) indicating appending points
    this.removes= null; // a single Cursor or null
    this._addEdges(node, [-1,1]);
  }
  toJSON() {
    const { left, right, removes } = this;
    return {
      left, right, removes,
    }
  }
  // side: -1/+1
  _remember(c, side) {
    const reps= (side<0)? this.left: this.right;
    reps.push(c);
  }
  // the basic idea is this:
  // starting from a node, look at its left and right siblings.
  // if the node is repeatable, remember that
  // if there is a missing optional sibling: remember that, and look at the next sibling.
  // if you hit an edge -- that is, if you have no sibling --
  // move up to the parent node, and repeat.
  _addEdges(node, sides) {
    const c= Cursor.At(node);
    if (c) {
      // at most one deletable element
      if (!this.removes && c.isDeletable()) {
        this.removes= c;
      }
      // can the current element be repeated?
      const repeats= c.isRepeatable();
      if (repeats) {
        sides.forEach(side => this._remember(c, side));
      }
      // next sides tracks whether we should look up to the parent or not.
      // -- only if we are at an edge.
      const nextSides= [];
      for (const side of sides) {
        let ok= false;
        for (let sib= c; sib; ) {
          // fix: optimize the internal token index lookup
          sib= sib.step(side);
          if (!sib) { // no sibling, we're an edge.
            ok = true;
          } else if (!sib.target) { // empty sibling, remember it.
            this._remember(sib, side);
          } else {
            // we have a sibling, we're not an edge
            // *unless* that sibling is a separate kid, and it repeats
            // in which case we should be able to add to it.
            if ((c.token !== sib.token) && sib.isRepeatable()) {
              this._remember(sib, side);
            }
            break; // since we're not on a (true) edge, end.
          }
        }
        if (ok) {
          nextSides.push(side);
        }
      };
      //
      if (nextSides.length) {
        this._addEdges(c.parent, nextSides);
      }
    } //~if c is valid
  }
}
