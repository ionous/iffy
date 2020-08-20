// a run ( aka a slat ) contains an array of parameters.
Vue.component('mk-run-ctrl', {
  template:
    `<span
    :class="bemBlock()"
    :data-tag="node.type"
    ><span
      v-for="el in els"
      class="mk-run-param"
      :data-dot="el.plain"
      :data-tag="el.param && el.param.type"
      >{{el.opener}}<mk-switch
        :node="el.kid"
        :token="el.token"
        :param="el.param"
      ></mk-switch
      >{{el.closer}}<mk-a-button
        v-if="el.ghost"
        :class="bemElem('ghost')"
        @activate="$emit('ghost', el.token)"
      >{{el.ghost}}</mk-a-button
    ></span
    ></span>`,
  methods: {
    // when the ghost is clicked, we want to expand it.
    onGhost(token) {
      const at= new Cursor(this.node, token);
      this.$root.redux.newElem(at);
    },
  },
  computed: {
    els() {
      let els= [];
      const { node, "$root": root  } = this;
      node.forEach(({token, param, kid})=> {
        var opener, closer, ghost,plain ;
        // plain text doesnt have param
        if (!param) {
          // handle caps for execute statements in a block.
          const filter= root.filter(node);
          token= filter(token);
        } else {
          const type= param.type;
          const filters = param.filters;
          if (filters) {
            if (filters.includes("quote")) {
              opener= `\u201C`;
              closer= `\u201D`;
            }
            if (filters.includes("ghost")) {
              const gtype= Types.get(param.type);
              ghost= Types.labelOf(gtype);
            }
          }
          plain= token.trim().replace(/ /g, '-').replace(/^\$/, '').toLowerCase();
        }
        els.push({
          kid,
          token,
          plain,
          param,
          opener,
          closer,
          ghost,
        });
      });
      return els;
    },
  },
  mixins: [bemMixin()],
  props: {
    node: RunNode,
    param: Object,
    token: String,
  }
});
