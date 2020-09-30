Vue.component('mk-folder-item', {
  template:
  `<li
    ><mk-cat-button
      :class="bemElem('button', open?'open':'closed')"
      :depth="depth"
      @activate="onActivated"
    >{{name}}</mk-cat-button
    ><slot
    ></slot
  ></li>`,
  mixins: [bemMixin()],
  props: {
    item: CatalogFolder,
    depth: Number
  },
  computed: {
    name() {
      const { item }= this;
      return item && item.name;
    },
    open() {
      const { item }= this;
      return item && item.contents;
    }
  },
  methods: {
    onActivated() {
      this.$emit("activated", this.name);
    },
  },
});
