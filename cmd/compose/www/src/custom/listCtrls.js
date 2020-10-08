// helper to create components based on em-node-table
function nodeListComponent(name, { token, type, listClass= NodeList, grip }) {
  return Vue.component(name, {
    template: `<em-node-table
      :list="list"
      :grip="grip"
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
      return {
        list: new listClass(this.$root.nodes, this.node, token, type),
        grip:grip,
      }
    },
  });
};

nodeListComponent('mk-story-ctrl', {
  listClass: ParagraphNodes
});

// paragraphs are actually the discrete lines of a story.
nodeListComponent('mk-paragraph-ctrl', {
  listClass: InlinePhraseList,
  grip:'\u2630' // hamburger heaven
});

nodeListComponent('mk-activity-ctrl', {
  token: "$EXE",
  type: "execute",
});

nodeListComponent('mk-pattern-rules-ctrl', {
  listClass: PatternRules // "$PATTERN_RULE", "pattern_rule"
});

nodeListComponent('mk-pattern-locals-ctrl', {
  token: "$LOCAL_DECL",
  type: "local_decl",
});
