Vue.component('mk-context', {
  template:
  `<div :class="cls.win"
  ><div class="mk-aux__title mk-aux__title--left"
    >context</div
    ><div
      ><div
          v-if="ctx"
          :class="bemElem('head')"
        >{{ctx.parent | titlecase}}: {{ctx.param| titlecase}} </div
      ><div
          v-if="cmd"
        ><div :class="bemElem('body')"
          ><b>{{cmd.label | titlecase}}</b
          > ( {{cmd.uses}} ): {{cmd.short}}<span
          v-if="cmd.long"
          > {{cmd.long}}</span
        ></div
      ></div
    ></div
  ></div>`,
  data() {
    return {
      cls: {
        win: [ this.bemBlock(), 'mk-aux' ],
      },
      currType: "",
      currField: null,
    };
  },
  computed: {
    cmd() {
      const t= allTypes.get(this.currType);
      return t && {
        label: Types.labelOf( t ),
        groups: Types.groupsOf( t ),
        short: Types.shortOf( t ),
        long: Types.longOf( t ),
        uses: t.uses,
      };
    },
    ctx() {
      const f= this.currField;
      return f && {
        // FIX? allow params to use the desc object setup?
        parent: Types.labelOf( f.parentType ),
        param: f.param.label,
      };
    }
  },
  methods: {
    onCmdSelected(cmdName) {
      // FIX: synchronize context display
      // this.currField= field;
      // this.currType= field && field.param.type;
    },
    onNodeSelected(node, param, token) {
      // FIX: synchronize context display
      // this.currType= typeName;
      // this.currField= null;
    }
  },
  mounted() {
    this.$root.$on("cmd-selected", this.onCmdSelected);
    this.$root.$on("node-selected", this.onNodeSelected);
  },
  beforeDestroy() {
    this.$root.$off("cmd-selected", this.onCmdSelected);
    this.$root.$off("node-selected", this.onNodeSelected);
  },
  mixins: [bemMixin()],
});

// groups... tbd. maybe in a run at the bottom of the control
// <span
//           v-if="cmd.groups.length"
//           :class="bemElem('groups')"
//           >Groups:</span
//           ><span v-for="(g,i) in cmd.groups"
//             ><template v-if="i"
//             >,</template
//             > <mk-a-button
//             >{{g| capitalize}}</mk-a-button
//         ></span
//         >
