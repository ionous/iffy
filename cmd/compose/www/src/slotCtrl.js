// when empty, present a dropdown of possible runs used to fill this slot
// otherwise expands into the selected run.
Vue.component('mk-slot-ctrl', {
  template:
  `<span
      :class="bemBlock()"
      :data-tag="node.item.type"
     ><template
       v-if="!hasChosen"
       ><mk-a-button
         :class="bemElem('item')"
          @activate="onActivated"
        >{{label}}</mk-a-button
        ><mk-auto-text
          v-if="editing"
          :key="node.item.id"
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
    hasChosen() {
      return this.node.item.value;
    },
    // label => typeName
    labelTypes() {
      const ret= {};
      const { node }= this;
      // get all of the slats that fit into this slot.
      const slats= Types.slats(this.slotType);
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
        placeholder: placeholder,
        // header:"boop",
        autoFocus: true,
      });
    },
    commandMap() {
       return this.mutation.commandMap;
    },
    slotType() {
      return this.field.param.type;
    },
    field() {
      return this.node.field;
    },
    label() {
      return this.field.param.label || Filters.capitalize( this.field.param );
    },
    mutation() {
      return this.$root.newMutation( this.node);
    },
    childNode() {
      return this.node.firstChild;
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
      this.$root.fieldSelected(this.field);
    },
    onInputChange(choice) {
      if (choice) {
        if (choice.startsWith("/")) {
          const cmd= this.commandMap[choice];
          this.mutation.mutate( cmd );
        } else {
          const { node } = this;
          const typeName = this.labelTypes[choice];
          const childItem= Types.createItem(typeName);
          this.$root.setChild( node, childItem );
          this.childNode= node.newKid( childItem );
        }
      }
    },
  },
  mixins: [bemMixin()],
  props: {
    node: {
      type:Node,
      required:true
    }
  }
});
