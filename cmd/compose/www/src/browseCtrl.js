// target is from TargetFinder
class DraggableCommand extends Draggable {
  constructor(type, target) {
    super();
    this.type= type;
    this.target= target;
  }
  getDragData() {
    return {
      'text/plain': this.type,
    };
  }
  getDragImage() {
    return this.target.el;
  }
}

class TabHandler {
  constructor(tabbable) {
   this.finder= null; // see bind()
   this.tabbable= tabbable;
  }
  // bind a container of items
  // each el in containerEl should have a data-drag-idx
  bind(containerEl) {
    this.finder= containerEl? new TargetFinder(containerEl): false;
  }
  // can only drag from
  dragOver() {
    return false;
  }
  dragStart(el, dt) {
    let ret;
    const { items } = this.tabbable;
    if (items) {
      const target= this.finder.findIdx(el, true); // { el, idx, edge };
      if (target) {
        const item= items[target.idx];
        if (item) {
          ret= new DraggableCommand(item.name, target);
        }
      }
    }
    return ret;
  }
}

// a group of tabs, each with their own lists or maps of lists
// acts as a drag group
class Tabbable {
  constructor(lists) {
    this.names= Object.keys(lists);
    this.lists= lists;
    this.items= [];// items => [{name, label}]
  }
  // change the active items
  updateTab(tab, item) {
    let contents;
    const list= this.lists[tab];
    if (Array.isArray(list)) {
      contents= this._outlineOf(list);
    } else if (!item) {
      contents= this._outlineOf(Object.keys(list));
    } else {
      contents= this._outlineOf(list[item]);
    }
    this.items= contents.sort((a,b)=> a.label.localeCompare(b.label));
    return this.items;
  }
  _outlineOf(names) {
    return names.map((k)=> ({
        name:k, label:Types.labelOf( Types.get(k) )
      }));
  }

}

// displays the Tabbable commands
Vue.component('mk-browser', {
  template:
  `<div :class="cls.win"
  ><div class="mk-aux__title mk-aux__title--right"
     >commands</div
  ><div :class="bemElem('nav')"
    ><span v-for="(x,i) in tabs"
      ><template v-if="i"
      >, </template
      ><mk-a-button
          :class="[ bemElem('btn'), tab===x?cls.btnSel:false ]"
          @activate="onTab(x)"
      >{{x| capitalize}}</mk-a-button
    ></span
  >.</div
  ><div :class="bemElem('subtitle')"
     >{{item|capitalize}} {{tab|capitalize}}</div
  ><ul
    class="mk-browser-list"
    ref="browserList"
    ><li
      v-for="(k,idx) in items"
      ><span
        v-if="k.name"
        :key="k.name"
      >&#x2753;<mk-a-button
          :data-drag-idx="idx"
          draggable="true"
          @activate="onItem(k.name)"
        >{{k.label| titlecase}}</mk-a-button
      ></span
      ><span v-else>&nbsp;{{k.label| titlecase}}</span
    ></li
  ></ul
  ></div>`,
  data(){
    this.$nextTick(() => {
      this.tabs= this.tabbable.names;
      this.items= this.tabbable.updateTab(this.tab, this.item)
    });
    const tab= "phrases"; // current tab. changed by the user.
    const item= ""; // current item within a tab.
    return {
      tab,
      item,
      tabs: [],
      items: [],
      // css class helper
      cls: {
        win: [ this.bemBlock(), 'mk-aux' ],
        btnSel: this.bemElem('btn', 'sel'),
      },
    };
  },
  methods: {
    onTab(tab) {
      this.tab= tab;
      this.item= "";
      this.items= this.tabbable.updateTab(this.tab, this.item);
    },
    onItem(item) {
      if (!this.item) {
        this.item= item; // selected a sub-group
      }
      this.items= this.tabbable.updateTab(this.tab, this.item);
      this.$root.cmdSelected(item);
    },
    onNodeSelected(node, param, token) {
      // FIX: synchronize browser display
    },
  },
  created() {
    const types= allTypes.all;
    // compile groups, slots, and strs
    const all= [];
    const groups= {};
    const slots= {};

    for (const typeName in types) {
      const type= types[typeName];
      const group= type.group;
      if (group && group === "hidden") {
        continue;
      }
      all.push(typeName);
      //
      if (group && group !== "hidden") {
        const gs= Array.isArray(group)? group: [group];
        for (const g of gs) {
          const els= groups[g] || [];
          els.push(typeName);
          groups[g]= els;
        }
      }
      if (type.uses==='run') {
        const spec= type.with;
        const slotNames = spec.slots || [];
        for (const slot of slotNames) {
          const els= slots[slot] || [];
          els.push(typeName);        // add our new type
          slots[slot]= els;
        }
      }
    }
    const phrases= slots["story_statement"];
    const guards= slots["bool_eval"];
    const actions= slots["execute"];
    const tabbable= new Tabbable({
      all,
      phrases,
      actions,
      guards,
      groups,
      slots,
    });
    this.tabbable= tabbable;
  },
  mounted() {
    const { "$root": root, "$refs": refs, tabbable } = this;
    root.$on("node-selected", this.onNodeSelected);
    this.handler= new DragHandler(root.dropper, new TabHandler(tabbable)).
                  listen(refs.browserList);
  },
  beforeDestroy() {
    this.$root.$off("node-selected", this.onNodeSelected);
    this.handler.silence();
    this.handler= null;
  },
  mixins: [bemMixin()],
});


