// maps typeName to array of names
class Tab extends Map {
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
    return ret.sort((a,b)=>  a.name.localeCompare(b.name));
  }
};
Vue.component('mk-browser', {
  template:
  `<div :class="cls.win"
  ><div class="mk-aux__title mk-aux__title--right"
     >commands</div
  ><div :class="bemElem('nav')"
    ><span v-for="(x,i) in tabNames"
      ><template v-if="i"
      >, </template
      ><mk-a-button
          :class="[ bemElem('btn'), tab===x?cls.btnSel:false ]"
          @activate="onFilter(x)"
      >{{x| capitalize}}</mk-a-button
    ></span
  >.</div
  ><div :class="bemElem('subtitle')"
     >{{item}} {{tab|capitalize}}</div
  ><ul
    class="mk-browser-list"
    ><li
      v-for="k in kids"
      ><span v-if="k.name">&#x2753;<mk-a-button
          @activate="onItem(k.name)"
        >{{k.label| titlecase}}</mk-a-button
      ></span
      ><span v-else>&nbsp;{{k.label| titlecase}}</span
    ></li
  ></ul
  ></div>`,
   computed: {
    // list of commands to display
    kids() {
      return this.tabs[this.tab].contents(this.item);
    },
  },
  data(){
    const types= allTypes.all;
    // compile groups, slots, and strs
    const all= new Tab({sansContent:true});
    const groups= new Tab(); // group name => [ runs that implement the group ]
    const slots= new Tab();
    const str= new Tab({txtContent:true});
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
    const tabs= {
        all,
        groups,
        slots,
        str,
      };
    return {
      // all possible tabNames. displayed at top of ctrl.
      tabNames: Object.keys(tabs),
      // current tab. changed by the user.
      tab: "groups",
      // current item within a tab.
      item: "",
      // css class helper
      cls: {
        win: [ this.bemBlock(), 'mk-aux' ],
        btnSel: this.bemElem('btn', 'sel'),
      },
      tabs,
    };
  },
  methods: {
    onFilter(x) {
      this.tab= x;
      this.item= "";
    },
    onItem(k) {
      if (!this.item) {
        this.item= k; // selected a sub-group
      }
      this.$root.cmdSelected(k);
    },
    onNodeSelected(node, param, token) {
      // FIX: synchronize browser display
    },
  },
  mounted() {
    this.$root.$on("node-selected", this.onNodeSelected);
  },
  beforeDestroy() {
    this.$root.$off("node-selected", this.onNodeSelected);
  },
  mixins: [bemMixin()],
});


