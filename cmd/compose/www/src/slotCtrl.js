// when empty, present a dropdown of possible runs used to fill this slot
// otherwise expands into the selected run.
Vue.component('mk-slot-ctrl', {
  template:
  `<span
      :class="bemBlock()"
      :data-tag="node.type"
    ><template
      v-if="!childNode"
      ><mk-a-button
         :class="bemElem('item')"
          @activate="onActivated"
        >{{label}}</mk-a-button
        ><mk-auto-text
          v-if="editing"
          :key="node.id"
          :autoText="autoText"
          @change="onInputChange"
          @reject="onActivated(false)"
      ></mk-auto-text
    ></template
    ><mk-switch
      v-else
      :node="childNode"
    ></mk-switch
  ></span>`,
  computed: {
    childNode() {
      return this.node.kid;
    },
    // label => typeName
    labelTypes() {
      const ret= {};
      const { node }= this;
      // determine all of the slats that fit into this slot.
      const slats= Types.slats(node.type);
      if (slats) {
        for ( const type of slats ) {
          const label= Types.labelOf(type);
          ret[label]= type.name;
        }
      }
      return ret;
    },
    autoText() {
      const placeholder= Types.labelOf(this.node.itemType);
      return new AutoTextOptions({
        choices:() => {
          const txt= Object.keys(this.labelTypes);
          const all= txt.concat( Object.keys(this.commandMap) );
          return all;
        },
        placeholder,
        // header:"boop",
        autoFocus: true,
      });
    },
    commandMap() {
       return this.mutation.commandMap;
    },
    label() {
      return this.param? (this.param.label || Filters.capitalize( this.param )):
             Types.labelOf(this.node.itemType);
    },
    mutation() {
      return this.$root.newMutation(this.node);
    },
  },
  data() {
    return {
      editing: false,
    };
  },
  methods: {
    onActivated(yes=true) {
      this.editing= yes;
      this.$root.nodeSelected(this.node);
    },
    onInputChange(choice) {
      if (choice) {
        if (choice.startsWith("/")) {
          const cmd= this.commandMap[choice];
          this.mutation.mutate( cmd );
        } else {
          const { node } = this;
          const typeName = this.labelTypes[choice];
          if (typeName) {
            this.$root.redux.newSlot( node, typeName );
          }
        }
      }
    },
  },
  mixins: [bemMixin()],
  props: {
    node: {
      type: SlotNode,
      required: true
    },
    param: Object,
    token: String,
  }
});
