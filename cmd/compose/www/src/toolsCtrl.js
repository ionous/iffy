Vue.component('mk-tools', {
  template:
  `<div class="mk-tools mk-aux">
      <button :class="bemElem('btn')"  disabled>Play</button>
      <button :class="bemElem('btn')" :disabled="!allowTesting" @click="onTest">Check</button>
      <span :class="bemElem('msg')" v-if="copying">copying...</span>
      <span :class="bemElem('msg')" v-if="msg">{{msg}}</span>
    </div>`,
  mixins: [bemMixin()],
  props: {
    catalog: Cataloger,
    currentFile: CatalogFile,
  },
  data() {
    return {
      msg: "",
      testing: false,
    }
  },
  computed: {
    copying() {
      const { copier } = this.$root;
      return copier.active;
    },
    allowTesting() {
      return this.currentFile && !this.testing;
    }
  },
  methods: {
    onTest() {
      const { currentFile } = this;
      if (!currentFile) {
        throw new Error("nothing to test");
      }
      this.msg= "Connecting...";
      this.testing= true;

      this.catalog.run("check", currentFile, {}, (ok)=>{
        this.testing= false;
        this.msg= ok? "": "An unknown error occurred.";
      });
    },
  },
});


