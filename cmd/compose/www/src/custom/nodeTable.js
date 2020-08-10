// nodeList.js
class NodeTable extends DragList {
  constructor( redux, node, items ) {
    super(items);
    this.redux= redux;
    this.nodes= redux.nodes;
    this.node= node;
  }
  makeBlank() {
    throw new Error("not implemented");
  }
  // at:index, from:{list,idx}
  transferTo(at, fromList, fromIdx) {
    const toIdx= at;
    const toList= this;

    if (toList === fromList) {
      const needBlank= Math.abs(fromIdx - toIdx) === 1;
      if (needBlank) {
        this.addBlank(toIdx);
      } else {
        this.move(toIdx, fromIdx);
      }
    } else {
      const { redux } = this;
      redux.invoke({
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
    redux.invoke({
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
    if (width<=0) {
      const e= new Error("invalid width");
      if (nothrow) { return e; }
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
    redux.invoke({
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
