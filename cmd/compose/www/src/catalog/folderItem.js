Vue.component('mk-folder-item', {
  template:`<li
    :class="bemBlock(open?'open':'closed')"
    ><mk-a-button
      @activate="onActivated"
    >{{name}}</mk-a-button
    ><slot
    ></slot
  ></li>`,
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
  props: {
    item: CatalogItem,
  },
  methods: {
    onActivated() {
      this.$emit("activated", this.name);
    },
  },
  mixins: [bemMixin()],
});
