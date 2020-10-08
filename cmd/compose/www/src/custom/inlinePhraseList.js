// an inline list of phrases.
class InlinePhraseList extends NodeList {
  constructor(nodes, para) {
    super(nodes, para, "$STORY_STATEMENT", "story_statement");
    this.inline= true;
  }
  // when we drag, we re/move everything from a given statement till the end of line.
  // returns a list of statements
  removeFrom(at, width= Number.MAX_VALUE) {
    return this.items.splice(at, width).map(el=> {
      el.parent= null;
      return el;
    });
  }
  // add a paragraph, or a line of statements
  // at the line of statements targeted
  addTo(at, paraEls) {
    let ret;
    const { node, items } = this;
    // adding a single paragraph?
    if (!Array.isArray(paraEls)) {
      // tack its elements to the end of the targeted line
      const para= paraEls;
      // remove all the kids from their parent array
      const kids= para.getKid("$STORY_STATEMENT");
      const els= kids.splice(0);
      ret= this.addTo( at, els );
    } else {
      const els= paraEls;
      items.splice(at, 0, ...els.map(el=> {
        el.parent= node;
        return el;
      }));
      ret= els.length;
    }
    return ret;
  }
}

