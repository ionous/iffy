
class TrashGroup {
  constructor(redux) {
    this.redux= redux;
  }
  dragOver() {
    return true;
  }
  // cant drag out of trash icon
  dragStart() {
    return false;
  }
  drop(from) {
    const { redux } = this;
    const { group, idx } = from;
    redux.invoke({
      paraEls: false, // inelegant to say the least.
      apply() {
        this.paraEls= group.list.removeFrom(idx);
      },
      revoke() {
        group.list.addTo( idx, this.paraEls );
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
//

Vue.component('mk-trash-can', {
  template:
  `<div
    :class="bemBlock()"
    v-show="dropper.dragging"
  ><span
    :class="bemElem('trash', hovering && 'over')"
    :data-drag-idx="-2"
  >&#x267A</span
  ></div>`,
  mounted() {
    this.handler.listen(this.$el);
  },
  beforeDestroy() {
    this.handler.silence();
  },
  computed: {
    hovering() {
      const at = this.dropper.target;
      const atList= at && (at.group === this.group);
      return atList;
    }
  },
  data() {
    const { "$root": root } = this;
    const dropper= root.dropper;
    const group= new TrashGroup(root.redux);
    const handler= new DragHandler(dropper, group);
    return {
      dropper,
      handler,
      group
    }
  },
  mixins: [bemMixin()],
});


