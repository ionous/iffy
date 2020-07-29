
Vue.component('mk-num-ctrl', {
  template:
  `<span
      :class="bemBlock()"
      :data-tag="node.item.type"
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
      return String(this.node.item.value);
    },
    mutation() {
      return this.$root.newMutation(this.node);
    },
    field() {
      return this.node.field;
    },
    label() {
      const param= this.checkParam();
      return param.label || Filters.capitalize( param );
    },
  },
  data() {
    return {
      editing: false
    };
  },
  methods: {
    checkParam() {
      const field= this.field;
      if (!field.param) {
        console.log(this.node);
        throw new Error("missing param");
      }
      return field.param;
    },
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
      this.$root.fieldSelected(this.field);
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
