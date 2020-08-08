class StoryList extends DragList {
  constructor(node, nodes, items) {
    super(items);
    this.node= node;
    this.nodes= nodes;
  }
  makeBlank() {
    return this.nodes.newFromType(this.node, "story_statement");
  }
}

// paragraphs could be removed
// instead just separate groups of lines by blank lines.
// curr story has one line in the first paragraph.
Vue.component('mk-paragraph-ctrl', {
  template:
  `<em-table
      :class="$root.shift && 'em-shift'"
      :list="list"
      :dropper="$root.dropper"
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
    const { node } = this;
    const items= this.node.getKid("$STORY_STATEMENT");
    const inline= false;
    return { // FIX -- other new Item(s) like merge.
      list: new StoryList(node, this.$root.nodes, items)
    }
  },
});


 // FIX!
