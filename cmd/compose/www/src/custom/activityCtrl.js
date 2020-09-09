class ActivityList extends NodeList {
  constructor(nodes, para) {
    super(nodes, para, "$EXE", "execute");
  }
}

// paragraphs are actually, basically, the discrete lines of a story.
Vue.component('mk-activity-ctrl', {
  template:
  `<em-node-table :list="list"
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
      list: new ActivityList(root.nodes, node),
      dropper: root.dropper,
    }
  },
});
