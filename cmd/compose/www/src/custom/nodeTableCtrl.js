// we have two lists types controlled by the same control
// one of "paragraphs" ( lines ) and one of "statements" ( line elements )
Vue.component('em-node-table', {
  data() {
    const { list }= this;
    return {
      items: list.items,
      classmod: list.inline? "inline":"block",
    };
  },
  props: {
    grip:String,
    list: NodeList,
  },
  mounted() {
    const { "$root": root, list }= this;
    this.handler= new DragHandler(root.dropper, new NodeTable(list)).
                  listen(this.$el);
  },
  beforeDestroy() {
    this.handler.silence();
    this.handler= null;
  },
  methods: {
    // generate a vue css class object for an item based on the current highlight settings.
    highlight(idx) {
      const { "$root": root, list }= this;
      let highlight= false;
      let edge= false;
      const at = root.dropper.target;
      const start= root.dropper.start;
      const atList= at && (at.list === list);
      const startList= start && (start.list === list);
      if (atList) {
        edge= idx === at.target.edge;
        highlight=(idx === at.target.idx) || edge;
      }
      const mod= this.classmod;
      const inline= this.list.inline;
      return {
        "em-row": true,
        ["em-row--"+mod] : true,
        "em-drag-mark": highlight,
        "em-drag-highlight": highlight,
        "em-drag-head": edge && (at.target.idx < 0),
        "em-drag-tail": edge && (at.target.idx > 0),
        "em-drag-start": startList && ((idx === start.target.idx) || (inline && idx > start.target.idx))
      }
    }
  },
  template:
  `<div :class="['em-node-table', 'em-node-table--'+classmod]"
    ><div
      :class="['em-row', 'em-row__header--'+classmod]"
      :data-drag-idx="-1"
      :data-drag-edge="0"
    >&nbsp;</div
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
      ></em-gutter
    ></div
    ><div
      :class="['em-row', 'em-row__footer--'+classmod]"
      :data-drag-idx="items.length"
      :data-drag-edge="items.length-1"
    ><slot
        name="footer"
    >&nbsp;</slot
    ></div
  ></div>`
});
