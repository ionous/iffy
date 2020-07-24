// use a pair of numbers for the gutter to manage the sizing.
// note: we have to put the draggable on the inner-most element
// https://bugs.chromium.org/p/chromium/issues/detail?id=982219
Vue.component('em-gutter', {
  template:
  `<div :class="cls"
    ><div class="em-max"
    >{{max}} </div
    ><div class="em-num"
    :draggable="draggable"
    @dragstart="startDrag"
    >{{grip || (num)}} </div
  ></div>`,
  props: {
    grip: String,
    num: Number,
    max: Number,
    draggable: {
      type: Boolean,
      default: true,
    }
  },
  methods: {
    startDrag(evt) {
      if (!this.draggable) {
        evt.stopPropagation();
      }
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
