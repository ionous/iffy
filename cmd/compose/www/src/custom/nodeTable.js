// nodeList.js
class NodeTable extends DragList {
  constructor( redux, node, items ) {
    super(items);
    this.redux= redux;
    this.nodes= redux.nodes;
    this.node= node;
  }
  // users should generally call "addBlank"
  makeBlank() {
    throw new Error("not implemented");
  }
  // returns number of elements added
  addTo(at, elOrEls) {
    throw new Error("not implemented");
  }
  // returns the element or elements removed
  removeFrom(at, width) {
    throw new Error("not implemented");
  }
  // at:index, from:{list,idx}
  transferTo(at, fromGroup, fromIdx) {
    const toIdx= at;
    const toList= this;
    const fromList= fromGroup.list;

    if (toList === fromList) {
      const needBlank= Math.abs(fromIdx - toIdx) === 1;
      if (needBlank) {
        this.addBlank(toIdx);
      } else {
        this.move(toIdx, fromIdx, fromList.inline?Number.MAX_VALUE:1 );
      }
    } else {
      const { redux } = this;
      redux.doit({
        added: 0, // inelegant to say the least.
        apply() {
          const paraEls= fromList.removeFrom(fromIdx);
          this.added= toList.addTo( toIdx, paraEls );
        },
        revoke() {
          const paraEls= toList.removeFrom( toIdx, this.added );
          fromList.addTo( fromIdx, paraEls );
        },
      });
    }
  }
  addBlank(at=-1) {
    const { redux, node, items } = this;
    const blank= this.makeBlank();
    if (at<0) {
      at= items.length;
    }
    redux.doit({
      apply() {
        blank.parent= node;
        items.splice(at,0,blank);
      },
      revoke() {
        blank.parent= null;
        items.splice(at,1);
      },
    });
    return blank;
  }
  // move items within this same list
  move(src, dst, width, nothrow) {
    const { redux, items } = this;
    if ((!width) || (width<0)) {
      const e= new Error("invalid width");
      if (nothrow) { console.error(e); return e; }
      throw e;
    }
    if ((dst > src) && (dst < src+width)) {
      const e= new Error("invalid dest");
      if (nothrow) { return e; }
      throw e;
    }
    if (src+width> items.length) {
      width= items.length-src;
    }
    if (dst > src) {
      dst -= width;
    }
    redux.doit({
      apply() {
        const rub= items.splice(src, width);
        items.splice(dst, 0, ...rub);
      },
      revoke() {
        const rev= items.splice(dst, width);
        items.splice(src, 0, ...rev);
      }
    });
  }
}
