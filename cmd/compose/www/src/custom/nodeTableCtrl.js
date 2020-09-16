// we have two lists types controlled by the same control
// one of "paragraphs" ( lines ) and one of "statements" ( line elements )
Vue.component('em-node-table', {
  data() {
    const { list } = this;
    const classMod= list.inline ? "inline" : "block";
    // bemStyle
    const cls= `em-node-table em-node-table--${classMod}`;
    const rowCls= `em-row em-row--${classMod}`;
    const rowGhost= rowCls + " em-row--ghost";
    const ghostText = Types.labelOf(Types.get(list.type));
    return {
      items: list.items,
      cls,
      rowCls,
      ghostText,
    };
  },
  props: {
    grip: String,
    list: NodeList,
  },
  mounted() {
    const { "$root": root, list } = this;
    this.handler = new DragHandler(root.dropper, new NodeTableEvents(list)).
      listen(this.$el);
  },
  beforeDestroy() {
    this.handler.silence();
    this.handler = null;
  },
  computed: {
    showFooter() {
      const { "$root": root, list } = this;
      let okay = false;
      if (!list.length) {
        okay= true;
      } else {
        // always show?
        const altView= true; // root.shift || root.dropper.dragging;
        if (altView && (!list.inline || !root.blockSearch.hasBlock(list.at(-1)))) {
          const from = root.dropper.start;
          okay = !from || list.acceptsType(from.getType());
        }
      }
      return okay;
    }
  },
  methods: {
    onGhost() {
      const { list } = this;
      list.insertAt(list.length, list.type); // add ghost
    },
    dragging(idx) {
      const { "$root": root, list } = this;
      const start = root.dropper.dragging;
      return start && start.contains && start.contains(list, idx);
    },
    // generate a vue css class object for an item based on the current highlight settings.
    highlight(idx) {
      let highlight = false;
      const { "$root": root, list } = this;
      if (root.dropper.dragging) {
        // are we the target?
        const at = root.dropper.target;
        const atList = at && (at.list === list);
        if (atList) {
          highlight = idx === at.target.idx;
        }
      }
      return highlight;
    }
  },
  template:
    `<div :class="cls"
    ><div v-for="(item,idx) in items"
        v-show="!dragging(idx)"
        :class="[rowCls, {'em-drag-highlight': highlight(idx)}]"
        :data-drag-idx="idx"
        :key="item.id"
      ><em-gutter
        :num="idx+1"
        :grip="grip"
        :max="10+items.length"
      ></em-gutter
      ><slot
        :idx="idx"
        :item="item"
      ></slot
    ></div
    ><div
      v-show="showFooter"
        :class="[rowCls, 'em-row--ghost', {'em-drag-highlight': highlight(items.length)}]"
        :data-drag-idx="items.length"
        :data-drag-edge="items.length-1"
      ><em-gutter
        v-if="!list.inline"
        :grip="grip"
        :num="items.length+1"
        :max="10+items.length"
        :draggable="false"
      ></em-gutter
      ><mk-a-button
        @activate="onGhost"
    >+ {{ghostText}}</mk-a-button
    ></div
  ></div>`
});
