
Vue.component('mk-tester', {
  template:
  `<div class="mk-aux" :class="bemBlock()"
  ></div>`,
   mixins: [bemMixin()],
});


