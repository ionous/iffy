const counts= [15, 30, 3, 5, 25, 7];
const items= counts.map((c) => new Lipsum(c));

// use a pair of numbers for the gutter to manage the sizing.
Vue.component('em-gutter', {
  template:
  `<div class="em-gutter"
    ><div class="em-len"
    >{{len}}</div
    ><div class="em-num"
    >{{num}}</div
  ></div>`,
  props: {
    len: Number,
    num: Number,
  },
});

Vue.component('em-item', {
  template:
  `<div class="em-item"
      ref="item"
      @dragstart.stop="onDragStart($event)"
    ><em-gutter
      :num="idx*idx*idx"
      :len="1234"
      draggable="true"
    ></em-gutter
    ><div
      class="em-content"
    >{{item.text}}</div
  ></div>`,
  methods: {
    onDragStart(evt) {
      const idx= this.idx;
      const item= this.item;
      const text= item.text;
      const json= JSON.stringify({item: idx, cnt: item.words.length});
      const dt= evt.dataTransfer;
      dt.setData('text/plain', text);
      dt.setData('application/json', json);
      dt.effectAllowed= 'move';
      // //
      const el= this.$refs.item;
      el.classList.add("dragging");
      dt.setDragImage(el,10,10);
      setTimeout(()=>{
        el.classList.remove("dragging");
      });
    },
  },
  props: {
    item: Object,
    idx: Number,
  }
});

Vue.component('em-table', {
  props: {
    items: Array,
  },
  template:
  `<div class="table"
    ><em-item
      v-for="(item,idx) in items"
      :key="item.id"
      :item="item"
      :idx="idx"
    ></em-item
   ></div>`
});

const app= new Vue({
  el: '#app',
  // methods: {},
  data: {
    items,
  }
});
