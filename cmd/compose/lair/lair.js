const counts= [15, 30, 3, 5, 25, 7  ];
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
  template:
  `<div class="table"
    ><template v-for="(x,i) in items"
      ><div class="item"  draggable="true"><div class="handle">{{i*i*i}}</div
      ><div draggable="false">{{x.words.join(" ")}}</div
      ></div
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

