class ParagraphTable extends NodeTable {
    constructor(redux, node) {
    super(redux, node, node.getKid("$PARAGRAPH"));
    this.inline= false;
  }
  makeBlank() {
    return this.nodes.newFromType(this.node, "paragraph");
  }
  // when we drag, we re/move a single paragraph ( a line ) at once.
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
  // add a paragraph, or a line of statements
  // at the paragraph targeted.
  addTo(at, paraEls) {
    const { node, items } = this;
    // adding a single paragraph?
    if (!Array.isArray(paraEls)) {
      const para= paraEls;
      para.parent= node;
      items.splice(at, 0, para);
    } else {
      // make a new paragraph...
      const els= paraEls;
      const para= this.makeBlank();
      // move els into the new paragraph...
      para.getKid("story_statement").splice(at, 0, ...els.map(el=> {
        el.parent= para;
        return el;
      }));
      // add the paragraph to us.
      this.addTo(at, para);
    }
  }
}

Vue.component('mk-story-ctrl', {
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
    // each item is a paragraph run
    return {
      list: new ParagraphTable(root.redux, node),
      dropper: root.dropper,
    }
  }
});
