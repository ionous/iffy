// we have two lists types controlled by the same control
// one of "paragraphs" ( lines ) and one of "statements" ( line elements )
Vue.component('em-node-table', {
  data() {
    const { list } = this;
    return {
      items: list.items,
      classmod: list.inline? "inline": "block",
      ghost: Types.labelOf(Types.get(list.type)),
    };
  },
  props: {
    grip:String,
    list: NodeList,
  },
  mounted() {
    const { "$root": root, list }= this;
    this.handler= new DragHandler(root.dropper, new NodeTable(root.nodes, list)).
                  listen(this.$el);
  },
  beforeDestroy() {
    this.handler.silence();
    this.handler= null;
  },
  computed: {
    tablecls() {
        const { "$root": root, classmod }= this;
        return {'em-node-table':true,
                [`em-node-table--${classmod}`]:true,
                'em-shift': root.shift}
    }
  },
  methods: {
    onGhost() {
      const { "$root": root, list }= this;
      list.insertAt(list.length, root.nodes.newFromType(list.type));
    },
    // generate a vue css class object for an item based on the current highlight settings.
    highlight(idx) {
      const { "$root": root, list }= this;
      let highlight= false;
      let edge= false;

      const start= root.dropper.start;
      
      const startList= start && (start.list === list);
      const atStart= startList && ((idx === start.target.idx) || (inline && idx > start.target.idx))


      // are we the target
      const at = root.dropper.target;
      const atList= at && (at.list === list);
      if (atList) {
        edge= idx === at.target.edge;
        highlight=(idx === at.target.idx) || edge;
      }
      const inline= this.list.inline;
      const mod= "em-row--"+this.classmod;
      return {
        "em-row": true,
        [mod] : true,
        "em-drag-mark": highlight,
        "em-drag-highlight": highlight,
        "em-drag-start": atStart,
        "em-row--ghost": idx === -1,
      };
    }
  },
  mixins: [bemMixin()],
  template:
  `<div :class="tablecls"
    ><div v-for="(item,idx) in items"
        :class="highlight(idx)"
        :data-drag-idx="idx"
        :key="item.id"
      ><em-gutter
        :num="idx+1"
        :grip="grip"
        :max="60+items.length"
      ></em-gutter
      ><slot
        :idx="idx"
        :item="item"
      ></slot
    ></div
    ><div v-if="$root.shift"
        :class="highlight(-1)"
        :data-drag-idx="items.length"
        :data-drag-edge="items.length-1"
      ><em-gutter
        v-if="!list.inline"
        :grip="grip"
        :max="60+items.length"
      ></em-gutter
      ><mk-a-button
        @activate="onGhost"
    >+ {{ghost}}</mk-a-button
    ></div
  ></div>`
});
