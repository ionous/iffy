Vue.component('mk-opt-ctrl', {
  template:
  `<span
      :class="bemBlock()"
      :data-tag="node.item.type"
    ><mk-switch
      v-if="hasPicked"
      :node="childNode"
    ></mk-switch
    ><mk-pick-inline v-else
      :node="node"
      @picked="onPick"
    ></mk-pick-inline
  ></span>`,
  computed: {
    hasPicked() {
      return !!this.node.item.value;
    },
  },
  data() {
    const { node } = this;
    const childItem= node.item.value;
    return {
      childNode: childItem? this.node.newKid(childItem): null
    };
  },
  methods: {
    onPick(token) {
      const { node } = this;
      const opt= node.itemType.with.params[token];
      // options can map straight to their value.
      const which= opt.type || opt;
      this.childNode= this.$root.setChild( node, which );
    },
  },
  mixins: [bemMixin()],
  props: {
    node: {
      type:Node,
      required:true,
    }
  }
});
