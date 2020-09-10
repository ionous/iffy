makeLang(new Make(new Types()));

const events= new Vue(); // global event bus
const nodes= new Nodes();
const redux= new Redux(Vue);

//
const app= new Vue({
  el: '#app',
  mixins: [shiftMixin()],
  methods: {
    newMutation(node, extras={}, after={}) {
      const state= new MutationState(node);
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
      return new Mutation(nodes, state, extras, after);
    },
    // used to sync context, browser, etc. controls
    ctrlSelected(ctrl) {
      this.$emit("node-selected", ctrl.node, ctrl.param, ctrl.token);
    },
    // used to sync context, browser, etc. controls
    // cmdName is the name of the type
    cmdSelected(cmdName) {
      this.$emit("cmd-selected", cmdName);
    },
    // find the filter for displaying labels ( ex. strCtrl.labelData. )
    filter(node) {
      let isAtStart= false;
      // find an edge
      while (node) {
        if (node.type === "story_statement") {
          isAtStart= true;
          break;
        }
        if (node.type === "execute") {
          isAtStart= true;
          break;
        }
        // elements exist to the left?
        const includePlainText= true;
        if (Cursor.At(node).step(-1, includePlainText)) {
          break;
        }
        node= node.parent;
      }
      return isAtStart? Filters.capitalize: Filters.none;
    },
  },
  created() {
    this.redux= redux;
    this.blockSearch= new BlockSearch("activity","paragraph","pattern_rules");
  },
  computed: {
    story() {
      return this.nodes.root;
    },
  },
  data: {
    nodes: nodes.unroll(getStory()),
    dropper: new Dropper(),
  }
});

const shortcuts= new Shortcuts(redux);

