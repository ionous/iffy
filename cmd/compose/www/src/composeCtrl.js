// composer is everything in the window
// ( as opposed to the editor, which is the part that does the composing )
// fix? better names?
Vue.component('mk-composer', {
   template:
   `<div
      class="mk-container"
      :class="bemBlock(sidebar||'nosidebar')"
      ><mk-navigator
        :initialTab="sidebar"
        :tabs="['Catalog', 'Compose', 'Test', '(x)']"
        @navigate="navigate"
      ></mk-navigator>
      <mk-tools
        :catalog="catalog"
        :currentFile="currentFile"
      ></mk-tools
        ><mk-browser v-if="sidebar==='Compose'"
        ></mk-browser
        ><mk-catalog
          v-else-if="sidebar==='Catalog'"
          :catalog="catalog"
          @opened-file="onOpenedFile($event.file)"
        ></mk-catalog
        ><mk-tester v-else-if="sidebar==='Test'"
        ></mk-tester
      ><div class="mk-composer"
        :class="{'em-shift': $root.shift}"
        ><mk-switch
          :node="currentFile && currentFile.story"
          :key="currentFile && currentFile.path"
        ></mk-switch
        ><mk-trash-can
        ></mk-trash-can
        ><div class="mk-breathing-space"
        ></div
      ></div
      ><mk-context
      ></mk-context
  ></div>`,
  mixins: [bemMixin("mk-container")],
  props: {
    copier: Object,
    catalog: Cataloger,
  },
  methods: {
    navigate(name) {
      this.sidebar= name;
    },
    onOpenedFile(file) {
      this.currentFile= file;
      const ext= ".if";
      window.document.title= `${file.name.slice(0, -ext.length)} - Iffy Composer - ${file.path}`;
    }
  },
  mounted() {
    this.$on("mk-button-activated",
      () => this.copier.cancel("button"));
  },
  data() {
    const sidebar= "Catalog";
    return {
      sidebar,
      currentFile: null,
    };
  },
});


