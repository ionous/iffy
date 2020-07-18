// use a pair of numbers for the gutter to manage the sizing.
Vue.component('elem', {
  template:
  `<em-table
      :inline="true"
      :items="subitems"
      :dropper="dropper"
      :grip="'\u201C'"
    ><template
      v-slot="{item,idx}"
      ><span class="em-content"
      >{{item.text}}</span
    ></template
  ></em-table>`,
  props: {
    idx:Number,
    item:Object,
    dropper:Object,
  },
  data(){
    const id= this.item.id;
    const text= this.item.text;
    const subitems= text.split(". ").map((x, i) => {
      // fake items:
      return {
        id: `${id}-${i}`,
        text: `${x.trim().replace(/\.|,$/,'')}.`,
      }
    });
    return {
      subitems: subitems,
    }
  }
});

// use a pair of numbers for the gutter to manage the sizing.
// note: we have to put the draggable on the inner-most element
// https://bugs.chromium.org/p/chromium/issues/detail?id=982219
Vue.component('em-gutter', {
  template:
  `<div class="em-gutter"
    ><div class="em-max"
    >{{max}}</div
    ><div class="em-num"
    draggable="true"
    >{{grip || num}}</div
  ></div>`,
  props: {
    grip: String,
    num: Number,
    max: Number,
  },
  beforeDestroy() {
    console.log(`gutter ${this.num}/${this.max} being destroyed`);
  }
});

Vue.component('em-table', {
  data() {
    const list= new DragList(this.items, this.inline, ()=> new Lipsum());
    const dropper= this.dropper;
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
      const {target:at, source:from} = this.dropper;
      const atList= at && (at.list === this.list);
      const fromList= from && (from.list === this.list);
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
        "em-drag-from": fromList && ((idx === from.idx) || (inline && idx > from.idx))
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
    ></div
  ></div>`
});
//
const app= new Vue({
  el: '#app',
  created() {
    document.addEventListener("keydown", (e) => {
      console.log("keydown", e.key === "Shift");
      this.shift= true;
    });
    document.addEventListener("keyup", (e) => {
      console.log("keyup", e.key === "Shift");
      this.shift= false;
    });
    window.addEventListener("blur", (e) => {
      this.shift= false;
    });
  },
  data: {
    groups: [
      Lipsum.list(15, 31, 3, 5, 8, 17),
      Lipsum.list(8, 12, 5, 42, 2, 17),
    ],
    dropper: new Dropper(),
    shift: false,
  },
});
