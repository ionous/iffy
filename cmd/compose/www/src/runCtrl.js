// a run ( aka a slat ) contains an array of parameters.
Vue.component('mk-run-ctrl', {
  template:
    `<span
    :class="bemBlock()"
    :data-tag="node.type"
    ><span
      v-for="el in els"
      class="mk-run-param"
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
      // FIX: this isnt right because cursor expects index for arrays
      // and ghosts are generally, always? arrays.
      const c= new Cursor(this.node, token);
      this.$root.redux.newElem(c);
    },
  },
  computed: {
    els() {
      let els= [];
      this.node.forEach(({token, param, kid})=> {
        var opener, closer, ghost;
        if (param) { // plain text doesnt have param
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
        }
        els.push({
          kid,
          token,
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
