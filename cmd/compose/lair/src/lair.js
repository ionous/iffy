
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

// simple content with numbered, draggable gutter
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
      drag: dropper.newGroup(this.group, {
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
          list.addTo(src, dst, rub, newGroup );
        },
      })
    }
  },
  props: {
    items: Array,
    group: String,
    dropper: Dropper,
  },
  mounted() {
    this.drag.handler.listen(this.$el);
  },
  beforeDestroy() {
    this.drag.handler.silence();
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
        :class="drag.highlight(idx)"
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
//
const app= new Vue({
  el: '#app',
  data: {
    g1: Lipsum.list(15, 30, 3, 5, 25, 7),
    g2: Lipsum.list(8, 12, 5, 42, 2, 17),
    dropper: new Dropper(),
  },
});
