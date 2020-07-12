const counts= [15, 30, 3, 5, 25, 7];
const allItems= counts.map((c) => new Lipsum(c));


// use a pair of numbers for the gutter to manage the sizing.
Vue.component('em-gutter', {
  template:
  `<div class="em-gutter"
    ><div class="em-len"
    >{{max}}</div
    ><div class="em-num"
    >{{num}}</div
  ></div>`,
  props: {
    num: Number,
    max: Number,
  },
});

// emits @dragon(idx, el,
Vue.component('em-item', {
  template:
  `<div class="em-item"
      :class="dropper.highlight(idx)"
      :data-drag-idx="idx"
    ><em-gutter
      :num="num"
      :max="max"
      draggable="true"
    ></em-gutter
    ><div
      class="em-content"
    >{{item.text}}</div
  ></div>`,
  props: {
    idx: Number,
    num: Number,
    max: Number,
    dropper: DragHelper,
    item: Object,
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
          const item= list.items[parseInt(at)];
          return {
            'text/plain': item.text,
          };
        },
        removeItem(src, dst, width=1) {
          return list.removeFrom(parseInt(src), parseInt(dst), width);
        },
        // note: addItem might happen in a group other than serialize and removeItem.
        addItem(src, dst, rub) {
          list.addTo(parseInt(src), parseInt(dst), rub);
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
      :data-drag-edge="-1"
      :class="dropper.highlight(-1)"
    ></div
    ><transition-group
      name="flip-list"
      ><em-item
        v-for="(item,idx) in items"
        :dropper="dropper"
        :key="item.id"
        :item="item"
        :idx="idx"
        :num="idx*idx*idx"
        :max="1234"
      ></em-item
    ></transition-group
    ><div
      class="em-table__footer"
      :data-drag-edge="items.length"
      :class="dropper.highlight(items.length)"
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
