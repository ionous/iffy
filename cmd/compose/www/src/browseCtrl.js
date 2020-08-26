// a list of { name, labels } representing commands used by the composer
class TabList extends Map {
  constructor({sansContent=false, txtContent=false}={}) {
    super();
    this.sansContent= sansContent;
    this.txtContent= txtContent;
  }
  // returns { label, name }
  outline() {
    return Array.from(this.keys(), (k) => ({
        name:k, label:Types.labelOf( Types.get(k) )
     }));
  }
  // return { label, name } of the sub items.
  contents(k) {
    let ret= [];
    if (!k || this.sansContent) {
        ret= this.outline();
    } else {
      const names= this.get(k);
      if (this.txtContent) {
        ret= Array.from(names, (k) => ({
            label:k
        }));
      } else {
        ret= Array.from(names, (k) => ({
            name:k, label:Types.labelOf( Types.get(k) )
        }));
      }
    }
    return ret.sort((a,b)=> a.label.localeCompare(b.label));
  }
}

// a group of tabs, each with their own list of {name, label} commands.
class Tabbable {
  constructor(lists) {
    this.finder= null; // see bind()
    this.names= Object.keys(lists);
    this.lists= lists;
    this.items= [];// items => [{name, label}]
  }
  // change the active items
  updateTab({tab, item}) {
    this.items= this.lists[tab].contents(item);
  }
  // bind a container of items
  // each el in containerEl should have a data-drag-idx
  bind(containerEl) {
    this.finder= containerEl? new TargetFinder(containerEl): false;
  }
  // can only drag from
  dragOver(start) {
    return false;
  }
  dragStart(el, dt) {
    let okay;
    if (this.items) {
      const start= this.finder && this.finder.findIdx(el, true); // { el, idx, edge };
      if (start) {
        Dropper.setDragData(dt, el, this._serializeItem(start));
        okay= true;
      }
    }
    return okay;
  }
  _serializeItem(start) {
    const item= this.items[start.idx];
    return {
      'text/plain': item.name,
    };
  }
}

// displays the Tabbable commands
Vue.component('mk-browser', {
  template:
  `<div :class="cls.win"
  ><div class="mk-aux__title mk-aux__title--right"
     >commands</div
  ><div :class="bemElem('nav')"
    ><span v-for="(x,i) in tabbable.names"
      ><template v-if="i"
      >, </template
      ><mk-a-button
          :class="[ bemElem('btn'), tab===x?cls.btnSel:false ]"
          @activate="onTab(x)"
      >{{x| capitalize}}</mk-a-button
    ></span
  >.</div
  ><div :class="bemElem('subtitle')"
     >{{item}} {{tab|capitalize}}</div
  ><ul
    class="mk-browser-list"
    ref="browserList"
    ><li
      v-for="(k,idx) in tabbable.items"
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
    const types= allTypes.all;
    // compile groups, slots, and strs
    const all= new TabList({sansContent:true});
    const groups= new TabList(); // group name => [ runs that implement the group ]
    const slots= new TabList();
    const str= new TabList({txtContent:true});
    //
    for (const typeName in types) {
      const type= types[typeName];
      all.set(typeName, type);
      //
      const group= type.group;
      if (group) {
        const gs= Array.isArray(group)? group: [group];
        for (const g of gs) {
          const typeNames= groups.get(g) || [];
          typeNames.push(typeName);
          groups.set(g, typeNames);
        }
      }
      switch (type.uses) {
        case 'run': {
          const slotNames = type.with.slots || [];
          for (const slot of slotNames) {
            const typeNames= slots[slot] || [];
            typeNames.push(typeName);        // add our new type
            slots.set(slot, typeNames); // write it back
          }
          break;
        }
        case 'opt': {
          const opts= Object.values( type.with.params ).map((x)=>x.type || x);
          slots.set(typeName, opts);
          break;
        }
        case 'str': {
          // FIX: it feels frustrating that definitions for user types
          // are in anyway different -- storage wise -- than custom types
          // ie. where is our "spec" db?
          const vals= Object.values( type.with.params ).map((x)=>x.label||x);
          str.set(typeName, vals);
          break;
        }
      };
    }
    const tabbable= new Tabbable({
      all,
      groups,
      slots,
      str,
    });
    const { "$root": root } = this;
    const dropper= root.dropper;
    const handler= new DragHandler(dropper, tabbable);
    this.$nextTick(function() {
      this.onTab("all");
    });
    return {
      // current tab. changed by the user.
      tab: "groups",
      // current item within a tab.
      item: "",
      // css class helper
      cls: {
        win: [ this.bemBlock(), 'mk-aux' ],
        btnSel: this.bemElem('btn', 'sel'),
      },
      tabbable,    // list of clickable tab names
      handler, // drag-drop listener
    };
  },
  methods: {
    onTab(tab) {
      this.tab= tab;
      this.item= "";
      this.tabbable.updateTab(this);
    },
    onItem(item) {
      if (!this.item) {
        this.item= item; // selected a sub-group
      }
      this.tabbable.updateTab(this);
      this.$root.cmdSelected(item);
    },
    onNodeSelected(node, param, token) {
      // FIX: synchronize browser display
    },
  },
  mounted() {
    this.$root.$on("node-selected", this.onNodeSelected);
    this.handler.listen(this.$refs.browserList);
  },
  beforeDestroy() {
    this.$root.$off("node-selected", this.onNodeSelected);
    this.handler.silence();
  },
  mixins: [bemMixin()],
});


