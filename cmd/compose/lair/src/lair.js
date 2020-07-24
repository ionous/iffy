
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
const app= new Vue({
  el: '#app',
  mixins: [shiftMixin()],
  data: {
    groups: [
      Lipsum.list(15, 31, 3, 5, 0, 8, 17),
      Lipsum.list(8, 12, 5, 42, 2, 17),
    ],

    dropper: new Dropper()
  },
});
