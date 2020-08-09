class ParagraphTable extends NodeTable {
    constructor(redux, node, paragraphs) {
    super(redux, node, paragraphs);
    this.inline= false;
  }
  makeBlank() {
    return this.nodes.newFromType(this.node, "paragraph");
  }
  // when we drag, we re/move a single paragraph ( a line ) at once.
  removeFrom(at) {
    return Node.Splice(null, this.items, at, 1);
  }
  // add a paragraph, or a line of statements
  // at the paragraph targeted.
  addTo(at, paraEls) {
    const { node, items } = this;
    // adding a single paragraph?
    if (!Array.isArray(paraEls)) {
      const para= paraEls;
      Node.Splice(node, items, at, 0, para);
    } else {
      // need a new paragraph, and add els to it.
      const els= paraEls;
      const para= this.makeBlank();
      const kids= para.getKid("story_statement");
      Node.Splice(para, kids, 0, kids.length, ...els);
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
    const items= node.getKid("$PARAGRAPH");
    return {
      list: new ParagraphTable(root.redux, node, items),
      dropper: root.dropper,
    }
  }
});
