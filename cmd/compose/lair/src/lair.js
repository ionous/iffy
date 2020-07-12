const counts= [15, 30, 3, 5, 25, 7];
const allItems= counts.map((c) => new Lipsum(c));

// use a pair of numbers for the gutter to manage the sizing.
Vue.component('em-gutter', {
  template:
  `<div class="em-gutter"
    ><div class="em-max"
    >{{max}}</div
    ><div class="em-num"
    >{{num}}</div
  ></div>`,
  props: {
    num: Number,
    max: Number,
  },
});

// simple content with numbered gutter
Vue.component('em-item', {
  template:
  `<div class="em-item"
    ><em-gutter
      :num="num"
      :max="max"
      draggable="true"
    ></em-gutter
    ><div
      class="em-content"
    >{{text}}</div
  ></div>`,
  props: {
    num: Number,
    max: Number,
    text: String,
  }
});

Vue.component('em-table', {
  data() {
    const list= new DragList(this.items, ()=> new Lipsum());
    const dropper= this.dropper;
    return {
      drag: new DragHandler(
        dropper, {
        serializeItem(at)  {
          const item= list.items[at];
          return {
            'text/plain': item.text,
          };
        },
        removeItem(src, dst, width=1) {
          return list.removeFrom(src, dst, width);
        },
        // note: addItem might happen in a group other than serialize and removeItem.
        addItem(src, dst, rub) {
          list.addTo(src, dst, rub);
        },
      })
    }
  },
  props: {
    items: Array,
    dropper: DragHelper,
  },
  mounted() {
    this.drag.listen(this.$el);
  },
  beforeDestroy() {
    this.drag.silence();
  },
  template:
  `<div class="em-table"
    ><div
      class="em-table__header"
      :data-drag-idx="-1"
      :data-drag-edge="0"
    ></div
    ><transition-group
      name="flip-list"
      ><em-item
        v-for="(item,idx) in items"
        :class="dropper.highlight(idx)"
        :data-drag-idx="idx"
        :key="item.id"
        :num="idx*idx*idx"
        :max="1234"
        :text="item.text"
      ></em-item
    ></transition-group
    ><div
      class="em-table__footer"
      :data-drag-idx="items.length"
      :data-drag-edge="items.length-1"
    ></div
  ></div>`
});

const app= new Vue({
  el: '#app',
  data: {
    items:allItems,
    dropper: new DragHelper(),
  },
});
