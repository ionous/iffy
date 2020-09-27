
Vue.component('mk-nav-button', {
  template:
  `<a tabindex="0"
    :class="bemBlock(selected?'selected':false)"
    @click="onActivated"
    @keydown.prevent.space="onActivated"
    @keydown.prevent.enter="onActivated"
    ><slot
    ></slot
  ></a>`,
  methods: {
    onActivated() {
      this.$root.events.$emit("mk-button-activated");
      this.$emit('activate');
    }
  },
  props: {
    selected: Boolean,
  },
  mixins: [bemMixin("mk-nav__btn")],
});


Vue.component('mk-navigator', {
  template:
  `<div class="mk-nav"
   ><mk-nav-button
     v-for="(name,i) in tabs"
     :key=i
     :selected="tab===i"
     @activate="activated(i)"
   >{{name}}</mk-nav-button
   ></div>`,
  methods: {
    activated(i) {
      this.tab=i;
      if ((i+1) === this.tabs.length) {
        this.$emit("navigate", false);
      } else {
        this.$emit("navigate", this.tabs[i]);
      }
    }
  },
  data() {
    return {
      tab: 0,
      tabs: ["Commands", "Files", "Tests", "(x)"],
    };
  },
  mixins: [bemMixin()],
});



