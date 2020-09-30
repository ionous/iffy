Vue.component('mk-file-item', {
  template:`<li
      class="mk-file-item"
    ><mk-a-button
      @activate="onActivated"
    >{{name}}</mk-a-button
  ></li>`,
  computed: {
    name() {
      const { item }= this;
      return item && item.name.slice(0, item.name.length-6);
    },
  },
  props: {
    item: CatalogItem,
  },
  methods: {
    onActivated() {
      this.$emit("activated", this.name);
    },
  }
});
