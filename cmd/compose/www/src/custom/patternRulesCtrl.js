class PatternRules extends NodeList {
  constructor(redux, node) {
    super(redux, node, "$PATTERN_RULE", "pattern_rule");
  }
}

// paragraphs are actually, basically, the discrete lines of a story.
Vue.component('mk-pattern-rules-ctrl', {
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
    // each item is a story statement slot
    return {
      list: new PatternRules(root.redux, node),
      dropper: root.dropper,
    }
  },
});
