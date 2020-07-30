// allow the user to pick from an item's tokens displayed horizontally
Vue.component('mk-pick-inline', {
  template:
  `<ol :class="bemBlock(lines.length>1 && 'pad')"
      ><li v-for="t in lines" :class="bemElem('item')"
        ><mk-a-button v-if="t.opt"
            :class="bemElem('opt')"
            :data-opt="t.opt"
            @activate="onPick(t.opt)"
          >{{t.text}}</mk-a-button
        ><template v-else
            :class="bemElem('text')"
          >{{t.text}}</template
      ></li
   ></ol>`,
  computed: {
    lines() {
      const { node } = this;
      const spec = node.itemType.with;
      // when there's only a single option,
      // use the field's label instead of our own inline label.
      const solo = (spec.tokens.length == 1) && node.param.label;
      return spec.tokens.map(t => {
        const opt= spec.params[t];
        return opt ? {
          opt: t,
          text: solo || opt.label || opt,
          clsOpt: "opt"
        }: {
          text: t,
          plain: "plain-text"
        };
      });
    }
  },
  methods: {
    onPick(opt) {
      const { node } = this;
      console.log("pickInline", node.key, "picked", opt);
      this.$emit("picked", opt);
    }
  },
  mixins: [bemMixin()],
  props: {
    node: {
      type: Node,
      required:true,
    }
  }
});
