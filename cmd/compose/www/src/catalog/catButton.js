/* a link that acts like a button.
some work with styles would be needed to share with mk-a-button
*/
Vue.component('mk-cat-button', {
  template:
  `<a class="mk-cat-button"
      :style="{ paddingLeft: (10*depth)+'px' }"
      role="button"
      @click="onActivated"
      @keydown.prevent.space="onActivated"
      @keydown.prevent.enter="onActivated"
      ><slot
      ></slot
    ></a>`,
  props: {
    depth:Number,
  },
  methods: {
    onActivated() {
      this.$root.events.$emit("mk-button-activated");
      this.$emit('activate');
    }
  },
});
