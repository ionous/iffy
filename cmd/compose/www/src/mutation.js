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
  constructor(redux, state, extras=null, after=null) {
    this.redux= redux;
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
      const { redux, state } = this;
      if (!which) {
        redux.deleteAt(state.removes);
      } else {
        const curse= (which<0) ? state.left[-which-1]: state.right[which-1];
        // FIX: probably have to handle terminal injections: starting node == mutation
        redux.newAt(curse, which<0);
      }
    }
  }
}

// places in the tree where mutations can happen
class MutationState {
  constructor(node) {
    this.left = [];     // array of Cursor(s) indicating insertion points
    this.right= [];     // array of Cursor(s) indicating appending points
    this.removes= null; // a single Cursor or null
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
  addEdges(node, sides) {
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
        this.addEdges(c.parent, nextSides);
      }
    } //~if c is valid
  }
}
