const counts= [15, 30, 3, 5, 25, 7];
const items= counts.map((c) => new Lipsum(c));

// note: css can generate the numbers dynamically
// and, content ex. attr(draggable),
// it doesnt seem that implicitly generated grid items can be targeted though
// https://github.com/w3c/csswg-drafts/issues/1943
// so, lets just do this the manual way
Vue.component('em-handle', {
  template:
  `<div class='em-handle'></div>`
  });

Vue.component('em-table', {
  props: {
    items: Array,
  },
  methods: {
    onDragStart(evt, idx) {
      const item= this.items[idx];
      const text= item.text;
      const json= JSON.stringify({item: idx, cnt: item.words.length});
      const dt= evt.dataTransfer;
      dt.setData('text/plain', text);
      dt.setData('application/json', json);
      dt.effectAllowed= 'move';
      //
      const el= this.$refs.item[idx]
      dt.setDragImage(el,-5,10);
      //
      console.log("drag start", idx, json);
    },
  },
  template:
  `<div class="table"
    ><template v-for="(item,idx) in items"
      ><div draggable="true"
        class="handle"
        @dragstart.stop="onDragStart($event, idx)"
      >{{idx*idx*idx}}</div
      ><div
        ref="item"
      >{{item.text}}</div
    ></template
   ></div>`
});

const app= new Vue({
  el: '#app',
  // methods: {},
  data: {
    items,
  }
});

