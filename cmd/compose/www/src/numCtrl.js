
Vue.component('mk-num-ctrl', {
  template:
  `<span
      :class="bemBlock()"
      :data-tag="node.type"
    ><mk-a-button
        @activate="onActivated"
      >{{itemText}}</mk-a-button
      ><mk-auto-text
        v-if="editing"
        :key="node.key"
        :autoText="autoText"
        :initialText="itemText"
        @change="onInputChange"
        @reject="onActivated(false)"
      ></mk-auto-text
  ></span>`,
  // -----------------------------------------------------
  computed: {
    autoText() {
      const ctrl= this;
      return new AutoTextOptions({
        autoFocus: true, // grab the focus when created.
        choices: () => Object.keys(ctrl.commandMap),
        permissive: true,
        placeholder: this.label,
        header: this.label,
      });
    },
    commandMap() {
      return this.mutation.commandMap;
    },
    itemText() {
      return String(this.node.value);
    },
    mutation() {
      return this.$root.newMutation(this.node);
    },
    label() {
      const { param } = this.node;
      return (param && param.label) || Filters.capitalize( param );
    },
  },
  data() {
    return {
      editing: false
    };
  },
  methods: {
    // value is text picked or typed.
    // we *only* send along our labels to the completion control
    onInputChange(choice) {
      if (choice.startsWith("/")) {
        const cmd= this.commandMap[choice];
        this.mutation.mutate( cmd );
      } else {
        const value= parseFloat(choice);
        if (!isNaN(value)) {
          this.$root.setPrim( this.node, value );
        }
      }
      this.editing= false;
    },
   onActivated(yes=true) {
      this.editing= yes;
      this.$root.nodeSelected(this.node);
    },
  },
  mixins: [bemMixin()],
  props: {
    node: {
      type:Node,
      required:true,
    },
  }
});
