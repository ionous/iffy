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
      const ext= ".if";
      return item && item.name.slice(0, -ext.length);
    },
  },
  methods: {
    onActivated() {
      this.$emit("activated", this.name);
    },
  }
});
