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
      :data-order="el.order"
      :data-role="el.role"
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

      // node.itemType -> Type, with.slots

      node.forEach(({token, param, role, kid})=> {
        var opener, closer, plain ;
        // plain text doesnt have param
        if (!param) {
          if (!els.length && role && role.charAt(0) == "C") {
            token= token.charAt(0).toUpperCase() + token.slice(1);
          }
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
          role,
          plain,
          param,
          opener,
          closer,
          order: `t${els.length}`,
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
