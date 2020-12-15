// a flow contains an array of parameters.
Vue.component('mk-flow-ctrl', {
  template:
    `<span
    :class="bemBlock()"
    :data-tag="node.type"
    ><span
      v-for="el in els"
      class="mk-flow-param"
      :data-dot="el.plain"
      :data-tag="el.param && el.param.type"
      >{{el.opener}}<mk-switch
        :node="el.kid"
        :token="el.token"
        :param="el.param"
      ></mk-switch
      >{{el.closer}}</span
    ></span>`,
  computed: {
    els() {
      let els= [];
      const { node, "$root": root  } = this;
      node.forEach(({token, param, kid})=> {
        var opener, closer, plain ;
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
        });
      });
      return els;
    },
  },
  mixins: [bemMixin()],
  props: {
    node: Flow,
    param: Object,
    token: String,
  }
});
