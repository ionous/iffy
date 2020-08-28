class ActivityList extends NodeList {
  constructor(redux, para) {
    super(redux, para, "$EXE");
    this.inline= false;
  }
  makeBlank() {
    return this.redux.nodes.newFromType("execute");
  }

  // when we drag, we re/move a single execute ( a line ) at once.
  // returns a single statement
  removeFrom(at) {
    var one;
    const rub= this.items.splice(at, 1);
    if (rub.length) {
      one= rub[0];
      one.parent= null;
     }
     return one;
  }
  addTo(at, exe) {
    const { node, items } = this;
    exe.parent= node;
    items.splice(at, 0, exe);
  }
}

// paragraphs are actually, basically, the discrete lines of a story.
Vue.component('mk-activity-ctrl', {
  template:
  `<em-node-table
      :class="$root.shift && 'em-shift'"
      :list="list"
  ><template
      v-slot="{item, idx}"
    ><mk-switch
      :node="item"
    ></mk-switch
    ></template
  ></em-node-table>`,
  props: {
    node: Node,
  },
  data() {
    const { node, "$root": root } = this;
    // each item is a story statement slot
    return {
      list: new ActivityList(root.redux, node),
      dropper: root.dropper,
    }
  },
});
