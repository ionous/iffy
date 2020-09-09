// standalone drop handler
class TrashHandler  {
  constructor(target) {
    this.target= target;
  }
  dragOver() {
    return this.target;
  }
  // cant drag out of trash icon
  dragStart() {
    return false;
  }
  dragDrop(from) {
    const { list, target: { idx } } = from;
    Redux.Run({
      paraEls: false, // inelegant to say the least.
      apply() {
        this.paraEls= list.removeFrom(idx);
      },
      revoke() {
        list.addTo( idx, this.paraEls );
      },
    });
  }
};

// &#x1F4C2; -- open file folder
// &#x1F4C1; -- file folder
// &#x1F5D1; -- trashcan
// U+267B  -- filled recycling
// 267A -- thin recycling
// 2672 -- empty recycling
Vue.component('mk-trash-can', {
  template:
  `<div
    :class="bemBlock()"
    v-show="showing"
  ><span
    :class="bemElem('trash', hovering && 'over')"
    :data-drag-idx="-2"
  >&#x267A</span
  ></div>`,
  created() {
    this.trashTarget= new Draggable();
  },
  mounted() {
    const { "$root": root } = this;
    this.handler= new DragHandler(root.dropper, new TrashHandler(this.trashTarget)).
                  listen(this.$el);
  },
  beforeDestroy() {
    this.handler.silence();
    this.handler= null;
  },
  computed: {
    showing() {
      const { "$root": root } = this;
      const { start } = root.dropper;
      return start && (start instanceof DraggableNode);
    },
    hovering() {
      const { "$root": root } = this;
      return root.dropper.target === this.trashTarget;
    }
  },
  mixins: [bemMixin()],
});


