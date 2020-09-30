// app stores a folder `{ name, contents: [] }`
// where elements of contents are either strings, or the same {} pair structure;
Vue.component('mk-catalog', {
  template:
  `<div
    class="mk-aux"
    :class="bemBlock()"
    ><mk-folder-ctrl
      v-if="folder.contents"
      :folder="folder"
      :backcat="backcat"
    ></mk-folder-ctrl
  ></div>`,
  mounted() {
    const { catalog } = this;
    catalog.getFiles(this.folder);
  },
  created() {
    const { "$root": root  } = this;
    if (!root._catalogSingleton) {
      root._catalogSingleton= new MockCatalog();
    }
    this.catalog= root._catalogSingleton;
  },
  data() {
    const that = this;
    return {
      backcat: {
        onFolder(parent,folder) {
          if (folder.contents) {
            folder.contents= false;
          } else {
            // this injects the list of sub-files into the passed folder
            that.catalog.getFiles(folder);
          }
        },
        onFile(parent,file) {
          // open the file.
        },
      },
      folder: new CatalogFolder(""),
    }
  },
  mixins: [bemMixin()],
});

