Vue.component('mk-file-item', {
  template:`<li
    :class="bemBlock()"
  ><mk-cat-button
    :class="bemElem('button')"
    :depth="depth"
    @activate="onActivated"
  >{{name}}</mk-cat-button
  ></li>`,
  mixins: [bemMixin()],
  props: {
    item: CatalogFile,
    depth: Number,
  },
  computed: {
    name() {
      const { item }= this;
      return item && item.name.slice(0, item.name.length-6);
    },
  },
  methods: {
    onActivated() {
      this.$emit("activated", this.name);
    },
  }
});
