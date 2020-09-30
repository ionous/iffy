
Vue.component('mk-composer', {
   template:
   `<div
      class="mk-container"
      :class="bemBlock(sidebar||'nosidebar')"
      ><mk-navigator
        :initialTab="sidebar"
        :tabs="['Compose', 'Catalog', 'Test', '(x)']"
        @navigate="navigate"
      ></mk-navigator>
      <mk-tools
      ></mk-tools
        ><mk-browser v-if="sidebar==='Compose'"
        ></mk-browser
        ><mk-catalog v-else-if="sidebar==='Catalog'"
        ></mk-catalog
        ><mk-tester v-else-if="sidebar==='Test'"
        ></mk-tester
      ><div class="mk-composer"
        :class="{'em-shift': $root.shift}"
        ><mk-switch
          :node="story"
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
  methods: {
     navigate(name) {
      this.sidebar= name;
    }
  },
  mounted() {
    this.$on("mk-button-activated",
      () => this.copier.cancel("button"));
  },
  props: {
    copier: Object,
    story: Object,
  },
  data() {
    const sidebar= "Catalog";
    return {
      sidebar
    };
  },
});


