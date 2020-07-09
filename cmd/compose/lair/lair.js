const counts= [15, 30, 3, 5, 25, 7];
let allItems= counts.map((c) => new Lipsum(c));

const dragHelper= new DragHelper();


// use a pair of numbers for the gutter to manage the sizing.
Vue.component('em-gutter', {
  template:
  `<div class="em-gutter"
    ><div class="em-len"
    >{{max}}</div
    ><div class="em-num"
    >{{num}}</div
  ></div>`,
  props: {
    num: Number,
    max: Number,
  },
});

Vue.component('em-item', {
  template:
  `<div class="em-item"
      @dragstart="onDragStart($event)"
      @dragend="onDragEnd($event)"
      @dragenter="onDragOver($event, true)"
      @dragover="onDragOver($event)"
    ><em-gutter
      :num="num"
      :max="1234"
      draggable="true"
    ></em-gutter
    ><div
      class="em-content"
    >{{item.text}}</div
  ></div>`,
  props: {
    num: Number,
    item: Object,
  },
  methods: {
    logEvent(evt) {
      return;
      const item= this.item;
      const target= evt.target;
      const dt= evt.dataTransfer;
      console.log(evt.type, item.id, target.nodeName, dt.dropEffect);
    },
    // the event targets the em-gutter (draggable=true)
    // the user is attempting to drag,
    // and bubbles up here to the item.
    onDragStart(evt) {
      this.logEvent(evt);
      // prepare data
      const item= this.item;
      const text= item.text;
      const json= JSON.stringify({id:item.id, cnt: item.words.length});
      // set drag data
      const dt= evt.dataTransfer;
      dt.setData('text/plain', text);
      dt.setData('application/json', json);
      dt.effectAllowed= 'all';
      // set the drag image
      dragHelper.start(item.id, evt.currentTarget, dt);
      // propagate to parent
    },
    // moving over a drop target; it has bubbled up to the item.
    onDragOver(evt, enter) {
      const dt= evt.dataTransfer;
      if (dt) {
        dt.dropEffect= "copy";
      }
      dragHelper.drag(this.item.id, evt.currentTarget);
      this.logEvent(evt);
      // evt.stopPropagation();
      evt.preventDefault();
    },
    // sent to the target which triggered the drag
    // ( the em-gutter ) and bubbles up to this item handler.
    onDragEnd(evt) {
      this.logEvent(evt);
      evt.stopPropagation();
      evt.preventDefault();
      dragHelper.end();
    },

  }
});

Vue.component('em-table', {
  props: {
    items: Array,
  },
  template:
  `<div
      class="em-table"
      @dragstart="onDragStart($event)"
      @drop="onDrop($event)"

    ><transition-group
      name="flip-list"
      ><em-item
        v-for="(item,i) in items"
        :key="item.id"
        :item="item"
        :num="i*i*i"
      ></em-item
    ></transition-group
    ><div
      class="em-table__footer"
      @dragenter="onDragOver($event, true)"
      @dragover="onDragOver($event)"
    ></div
  ></div>`,
   methods: {
    onDragStart(evt) {
      // dragHelper.setBounds(evt.currentTarget);
      evt.stopPropagation();
    },
    // moving over a drop target; it has bubbled up to the item.
    onDragOver(evt, enter) {
      dragHelper.drawDecor(0);
      evt.stopPropagation();
    },
    // since this handler is on the table
    // we either get here by a drop on the item, or by bubble up from one of its elements.
    onDrop(evt) {
      var item; // find item we are dropping upon
      const tableEl= evt.currentTarget;
      for (let el= evt.target; !item && el !== tableEl; el= el.parentElement) {
        const vue= el.__vue__;
        item= vue && vue.item;
      }
      const items= this.items;
      const dstIdx= items.findIndex((i)=> i === item);
      // console.log("drop at", dstIdx);
      if (dstIdx >=0) {
        // note: dropEffect isnt registering as move here for some reason.
        const dt= evt.dataTransfer;
        /*if (dt.dropEffect==='move') */{
          // add the item to the list
          // find where we are dropping the item.
          // get the drag/drop data.
          const data= dt.getData('application/json');
          const json= JSON.parse(data);
          // find the original item
          const srcIdx= items.findIndex((item)=> item.id === json.id);

          // by default we are removing src, and inserting above dst
          // if dst is the next item, we'd stay in the same spot
          // so instead, we insert a blank line
          if (srcIdx+1 === dstIdx) {
            const blank= new Lipsum();
            items.splice(srcIdx, 0, blank);
          } else if (srcIdx< dstIdx) {
            const rub= items.splice(srcIdx,1);
            items.splice(dstIdx-1,0,rub[0]);
          } else if (srcIdx > dstIdx) {
            const rub= items.splice(srcIdx,1);
            items.splice(dstIdx,0,rub[0]);
          }
        }
      }
      evt.stopPropagation();
    },
  }
});

const app= new Vue({
  el: '#app',
  // methods: {},
  data: {
    items:allItems,
  },
});
