class ActivityTable extends NodeTable {
  constructor(redux, para) {
    super(redux, para, para.getKid("$EXE"));
    this.inline= false;
  }
  makeBlank() {
    return this.nodes.newFromType("execute");
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
  `<em-table
      :class="$root.shift && 'em-shift'"
      :list="list"
      :dropper="dropper"
  ><template
      v-slot="{item, idx}"
    ><mk-switch
      :node="item"
    ></mk-switch
    ></template
  ></em-table>`,
  props: {
    node: Node,
  },
  data() {
    const { node, "$root": root } = this;
    // each item is a story statement slot
    return {
      list: new ActivityTable(root.redux, node),
      dropper: root.dropper,
    }
  },
});
