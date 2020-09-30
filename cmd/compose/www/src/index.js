makeLang(new Make(new Types()));

const events= new Vue(); // global event bus
const nodes= new Nodes();
const redux= new Redux(Vue);

//
const app= new Vue({
  el: '#app',
  methods: {
    newMutation(node, extras={}, after={}) {
      const state= new MutationState(node);
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
        const c= Cursor.At(node);
        if (c) { // cursor can be null if parent is null.
          const includePlainText= true;
          if (c.step(-1, includePlainText)) {
            break;
          }
        }
        node= node.parent;
      }
      return isAtStart? Filters.capitalize: Filters.none;
    }
  },
  created() {
    this.redux= redux;
    this.blockSearch= new BlockSearch("activity","paragraph","pattern_rules");
    this.events= events;
    const catalog= (typeof MockCatalog !== "undefined")?
                new MockCatalog(nodes):
                new RemoteCatalog(nodes);
    this.shortcuts= new Shortcuts(redux, catalog, this.copier);
    this.catalog= catalog;
  },
  data() {
    return {
      nodes: nodes,
      dropper: new Dropper(this),
      shift: false,
      copier: {
        active: false,
        cancel(reason) {
          this.active= false;
          console.log("copier deactivated", reason);
        },
        start(reason) {
          this.active= true;
          console.log("copier activated", reason);
        },
      },
    };
  }
});
