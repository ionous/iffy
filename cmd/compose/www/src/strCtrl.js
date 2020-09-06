// allows the user to pick from a set of predetermined strings arranged horizontally.
// once picked, it becomes an autocomplete style dropdown box to mutate the pick.
// fix: clicking the a-button causes a "one-frame" blue highlight leading to a visual pop
//
Vue.component('mk-str-ctrl', {
  template:
  `<span
      :class="bemBlock()"
      :data-tag="node.type"
    >{{framing.opener}}<span v-if="prefix"
      >{{prefix}}</span
    ><mk-pick-inline
      v-if="!hasPicked && !editing"
      :node="node"
      :param="param"
      :token="token"
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
  >{{framing.closer}}</span>`,
  // -----------------------------------------------------
  computed: {
    framing() {
      let opener, closer;
      const { node, pick } = this; // ex.$TEST_NAME
      const spec= node.itemType.with;
      const param= spec.params[pick];
      const filters = param && param.filters;
      if (filters) {
        if (filters.includes("quote")) {
          opener= `\u201C`;
          closer= `\u201D`;
        }
      }
      return { opener, closer };
    },
    hasPicked() {
      // technically, an unpicked value is null
      // but empty strings behave similar
      // fix: what about the appearance of empty values? multiple spaces, etc.
      return !!this.node.value;
    },
    autoText() {
      return new AutoTextOptions({
        autoFocus: true, // grab the focus when created.
        choices:() => {
          const {itemText, labelTokens}= this;
          const choices= Object.keys(labelTokens).filter(t=>t).sort((a,b)=>{
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
    // used for link text
    itemText() {
      // find entry for value
      const { labelTokens, filteredValue } = this; // ex. "", $NAME
      const labelToken= Object.entries(labelTokens).find(
        ([label,token]) => token===filteredValue
      );
      // text to return is either the label, or the raw value
      const itemText= (labelToken!== undefined)? labelToken[0]: filteredValue;
      console.assert("string" === typeof itemText, "expected text for item text");
      return itemText;
    },
    filteredValue() {
      const { value } = this.node;
      return this.filter(value || "");
    },
    mutation() {
      // fix: is there's a way to ask for ancestor properties?
      // the we could have a root "mutationFactory" instead of root.changes;
      // but the parent txtCtrl could happily override.
      return this.mutationFactory? this.mutationFactory():
              this.$root.newMutation( this.node);
    },
    // FIX! prefix, label, and suffix use node.param, which isnt a thing.
    // there's this.param, which is the parent's run argument
    // and the spec of the selected token
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
    const { node, token } = this;
    // fix: filtering depends on placement.
    // it is possible that adding an optional element in front of this node could change the filter
    // the only? way to detect that would be to either .watch for it ( somehow ) or track vue index in parent.
    const filter= this.$root.filter(node);
    // the possible choices for this str and stored in the item spec
    // ex. { "the": "$THE" }
    const labelTokens= {};
    const spec= node.itemType.with;
    for (const token of spec.tokens) {
      const param= spec.params[token];
      if (param) {
        // for recapitulation ( where the param value is null and the user can type anything. )
        // use the empty string as the label.
        const label= Node.LabelFromParam(param);
        const filteredLabel= filter(label);
        labelTokens[filteredLabel]= token;
      }
    }
    const pick= (node.value in spec.params)? node.value: labelTokens[""];
    return {
      editing: false,
      filter,
      labelTokens,
      pick,
    };
  },
  methods: {
    // value is text picked or typed.
    // we *only* send along our labels to the completion control
    onInputChange(choice) {
      const { node, labelTokens } = this;
      if (choice.startsWith("/")) {
        const cmd= this.commandMap[choice];
        this.mutation.mutate( cmd );
      } else {
        const pickedLabel= choice in labelTokens;
        // use the "anything" token
        if (!pickedLabel) {
          this.pick= labelTokens[""];
        } else {
          // the token for the specified label
          choice= labelTokens[choice];
          this.pick= choice;
        }
        this.$root.redux.setPrim( node, choice );
      }
      this.editing= false;
    },
    // which of the tokens were selected?
    // note: for "recapitulation" the token is the $(NAME_OF_TYPE)
    // and the param value is null; that's true for single entry "type anything" str controls too.
    onPickInline(token) {
      // skip setting the user data entry key
      const { node } = this;
      const spec= node.itemType.with;
      const param= spec.params[token];
      if (param.value !== null) {
        this.$root.redux.setPrim( node, token );
      }
      this.pick= token;
      this.editing= true;
      this.$root.ctrlSelected(this);
    },
    onActivated(yes=true) {
      this.editing= yes;
      this.$root.ctrlSelected(this);
    },
  },
  mixins: [bemMixin()],
  props: {
    node: {
      type:PrimNode,
      required:true,
    },
    param: Object,
    token: String,
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
