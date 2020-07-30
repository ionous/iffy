// a run ( aka a slot ) contains an array of parameters.

Vue.component('mk-run-ctrl', {
  template:
    `<span
    :class="bemBlock()"
    :data-tag="node.type"
    ><span
      v-for="param in params"
      class="mk-run-param"
      :data-tag="param.type"
      >{{param.head}}<mk-switch
        :node=param.node
        :key="param.node.key"
      ></mk-switch
      >{{param.tail}}<mk-a-button
        v-if="param.ghost"
        :class="bemElem('ghost')"
        @activate="$emit('ghost', param.node.token)"
      >{{param.ghost}}</mk-a-button
    ></span
    ></span>`,
  methods: {
    // when the ghost is clicked, we want to expand it.
    onGhost(token) {
      this.$root.newGhost(this.node, token);
    },
  },
  computed: {
    params() {
      return this.node.kids.map((kid) => {
        var head, tail, ghost, type;
        const { param }= kid;
        if (param) {
          type= param.type;
          const { filters }= param;
          if (filters) {
            if (filters.includes("quote")) {
              head= `\u201C`;
              tail= `\u201D`;
            }
            if (filters.includes("ghost")) {
              const gtype= Types.get(param.type);
              ghost= Types.labelOf(gtype);
            }
          }
        }
        return {
          node:kid,
          head,
          tail,
          ghost,
          type,
        };
      });
    },
  },
  mixins: [bemMixin()],
  props: {
    node: Node,
  }
});
