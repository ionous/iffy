// the folder control holds lists of file and folder items.
// each folder item contains a folder ctrl.
Vue.component('mk-folder-ctrl', {
  template:
  `<ol class="mk-folder-ctrl"
    ><mk-folder-item
      v-for="item in folders"
      :key="item.fullpath"
      :item="item"
      @activated="onFolder(item)"
      ><mk-folder-ctrl
        :folder="item"
        :backcat="backcat"
      ></mk-folder-ctrl
    ></mk-folder-item
    ><mk-file-item
      v-for="item in files"
      :key="item.fullpath"
      :item="item"
      @activated="onFile(item)"
    ></mk-file-item
  ></ol>`,
  props: {
    folder: Object,
    backcat: Object,
  },
  computed: {
    folders() {
      return this.items(true);
    },
    files() {
      return this.items(false);
    }
  },
  methods: {
    items(isFolder) {
      const { folder } = this;
      return folder.contents? folder.contents.filter((el)=> {
        return (el instanceof CatalogFolder) === isFolder;
      }): [];
    },
    onFolder(item) {
      const { backcat, folder } = this;
      backcat.onFolder(folder,item);
    },
    onFile(item) {
      const { backcat, folder } = this;
      backcat.onFile(folder,item);
      console.log("FILE", item.name);
    },
  },
});
