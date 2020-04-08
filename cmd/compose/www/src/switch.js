//------------------------------------------------
// simple slot based switch statement
// see: https://github.com/vuejs/vue/issues/8097#issuecomment-486018518
Vue.component('v-switch', {
  functional: true,
  props: {
    value: { type: [String, Number], required: true }
  },
  render(createElement, { data, props, scopedSlots:slots }) {
    const { value } = props;
    const slotFn = value in slots ?
                    slots[value] :
                    slots.default;
    return slotFn ? slotFn(data.attrs) : null;
  }
});

//------------------------------------------------
// fix? replace switch with a custom render that looks up the components by name
Vue.component('mk-switch', {
  props: {
    node: {
      type: Node,
      required: true
    }
  },
  computed: {
    uses() {
      return this.node.itemType.uses;
    }
  },
  template:
  `<v-switch :value="uses"
    ><template #num
      ><mk-num-ctrl
          :node=node
      ></mk-num-ctrl
    ></template
    ><template #opt
      ><mk-opt-ctrl
          :node=node
      ></mk-opt-ctrl
    ></template
    ><template #run
      ><mk-run-ctrl
          :node=node
      ></mk-run-ctrl
    ></template
    ><template #slot
      ><mk-slot-ctrl
          :node=node
      ></mk-slot-ctrl
    ></template
    ><template #str
      ><mk-str-ctrl
          :node=node
      ></mk-str-ctrl
    ></template
    ><template #txt
      ><mk-txt-ctrl
          :node=node
      ></mk-txt-ctrl
    ></template
  ></v-switch>`
});
