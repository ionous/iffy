makeLang(new Make(new Types()));

const events= new Vue(); // global event bus
const nodes= new Nodes();
const redux= new Redux(Vue, nodes);

//
const app= new Vue({
  el: '#app',
  mixins: [shiftMixin()],
  methods: {
    newMutation(node, extras={}, after={}) {
      const state= new MutationState();
      state.addEdges(node, [-1,1]);
      // REFACTOR
      // fix: ways to make this more generic?
      // fix -- this exists to create paragraphs,
      // changing this and the node hierarchy to handle drag and drop
      // const sides= [{
      //   side: "left",
      //   label: "/break before",
      // },{
      //   side: "right",
      //   label: "/break after"
      // }];
      // for (let k=0; k< sides.length; ++k) {
      //   const {side, label} = sides[k];
      //   const fields= state[side];
      //   for (let i=0; i< fields.length; ++i) {
      //     const field= fields[i];
      //     const { item } = field;
      //     if (item && item.type==="story_statement") {
      //       if (Sibling.HasAdjacentEls(item, field, k*2-1)) {
      //         let target= node;
      //         while (target.id !== item.id) {
      //           target= target.parentNode;
      //         }
      //         after[label]= () => {
      //           // note: the new item has a blank entry which gets overwritten.
      //           const containerType= target.parentNode.type;
      //           const para= Types.createItem(containerType); // ex. paragraph
      //           redux.split( target, para, !k );
      //         };
      //       }
      //     }
      //   }
      // }
      return new Mutation(redux, state, extras, after);
    },
    nodeSelected(node) {
      this.$emit("node-selected", node);
    },
    cmdSelected(cmdName) {
      this.$emit("cmd-selected", cmdName);
    },
    // find the filter for displaying labels ( ex. strCtrl.labelData. )
    filter(node) {
      let isAtStart= false;
      while (node.field && node.field.tokenIndex <= 0) {
        if (node.field.isRepeatable() && Sibling.HasAdjacentEls(node, -1)) {
          break;
        }
        node= node.parentNode;
        if (!node) {
          break;
        }
        if (!node.field || (node.type === "story_statement")) {
          isAtStart= true;
          break;
        }
      }
      return isAtStart? Filters.capitalize: Filters.none;
    },
    dumpStory() {
      return this.story.serialize();
    }
  },
  computed: {
    story() {
      return this.nodes.root;
    },
  },
  data: {
    redux: redux,
    nodes: nodes.unroll(getStory()),
    dropper: new Dropper(),
  }
});

const shortcuts= new Shortcuts(redux);

