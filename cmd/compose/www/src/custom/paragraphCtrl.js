
//
Vue.component('xmk-paragraph-ctrl', {
  template:
  `<em-table
      :class="$root.shift && 'em-shift'"
      :items="items"
      :dropper="$root.dropper"
  ><template
      v-slot="{item,idx}"
    >....</template
  ></em-table>`,
  props: {
    node: Node,
  },
  computed: {
    items() {
      return this.node.getChildAt("$STORY_STATEMENT");
    }
  }
});


// "$STORY_STATEMENT": [{
//   "id": "id-a-0",
//   "type": "story_statement",
//   "value": {
//    }
// ]
