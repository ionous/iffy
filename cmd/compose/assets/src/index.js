makeLang(new Make(new Types()));

const redux= new Redux(Vue);
const events= new Vue(); // global event bus
//
const app= new Vue({
  el: '#app',
  methods: {
    newMutation(node, extras={}, after={}) {
      const state= new MutationState(node);
      // fix: ways to make this more generic?
      const sides= [{
        side: "left",
        label: "/break before",
      },{
        side: "right",
        label: "/break after"
      }];
      for (let k=0; k< sides.length; ++k) {
        const {side, label} = sides[k];
        const fields= state[side];
        for (let i=0; i< fields.length; ++i) {
          const field= fields[i];
          const { item } = field;
          if (item && item.type==="story_statement") {
            if (Sibling.HasAdjacentEls(item, field, k*2-1)) {
              let target= node;
              while (target.item.id !== item.id) {
                target= target.parentNode;
              }
              after[label]= () => {
                // note: the new item has a blank entry which gets overwritten.
                const containerType= target.parentNode.item.type;
                const para= Types.createItem(containerType); // ex. paragraph
                redux.split( target, para, !k );
              };
            }
          }
        }
      }
      return new Mutation(redux, state, extras, after);
    },
    setPrim(node, value) {
      redux.setPrim( node.item, value );
    },
    setChild(node, typeName) {
      console.log("setChild", typeName);
      const childItem= Types.createItem(typeName);
      redux.setChild( node.item, childItem );
      return node.newKid(childItem);
    },
    // ghosts provide trailing links for easily adding new content.
    // clicking a ghost expands into corresponding element.
    // fix? bind these better....
    isGhost(node, token) {
      const param= node.itemType.with.params[token];
      return param.filters  && param.filters.includes("ghost");
    },
    newGhost(node, token) {
      const field= new ItemField( node.item, token );
      const newItem= Types.createItem(field.param.type);
      redux.addRepeat(field, newItem);
    },
    fieldSelected(itemField) {
      this.$emit("field-selected", itemField);
    },
    cmdSelected(cmdName) {
      this.$emit("cmd-selected", cmdName);
    },
    // find the filter for displaying labels ( ex. strCtr.labelData. )
    filter(node) {
      let isAtStart= false;
      while (node.field && node.field.tokenIndex <= 0) {
        if (node.field.isRepeatable() && Sibling.HasAdjacentEls(node.item, node.field, -1)) {
          break;
        }
        node= node.parentNode;
        if (!node) {
          break;
        }
        if (!node.field || (node.item.type === "story_statement")) {
          isAtStart= true;
          break;
        }
      }
      return isAtStart? Filters.capitalize: Filters.none;
    }
  },
  data: {
    story: new Node(getStory()),
  }
});

const shortcuts= new Shortcuts(redux);
