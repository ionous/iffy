class PatternRules extends NodeList {
  constructor(nodes, node) {
    super(nodes, node, "$PATTERN_RULE", "pattern_rule");
  }
  // fromType is class Type.
  // returns true for pattern_rule or any "bool_eval" type
  acceptsType(typeName) {
    const okay= (typeName === "pattern_rule") ||
                allTypes.areCompatible(typeName, "bool_eval");
    return okay;
  }
  // can insert a pattern rule or any bool_eval
  insertAt(at, typeName) {
    const rule= this.nodes.newFromType("pattern_rule", 0);
    if (typeName !== "pattern_rule") {
      const slot= this.nodes.newFromType("bool_eval");
      slot.putSlot(this.nodes.newFromType(typeName));
      rule.putField("$GUARD", slot);
    }
    this.spliceInto(at, rule);
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
      list: new PatternRules(root.nodes, node),
      dropper: root.dropper,
    }
  },
});
