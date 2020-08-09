class StatementTable extends NodeTable {
  constructor(redux, node, storyStatements) {
    super(redux, node, storyStatements);
    this.inline= true;
  }
  makeBlank() {
    return this.nodes.newFromType(this.node, "story_statement");
  }
  // when we drag, we re/move everything from a given statement till the end of line.
  removeFrom(at) {
    return Node.Splice(null, this.items, at, Number.MAX_VALUE);
  }
  // add a paragraph, or a line of statements
  // at the line of statements targeted
  addTo(at, paraEls) {
    const { node, items } = this;
    // adding a single paragraph?
    if (!Array.isArray(paraEls)) {
      // tack its elements to the end of the targeted line
      const para= paraEls;
      // remove all the kids from their parent array
      const els= para.getKid("story_statement").splice(0, Number.MAX_VALUE);
      this.addTo( at, els );
    } else {
      const els= paraEls;
      Node.Splice(node, items, at, 0, els.length);
    }
  }
}

// paragraphs are actually, basically, the discrete lines of a story.
Vue.component('mk-paragraph-ctrl', {
  template:
  `<em-table
      :class="$root.shift && 'em-shift'"
      :list="list"
      :dropper="dropper"
      :grip="'\u2630'"
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
    const items= this.node.getKid("$STORY_STATEMENT");
    return {
      list: new StatementTable(root.redux, node, items),
      dropper: root.dropper,
    }
  },
});
