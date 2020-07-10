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
      :class="focus.highlight(idx)"
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
    focus: DragHelper,
    item: Object,
  }
});
Vue.component('em-table', {
  data() {
    const items= this.items;
    const helper= this.dragHelper;
    return {
      drag: new DragHandler(
        helper, {
        serializeItem(idx)  {
          const item= items[idx];
          return {
            'text/plain': item.text,
            'application/json': JSON.stringify(item.words),
          };
        },
        removeItem(idx) {
          // needs some thought.
        },
      })
    }
  },
  props: {
    items: Array,
    dragHelper: DragHelper,
  },
  template:
  `<div
      class="em-table"
      @dragstart="drag.onDragStart($event)"
      @dragenter="drag.onDragItem($event)"
      @dragover="drag.onDragItem($event)"
      @dragleave="drag.onDragLeave($event)"
      @dragend="drag.onDragEnd($event)"
      @drag="drag.onDragUpdate($event)"
    ><transition-group
      name="flip-list"
      ><em-item
        v-for="(item,idx) in items"
        :focus="dragHelper"
        :key="item.id"
        :item="item"
        :idx="idx"
        :num="idx*idx*idx"
        :max="1234"
      ></em-item
    ></transition-group
    ><div
      class="em-table__footer"
      :data-drag-idx="-1"
      :data-drag-prev="items.length-1"
    ></div
  ></div>`
});

const app= new Vue({
  el: '#app',
  data: {
    items:allItems,
    dragHelper: new DragHelper(),
  },
});
