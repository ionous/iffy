 // FIX -- other new Item(s) like merge.
 class BlockStory extends DragList {
    constructor(node, nodes, items) {
    super(items);
    this.node= node;
    this.nodes= nodes;
    this.inline= false;
  }
  makeBlank() {
    return this.nodes.newFromType(this.node, "paragraph");
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
    const items= node.getKid("$PARAGRAPH");
    return {
      list: new BlockStory(node, root.nodes, items),
      dropper: root.dropper,
    }
  }
});
