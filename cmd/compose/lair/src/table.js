Vue.component('em-table', {
  data() {
    const { items, inline, dropper }= this;
    const makeBlank= ()=> new Item(inline?"":[]);
    const list= new DragList(items, inline, makeBlank);
    return {
      list,
      handler: new DragHandler(new DragGroup(list, dropper)),
    };
  },
  computed: {
    classmod() {
      return this.inline? "inline":"block";
    }
  },
  props: {
    grip:String,
    items: Array,
    dropper: Dropper,
    inline: Boolean,
  },
  mounted() {
    this.handler.listen(this.$el);
  },
  beforeDestroy() {
    this.handler.silence();
  },
  methods: {
    // generate a vue css class object for an item based on the current highlight settings.
    highlight(idx) {
      let highlight= false;
      let edge= false;
      const {target:at, start} = this.dropper;
      const atList= at && (at.list === this.list);
      const startList= start && (start.list === this.list);
      if (atList) {
        edge= idx === at.edge;
        highlight=(idx === at.idx) || edge;
      }
      const mod= this.classmod;
      const inline= this.inline;
      return {
        "em-row": true,
        ["em-row--"+mod] : true,
        "em-drag-mark": highlight,
        "em-drag-highlight": highlight,
        "em-drag-head": edge && (at.idx < 0),
        "em-drag-tail": edge && (at.idx > 0),
        "em-drag-start": startList && ((idx === start.idx) || (inline && idx > start.idx))
      }
    }
  },
  template:
  `<div :class="['em-table', 'em-table--'+classmod]"
    ><div
      :class="['em-row', 'em-row__header--'+classmod]"
      :data-drag-idx="-1"
      :data-drag-edge="0"
    ></div
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
    ></slot
    ></div
  ></div>`
});
//
