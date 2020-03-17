/* a link that acts like a button.
the container control can use @activate to listen to press, etc.
customizing the handler to differentiate multiple buttons.
*/
Vue.component('mk-a-button', {
  template:
  `<a tabindex="0"
      role="button"
      class="mk-a-button"
      @click="onActivated"
      @keydown.prevent.space="onActivated"
      @keydown.prevent.enter="onActivated"
      ><slot
      ></slot
      ></a>`,
  methods: {
    onActivated() {
      this.$emit('activate');
    }
  },
});
