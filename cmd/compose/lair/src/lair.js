// use a pair of numbers for the gutter to manage the sizing.
Vue.component('elem', {
  template:
  `<em-table
      :classmod="'inline'"
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
    const list= new DragList(this.items, ()=> new Lipsum());
    const dropper= this.dropper;
    return {
      drag: dropper.newGroup({
        serializeItem(at)  {
          const item= list.items[at];
          return {
            'text/plain': item.text,
          };
        },
        removeItem(src, dst, width=1, newGroup) {
          return list.removeFrom(src, dst, width, newGroup);
        },
        // note: addItem might happen in a group other than serialize and removeItem.
        addItem(src, dst, rub, newGroup) {
          list.addTo(src, dst, rub, newGroup);
        },
      })
    }
  },
  props: {
    grip:String,
    items: Array,
    dropper: Dropper,
    classmod: String,
  },
  mounted() {
    this.drag.handler.listen(this.$el, this.classmod==="inline");
  },
  beforeDestroy() {
    this.drag.handler.silence();
  },
  methods: {
    // generate a vue css class object for an item based on the current highlight settings.
    highlight(idx) {
      let highlight= false;
      let edge= false;
      const {target:at, source:from} = this.dropper;
      const atGroup= at && (at.group === this.drag);
      const fromGroup= from && (from.group === this.drag);
      if (atGroup) {
        edge= idx === at.edge;
        highlight=(idx === at.idx) || edge;
      }
      const mod= this.classmod;
      const inline= mod==="inline";
      return {
        "em-row": true,
        ["em-row--"+mod] : !!mod,
        "em-drag-mark": highlight,
        "em-drag-highlight": highlight,
        "em-drag-head": edge && (at.idx < 0),
        "em-drag-tail": edge && (at.idx > 0),
        "em-drag-from": fromGroup && ((idx === from.idx) || (inline && idx > from.idx))
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
