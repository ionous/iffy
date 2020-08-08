// FIX -- other new Item(s) like merge.
class InlineStory extends DragList {
  constructor(node, nodes, items) {
    super(items);
    this.node= node;
    this.nodes= nodes;
    this.inline= true;
  }
  makeBlank() {
    return this.nodes.newFromType(this.node, "story_statement");
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
    const items= this.node.getKid("$STORY_STATEMENT");
    return {
      list: new InlineStory(node, this.$root.nodes, items),
      dropper: this.$root.dropper,
    }
  },
});
