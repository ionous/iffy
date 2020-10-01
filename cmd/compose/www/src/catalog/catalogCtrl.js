// interacts with the list of all story files and folders
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
  props: {
    catalog: Cataloger,
  },
  mounted() {
    const { catalog } = this;
    catalog.loadFolder(this.folder);
  },
  data() {
    const that = this;
    return {
      backcat: {
        onFolder(parent, folder) {
          if (folder.contents) {
            folder.contents= false;
          } else {
            // injects the list of sub-files into the passed folder
            that.catalog.loadFolder(folder);
          }
        },
        onFile(parent, file) {
          // injects the story data into the passed file
          that.catalog.loadStory(file);
          that.$emit("opened-file", {file});
        },
      },
      folder: new CatalogFolder(""),
    }
  },
  mixins: [bemMixin()],
});

