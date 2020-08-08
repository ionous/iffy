
// paragraphs could be removed
// instead just separate groups of lines by blank lines.
// curr story has one line in the first paragraph.
Vue.component('mk-story-ctrl', {
  template:
  `<div>
    <mk-switch
      :node="unified"
   ></mk-switch
   ></div>`,
  props: {
    node: Node,
  },
  computed: {
    // FIX: save back into paragraphs
    unified() {
      const { node }= this;
      const parent= this.$root.nodes.newFromType(null, "paragraph");
      const lines= parent.getKid("$STORY_STATEMENT");

      let newLines= [];
      const ps= node.getKid("$PARAGRAPH");
      ps.forEach((p)=> {
        // each p is a node of type paragraph containing a story statement array
        const ts= p.getKid("$STORY_STATEMENT");
        ts.forEach((t)=> {
          newLines.push(t);
        });
      });
      if (newLines.length >0) {
        lines.splice(0,1, ...newLines);
      }
      return parent;
    }
  }
});
