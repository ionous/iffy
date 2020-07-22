
// use a pair of numbers for the gutter to manage the sizing.
// note: we have to put the draggable on the inner-most element
// https://bugs.chromium.org/p/chromium/issues/detail?id=982219
Vue.component('em-gutter', {
  template:
  `<div :class="cls"
    ><div class="em-max"
    >{{max}} </div
    ><div class="em-num"
    :draggable="draggable"
    @dragstart="startDrag"
    >{{grip || (num)}} </div
  ></div>`,
  props: {
    grip: String,
    num: Number,
    max: Number,
    draggable: {
      type: Boolean,
      default: true,
    }
  },
  methods: {
    startDrag(evt) {
      if (!this.draggable) {
        evt.stopPropagation();
      }
    }
  },
  computed: {
    cls() {
      const digi= !this.grip;
      return {
        "em-gutter": true,
        "em-gutter--digi":  digi,
        "em-gutter--grip": !digi,
      };
    }
  }
});

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
Vue.component('main-panel', {
  template:
  `<em-table
      :class="$root.shift && 'em-shift'"
      :items="items"
      :dropper="dropper"
    ><template
      v-slot="{item,idx}"
      ><em-table v-if="item.content.length"
          :inline="true"
          :items="item.content"
          :dropper="dropper"
          :grip="gripOf(item.outputType)"
        ><template
          v-slot="{item,idx}"
          ><span class="em-content"
          >{{item.content}} </span
        ></template
        ><template
          v-slot:footer
          ><em-gutter
            class="em-gutter--unfilled"
            :num="items.length"
            :draggable="false"
            :grip="gripOf(items[items.length-1].inputType)"
            :max="60+items.length"
          ></em-gutter
        ></template
      ></em-table
      ><em-gutter v-else
        class="em-gutter--unfilled"
        :num="idx+1"
        :draggable="false"
        :grip="gripOf(item.outputType)"
        :max="60+items.length"
      ></em-gutter
    ></template
  ></em-table>`,
  props: {
    items: Array,
    dropper: Dropper,
  },
  methods: {
    gripOf(t) {
      return itemTypes[t]
    }
  },
  computed: {
    cls() {
      const digi= !this.grip;
      return {
        "em-gutter": true,
        "em-gutter--digi":  digi,
        "em-gutter--grip": !digi,
      };
    }
  }
});
//
const app= new Vue({
  el: '#app',
  created() {
    document.addEventListener("keydown", (e) => {
      const shift= e.key === "Shift";
      this.shift= this.shift || shift;
      // console.log("keydown", e.key, shift);
    });
    document.addEventListener("keyup", (e) => {
      const shift= e.key === "Shift";
      this.shift= this.shift && !shift;
      // console.log("keyup", e.key, shift);
    });
    window.addEventListener("blur", (e) => {
      this.shift= false;
    });
  },
  data: {
    groups: [
      Lipsum.list(15, 31, 3, 5, 0, 8, 17),
      Lipsum.list(8, 12, 5, 42, 2, 17),
    ],

    dropper: new Dropper(),
    shift: false,
  },
});
