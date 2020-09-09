class ParagraphNodes extends NodeList {
  constructor(nodes, story) {
    super(nodes, story, "$PARAGRAPH", "paragraph");
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
      const els= paraEls;
      // make a new paragraph...
      const para= this.nodes.newFromType("paragraph");
      // move els into the new paragraph
      const kids= para.getKid("$STORY_STATEMENT");
      // noting: we have to remove the default created els first.
      kids.splice(0, Number.MAX_VALUE, ...els.map(el=> {
        el.parent= para;
        return el;
      }));
      // add the paragraph to us.
      this.addTo(at, para);
    }
    return 1;
  }
}

Vue.component('mk-story-ctrl', {
  template:
  `<em-node-table
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
    // each item is a paragraph run
    return {
      list: new ParagraphNodes(root.nodes, node),
      dropper: root.dropper,
    }
  }
});
