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
  // add label -> field for
  static add(out, prefix, side, list) {
    let dupes= {};
    for (let i=0; i< list.length; ++i) {
      const field= list[i];
      const { param, item } = field;
      if (!param && !item) {
        throw new Error( "what does this look like?" );
      }
      let txt= param? (param.label || param.type): item.type;
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
  mutate(index) {
    if (typeof index === 'function') {
      index(this)
    } else if (!index) {
      const field= this.state.removes;
      this.redux.deleteField(field);
    } else {
      const field= (index<0) ? this.state.left[-index-1]:this.state.right[index-1];
      const newItem= Types.createItem(field.param.type);
      if (!newItem) {
        throw new Error("couldn't create item");
      }
      if (field.isRepeatable()) {
        this.redux.addRepeat(field, newItem, index<0);
      } else {
        this.redux.addField(field, newItem);
      }
    }
  }
}

// backwards compat: tack item on to field
// really we should use a separate object {item, field}
function newMutableItem(item, field) {
    field.item= item;
    field.toJSON= function() {
      return {
        parent: field.parentItem ? field.parentItem.id : -1,
        token: field.token,
        item: item? item.id: null
      };
    }
    return field;
}

class MutationState {
  constructor(node) {
    this.left = [];     // array of ItemField(s) indicating insertion points
    this.right= [];     // array of ItemField(s) indicating appending points
    this.removes= null; // a single ItemField or null
    this.addEdges(node, [-1,1]);
  }
  pushRepeater(item, field, sides) {
    for (const side of sides) {
      this.pushField(item, field, side);
    }
  }
  // side: -1/+1
  pushField(item, field, side) {
    const reps= (side<0)? this.left: this.right;
    reps.push(newMutableItem(item,field));
  }
  // internal recursive
  addEdges(node, sides) {
    const field= node.field;
    if (field) {
      // at most one deletable element
      if (!this.removes && field.isDeletable()) {
        this.removes= newMutableItem(node.item, field);
      }
      // can the element be repeated?
      const repeats= field.isRepeatable();
      if (repeats) {
        this.pushRepeater(node.item, field, sides);
      }
      // check if the node is an edge
      const nextSides= [];
      let it= new Sibling(node.item, field);
      for (const side of sides) {
        // note: cant be an edge if we are in an array with sibling elements.
        if (!repeats || !it.hasAdjacentEls(side)) {
          let sib= it.step(side);
          // if it's not optional than it will have a value
          // if its repeatable then we can add to that value
          // fix? can this and the empty sib loop be merged better?
          if (sib && sib.field.isRepeatable() && !sib.field.isOptional()) {
            this.pushField(sib.item, sib.field, side);
          } else {
            for (; sib && MutationState.emptySib(sib); sib= sib.step(side)) {
              this.pushField(sib.item, sib.field, side);
            }
            // missing a valid sibling: we are an edge
            // so we will want to check this edge in our parent as well.
            if (!sib) {
              nextSides.push(side);
            }
          }
        }
      }
      if (nextSides.length) {
        this.addEdges(node.parentNode, nextSides);
      }
    } //~if field is valid
  }

  static emptySib(sib) {
    const { field } = sib;
    return field.isEmpty()  && (field.isOptional() || field.isRepeatable());
  }
}
