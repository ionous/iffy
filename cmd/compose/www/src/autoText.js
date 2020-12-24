class AutoTextOptions {
  constructor({
      autoFocus=false,      // grab the focus when created.
      choices=null,         // fn() => array { string; }
      permissive=false,     // when true, input not limited to choices.
      header= "",           // optional header element
      placeholder="Type something here",
  }) {
    this.autoFocus= autoFocus;
    this.header= header;
    this.choices= choices;
    this.permissive= permissive;
    this.placeholder= placeholder;
  }
}

// input box for items with autoComplete
// emits:
//  . change
//  . changing
//  . reject
Vue.component('mk-auto-text', {
  template:
  `<span :class="bemBlock(autoText.header?'header':false)"
      tabindex="-1"
      @focusout="onFocusOut"
    ><div v-if="autoText.header"
      :class="bemElem('header')"
      >{{autoText.header}}</div
    ><input type="text" :class="bemElem('input')"
        :placeholder="autoText.placeholder"
        :style="{ minWidth: lastText? null:'5rem'}"
        ref="input"
        @input="inputChanged"
        @focus="onFocus"
        @keydown.prevent.esc="rejectSelection"
        @keydown.prevent.enter="acceptSelection"
        @keydown.up="moveSelection(-1)"
        @keydown.down="moveSelection(1)"
        @keydown.delete="moveSelection(0)"
    ><ul ref="itemList"
        :class="bemElem('dropdown')"
        v-if="options.length"
        @click.stop="changeText($event.target.dataset.opt)"
    ><li v-for="(opt,i) in options"
        :class="[ opt.cls, sel===i? cls.itemSel:'' ]"
        :data-opt="opt.val"
      >{{ opt.val }}</li
    ></ul
  ></span>`,
  data(){
    return {
      lastText: this.initialText,
      choices: false,
      options: [],
      inError: false,
      sel: -1, // would a "last" string would be better
      changing: false,
      cls: {
        item: this.bemElem('dropitem'),
        // fix? shouldnt this be: bemElem('dropitem, 'cmd')
        itemCmd: this.bemElem('dropitem--cmd'),
        itemSel: this.bemElem('dropitem--sel'),
        initial: this.bemElem('dropitem--org'),
      }
    }
  },
  methods: {
    getChoices() {
      if (this.choices===false) {
        this.choices= this.autoText.choices();
      }
      return this.choices;
    },
    // user has clicked outside of the control
    // after extensive user testing it seems to feel best to accept the user input.
    onFocusOut(evt) {
      if (!evt.relatedTarget || !this.$el.contains(evt.relatedTarget)) {
        this.deselected();
      }
    },
    onFocus() {
      this.select();
    },
    // user has clicked outside of the control
    // after extensive user testing it seems to feel best to accept the user input.
    deselected() {
      // note: accepting the selection can change other parts of the story
      // which can drag focus away from the auto text, which can deselect it.
      // and we _dont_ want to trigger accept selection again.
      if (this.changing) {
        this.acceptSelection();
      } else {
        this.rejectSelection();
      }
      // to cancel instead:
      //this.setText(this.lastText);
    },
    select() {
      console.log("itemInput", this.$vnode.key, "selecting");
      this.input.select();
    },
    // programmatically set the input's displayed text
    // doesn't raise events
    setText(text) {
      this.options= [];  // reset our autocomplete
      this.sel= -1;
      this.changing= false;
      this.choices= false;
      //
      this.input.value= text; // our display
      this.checkForError(text);
      //Stretchy.resize(this.input);
    },
    // try to accept the input's current text or selection
    // ex. on <enter>
    acceptSelection() {
      const text= this.getCurrentValue();
      if (!this.checkForError(text)) {
        this.changeText(text);
      } else {
        console.log("itemInput", this.$vnode.key, "can't accept", JSON.stringify(text));
      }
    },
    // reset to the prior selection
    rejectSelection() {
      const text=  this.lastText;
      console.log("itemInput", this.$vnode.key, "rejecting", JSON.stringify(text));
      this.setText(text);
      this.$emit("reject", text);
    },
    // return the typed input or currently selected text
    // fix? could maybe turn into a computed property...
    // if we used v-model to change input.value into a reactive value
    getCurrentValue() {
      return (this.sel>=0 && this.sel<this.options.length)?
              this.options[this.sel].val:
              this.input.value.trimRight();
    },
    // change the currently selected drop down item
    moveSelection(dir) {
      // hitting the delete key calls moveSelection(0)
      const choices= this.getChoices();
      if (!dir || !choices.length) {
        this.sel= -1;
      } else {
        if (this.sel>= 0) {
          // no-wrap, it confuses things.
          this.sel= Math.min( Math.max( this.sel+dir, 0 ), this.options.length-1 );
        } else if (dir) {
          // first down press?
          this.input.value= "";
          this.updateOptions();
          this.checkForError("");
          this.sel= 0;
        }

        this.$nextTick(function() {
          const els= this.$refs.itemList;
          const el= els && els.children[this.sel];
          if (el) {
            if (el.scrollIntoViewIfNeeded) {
               el.scrollIntoViewIfNeeded();
            } else {
              el.scrollIntoView();
            }
          }
        });
      }
    },
    // ex. enter pressed, or selected from drop down
    // doesnt check for errors of text.
    changeText(text) {
      console.log("itemInput", this.$vnode.key, "changed", JSON.stringify(text));
      if (text) {
        const nextText= text.startsWith("/")? this.lastText: text.trim();
        this.lastText= nextText;
        this.setText(nextText);
        // testing not selecting ever
        // alternatively, will have to pass a bool for deselect to not select
        // this.select();
        this.$emit('change', text);
      }
    },
    firstIndexOf(text) {
      const choices= this.getChoices();
      return choices.findIndex(c => c.startsWith(text));
    },
    // set the input to :invalid if needed
    checkForError(text) {
      let err;
      if (text.startsWith("/")) {
        err= this.firstIndexOf(text) < 0;
      } else {
        err= text && (!this.autoText.permissive) && this.firstIndexOf(text) < 0;
      }
      if (err !== this.inError) {
        const msg= err? "select from one of the available options": "";
        this.input.setCustomValidity(msg);
        this.inError= err;
      }
      return err;
    },
    inputChanged() {
      this.updateOptions();
      const text= this.getCurrentValue();
      this.checkForError(text);

      if (!this.changing) {
        this.$emit('changing');
        this.changing= true;
      }
    },
    // could syllable separate, calc word distance when (near) empty,
    // weight by recent, etc.
    updateOptions() {
      const raw= this.input.value;
      let choices = this.getChoices();
      if (raw.length) {
        choices= choices.filter(w=> w.toLowerCase(w).indexOf(raw.toLowerCase())>=0);
      }
      const cls= this.cls;
      this.options= choices.map(c =>  {
        return {
          val: c,
          cls: {
            [cls.item]: true,
            [cls.itemCmd]: c.startsWith('/'),
            [cls.initial]: c === this.initialText
          }
        };
      });
    }
  },
  mixins: [bemMixin()],
  mounted() {
    // data is inited before mounting
    this.input= this.$refs.input;
    this.setText(this.lastText);
    if (this.autoText.autoFocus) {
      this.input.focus();
    }
  },
  props: {
    autoText: AutoTextOptions,
    initialText: {
      type:String,
      default:""
    }
  },
  watch: {
    // because we aren't using v-bind, even though this.initialText updates correctly
    // when users of the component change things, the input.value doesn't.
    initialText(newText) {
      if (this.lastText!== newText) {
        this.lastText= newText;
        this.setText(newText);
      }
    }
  },
});
