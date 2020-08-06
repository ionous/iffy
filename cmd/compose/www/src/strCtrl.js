// allows the user to pick from a set of predetermined strings arranged horizontally.
// once picked, it becomes an autocomplete style dropdown box to mutate the selection.
// fix: clicking the a-button causes a "one-frame" blue highlight leading to a visual pop
//
Vue.component('mk-str-ctrl', {
  template:
  `<span
      :class="bemBlock()"
      :data-tag="node.type"
    ><span v-if="prefix"
      >{{prefix}}</span
    ><mk-pick-inline
      v-if="!hasPicked && !editing"
      :node="node"
      @picked="onPickInline"
    ></mk-pick-inline
    ><template v-else
      ><mk-a-button
        @activate="onActivated"
      >{{showText? itemText: label}}</mk-a-button
      ><mk-auto-text
        v-if="editing"
        :key="node.id"
        :autoText="autoText"
        :initialText="showText? itemText:''"
        @change="onInputChange"
        @reject="onActivated(false)"
      ></mk-auto-text
    ></template
    ><span v-if="suffix"
      >{{suffix}}</span
  ></span>`,
  // -----------------------------------------------------
  computed: {
    hasPicked() {
      return this.node.value !== null;
    },
    autoText() {
      return new AutoTextOptions({
        autoFocus: true, // grab the focus when created.
        choices:() => {
          const {itemText, labelData}= this;
          const choices= Object.keys(labelData.map).filter(t=>t).sort((a,b)=>{
            return a===itemText? -1: a.localeCompare(b);
          });
          if (choices.length && choices[0] !== itemText) {
            choices.unshift(itemText);
          }
          return choices.concat( Object.keys(this.commandMap) );
        },
        permissive: this.permissive,  // accepts any input, not just choices
        placeholder: this.label,
        header: this.label,
      });
    },
    commandMap() {
      return this.mutation.commandMap;
    },
    itemText() {
      // find entry for value
      const { labelData } = this; // ex. "", $COMMON_NAME
      const { map, value } = labelData;
      const labelToken= Object.entries(map).find(
        ([label,token]) => token===value
      );
      // text to return is either the label, or the raw value
      return (labelToken!== undefined)? labelToken[0]: value;
    },
    // the possible choices for this str and stored in the item spec
    // ex. { "the": "$THE" }
    labelData() {
      const lts= {};
      const { tokens, params }= this.node.itemType.with;
      const filter= this.$root.filter(this.node);
      for (const token of tokens) {
        const param= params[token];
        if (param) {
          // for recapitulation ( where value is null )
          // use the empty string as the label.
          const label= (param.value !== null)? (param.label || param): "";
          lts[filter(label)]= token;
        }
      }
      return {
        map: lts,
        value: filter(this.node.value)
      };
    },
    mutation() {
      // fix: is there's a way to ask for ancestor properties?
      // the we could have a root "mutationFactory" instead of root.changes;
      // but the parent txtCtrl could happily override.
      return this.mutationFactory? this.mutationFactory():
              this.$root.newMutation( this.node);
    },
    prefix() {
      const { param } = this.node;
      return param && param.prefix;
    },
    label() {
      const { param } = this.node;
      return (param && param.label) || Filters.capitalize( param );
    },
    suffix() {
      const { param } = this.node;
      return param && param.suffix;
    }
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
      const { node } = this;
      if (choice.startsWith("/")) {
        const cmd= this.commandMap[choice];
        this.mutation.mutate( cmd );
      } else {
        const lts= this.labelData.map;
        if (choice in lts) {
          choice= lts[choice];
        }
        this.$root.setPrim( node, choice );
      }
      this.editing= false;
    },
    onPickInline(token) {
      // skip setting the user data entry key
      const { node } = this;
      const param= node.itemType.with.params[token];
      if (param.value !== null) {
        this.$root.setPrim( node, token );
      }
      this.editing= true;
      this.$root.nodeSelected(node);
    },
    onActivated(yes=true) {
      const { node } = this;
      this.editing= yes;
      this.$root.nodeSelected(node);
    },
  },
  mixins: [bemMixin()],
  props: {
    node: {
      type:PrimNode,
      required:true,
    },
    permissive: {
      type: Boolean,
      default: true,
    },
    // multiline text controls inject their own commands.
    mutationFactory: Function,
    // multiline text controls show the text elsewhere.
    showText: {
      type: Boolean,
      default: true,
    },
  }
});
