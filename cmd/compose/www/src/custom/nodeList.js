
class NodeList {
  // token should target an array of the passed type
  constructor( redux, node, token, type ) {
    this.redux= redux;
    this.node= node;
    this.token= token;
    this.type= type
    this.items= node.getKid(token);
    this.inline= false;
  }
  // users should generally call "addBlank"
  makeBlank() {
    return this.redux.nodes.newFromType(this.type);
  }
  newAt(idx, cmd) {
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
  // at:index, from:Draggable
  transferTo(toIdx, fromList, fromIdx, width=1) {
    const toList= this;
    //
    if (toList === fromList) {
      const needBlank= Math.abs(fromIdx - toIdx) === 1;
      if (needBlank) {
        this.addBlank(toIdx);
      } else {
        this.move(toIdx, fromIdx, width);
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
  addBlank(at=-1, newItem=null) {
    const { redux, node, items } = this;
    const blank= newItem || this.makeBlank();
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
