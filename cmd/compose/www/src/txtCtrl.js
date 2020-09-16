// multiline text: inject text _after_ a normal str-ctrl
// turn off txt
Vue.component('mk-txt-ctrl', {
  template:
  `<span
      :class="bemBlock(node.value? '':'empty')"
      :data-tag="node.type"` +
    // removed for say, because this was duplicating the a-link
    // even though there was no prefix/suffix text.
    // tdb if its really needed.
    // ><mk-str-ctrl
    //   ref="str"
    //   :node=node
    //   :permissive=false
    //   :showText=false
    //   :mutationFactory=mutationFactory
    // ></mk-str-ctrl
    `><mk-txt-edit
      v-if="editing"
      :initialText="node.value"
      :header="label"
      @close="acceptText($event)"
    ></mk-txt-edit
    ><mk-txt-lines
      :text="node.value"
      :placeholder="label"
      @click="onClick"
      @keydown.prevent.space="onClick"
      @keydown.prevent.enter="onClick"
    ></mk-txt-lines
  ></span>`,
  components: {
    // ---------------------------------
    'mk-txt-edit': {
      data() {
        return {
          // fix: initalText is sometimes null; can a validator on the prop work?
          currentText: this.initialText || "",
          closing: false,
        };
      },
      mixins: [bemMixin()],
      methods: {
        onEscape() {
          this.close(this.initialText, "escape");
        },
        onEnter() {
          this.close(this.currentText, "enter");
        },
        onOtherEnter() {
          console.log("other enter");
          if (!this.closing) {
            this.currentText+= "\n";
            Stretchy.resize(this.$refs.textarea);
          }
        },
        // user has clicked outside of the control
        // after extensive user testing it seems to feel best to accept the user input.
        onFocusOut(evt) {
          if (!evt.relatedTarget || !this.$el.contains(evt.relatedTarget)) {
            this.close(this.currentText, "clickOut");
          }
        },
        close(text, reason) {
          if (!this.closing) {
            console.log("closing", reason, text);
            this.closing= true;
            this.$emit('close', text);
          }
        }
      },
      mounted() {
        const { $refs: { textarea:input }, initialText } = this;
        const text= initialText || "";
        input.value= text;
        input.setSelectionRange(text.length, text.length);
        Stretchy.resize(input);
        input.focus();
      },
      props: {
        initialText: String,
        header:String,
      },
      template:
      `<span :class="bemBlock('header')"
          tabindex="-1"
          @focusout="onFocusOut"
        ><div v-if="true"
          :class="bemElem('header')"
          >{{header}}</div
        ><textarea ref="textarea"
          :class="bemElem('input')"
          v-model="currentText"
          @keydown.esc="onEscape"
          @keydown.prevent.exact.enter="onEnter"
          @keydown.prevent.enter="onOtherEnter"
        ></textarea
      ></span>`
    },
    // ---------------------------------
    'mk-txt-lines': {
      mixins: [bemMixin()],
      props: {
        text: String,
        placeholder: String,
      },
      render(createElement) {
        let mod= this.text? "" : "empty";
        const text= this.text || this.placeholder || "";
        const lines= text.split('\n');
        // between every line we have to createElement.
        const els= [lines[0]];
        for ( let i= 1; i< lines.length; ++i ) {
          els.push( createElement("br") );
          els.push( lines[i] );
        }

        return createElement("span", {
            class: [ this.bemBlock(mod) ],
            attrs: {
              tabindex: "0",
            },
            on: {
              click:(evt)=>{
                this.$emit('click', evt);
              },
              keydown:(evt)=>{
                this.$emit('keydown', evt);
              },
            },
          }, els);
      },
    }
  },
  data() {
    return {
      editing: false,
      label: this.param && this.param.label,
    };
  },
  methods: {
    onClick() {
      console.log("mk-txt-ctrl: clicked");
      this.editing= true;
    },
    acceptText(text) {
      const { node } = this;
      console.log("mk-txt-ctrl: acceptText");
      // ensure that we dont do this twice
      if (this.editing) {
        this.node.value= text || "";
        if (this.editing !== 1)  { // for testing.
          this.editing= false;
        }
      }
    },
    // create the mutator command list
    // adding in the "edit" command
    mutationFactory() {
      const { node } = this;
      let text= this.node.value;
      if (!text) {
        text= "/edit";
      } else {
        const words= text.split(/\s+/);
        const out= [];
        let len=0;
        for (const w of words) {
          len+= w.length+1;
          if (len >= 20) {
            if (!out.length) {
              out.push( w.slice(0, 20) );
            }
            out[out.length-1]+= "...";
            break;
          }
          out.push(w);
        }
        text= `/edit "${out.join(" ")}"`;
      }
      return this.$root.newMutation( this.node, {
        [text]: () => {
          this.editing= true;
        }
      });
    }
  },
  mixins: [bemMixin()],
  mounted() {
    // this.editing= 1; // for testing, auto-open.
  },
  props: {
    node: {
      type: PrimNode,
      required: true
    },
    param: Object,
    token: String,
  }
});
