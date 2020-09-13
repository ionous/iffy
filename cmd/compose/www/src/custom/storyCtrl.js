// a list of blocks containing, each containing a list of phrases.
// fix; would something like LinesList be more meaningful?
class ParagraphNodes extends NodeList {
  constructor(nodes, story) {
    super(nodes, story, "$PARAGRAPH", "paragraph");
  }
  // fromType is class Type.
  // returns true for any paragraph or story statement type.
  acceptsType(typeName) {
    const okay= allTypes.areCompatible(typeName, this.type) ||
                allTypes.areCompatible(typeName, "story_statement");
    return okay;
  }
  // can insert a "paragraph", a "story_statement", or any individual phrase
  insertAt(at, typeName) {
    // if we insert a paragraph, we add a node directly to the list
    // if we insert a story statement, we first insert a paragraph
    // if we insert a specific phrase, we firstew insert a story story statement
    let newItem;
    if (allTypes.implements(typeName, "story_statement")) {
      newItem= this.nodes.newFromType(typeName);
    }
    if (typeName !== "paragraph") {
      const slot= this.nodes.newFromType("story_statement");
      if (newItem) {
        slot.putSlot(newItem);
      }
      newItem= slot;
    }
    const para= this.nodes.newFromType("paragraph", 0);
    if (newItem) {
      para.splices("$STORY_STATEMENT", 0, 0, newItem);
    }
    this.spliceInto(at, para);
  }
  // add a paragraph, or a line of statements
  // at the paragraph targeted.
  addTo(at, paraEls) {
    const { node, items } = this;
    // adding a single paragraph?
    if (!Array.isArray(paraEls)) {
      const para= paraEls;
      para.parent= node;
      items.splice(at, 0, para);
    } else {
      const els= paraEls;
      // make a new paragraph...
      const para= this.nodes.newFromType("paragraph");
      // move els into the new paragraph
      const kids= para.getKid("$STORY_STATEMENT");
      // noting: we have to remove the default created els first.
      kids.splice(0, Number.MAX_VALUE, ...els.map(el=> {
        el.parent= para;
        return el;
      }));
      // add the paragraph to us.
      this.addTo(at, para);
    }
    return 1;
  }
}

Vue.component('mk-story-ctrl', {
  template:
  `<em-node-table
      :list="list"
  ><template
      v-slot="{item, idx}"
    ><mk-switch
      :node="item"
    ></mk-switch
    ></template
  ></em-node-table>`,
  props: {
    node: Node,
  },
  data() {
    const { node, "$root": root } = this;
    // each item is a paragraph run
    return {
      list: new ParagraphNodes(root.nodes, node),
      dropper: root.dropper,
    }
  }
});
