

class DragHelper {
  constructor() {
    // drag info
    this.startEl= false;
    // this.containerEl= false;
    // focus
    this.focusId= false;
    this.focusEl= false;
    // drawDecoring info
    this.decorEl= false;
    this.edgeIdx = -1;
    this.decorId= false;
    // edge detection
    this.lastX= 0;
    this.lastY= 0;
    // event function, fix: maybe there's a better way?
    // const self= this;
    // this.update= function(evt) {
    //   self._update(evt);
    // }
    // class data for the focused element
    this.focusClasses= ["em-drag-highlight", "em-drag-mark"];
    this.imageClasses= ["em-drag-image", "em-drag-mark"];
    this.edgeClasses= ["em-item--bot", "em-item--top"];
  }
  // id of drag item, el of drag item, event dataTransfer
  start(id, el, dt) {
    el.classList.add(...this.imageClasses);
    dt.setDragImage(el,10,10); // fix? maybe should be click relative?
    setTimeout(()=>{
      el.classList.remove(...this.imageClasses);
    });
    this.startEl= el;
  }
  // id of item entered, el of item entered.
  drag(id, el) {
    if (el !== this.focusEl) {
      this.focusId= id;
      this.focusEl= el;
    }
    // when we enter a valid item we have a top border, or no border.
    const edge= (el === this.startEl)? -1: 1;
    this.drawDecor(edge);
  }
  // id matches start.
  end() {
    this.clearDecor();
    this.setBounds(false);
    this.focusEl= false;
    this.focusId= false;
    this.startEl= false;
  }
  // happens after ".start"
  // setBounds(el) {
  //   this.containerEl= el;
  //   const call= el ? "addEventListener" : "removeEventListener";
  //   document[call]("drag", this.update, { passive: true});
  // }
  // fix? getComputedStyle requires unit conversion ( em, etc-> px )
  // getLineHeight() {
  //   return 20;
  // }
  // fix: firefox sets clientX/Y to zero during drag events
  // and mouse move is absorbed.... so...
  // _update(evt) {
  //   // if we're hovering over the start element:
  //   // dont show a injection border.
  //   // check for being below
  //   const container= this.containerEl;
  //   if (container && evt.clientX !== 0 && evt.clientY !== 0) {
  //     let edge= -1;
  //     const focus= this.focusEl;
  //     if (focus) {
  //       const bounds= container.getBoundingClientRect();
  //       if (evt.clientX > 0/*bounds.x*/) {
  //         // we prefer a top border
  //         if (focus !== this.startEl) {
  //           edge= 1;
  //         }
  //         const b= evt.clientY - bounds.bottom;
  //         if (b >=0) {
  //           const height= this.getLineHeight();
  //           edge= (b < height)? 0: -1; // top edge of the first element.
  //         }
  //       }
  //     }
  //     this.drawDecor(edge, true);
  //   }
  // }
  // 0:bottom, 1:top, -1:out of bounds
  drawDecor(newEdgeIdx, fromUpdate=false) {
    const newDecor= this.focusEl;
    const oldDecor= this.decorEl;
    const focusId= this.focusId;
    //
    if (oldDecor !== newDecor || (this.edgeIdx !== newEdgeIdx)) {
      // if anything's changed, remove it all.
      // ( we might add some of it back almost immediately. )
      this.clearDecor();

      if (newDecor) {
        if (newEdgeIdx >=0) {
          const name= this.edgeClasses[newEdgeIdx];
          newDecor.classList.add(name);
        }

        newDecor.classList.add(...this.focusClasses);
        console.log("decorated", focusId, newEdgeIdx, fromUpdate);

        this.decorEl= newDecor;
        this.decorId= focusId;
        this.edgeIdx= newEdgeIdx;
      }
    }
  }
  clearDecor() {
    const oldDecor= this.decorEl;
    if (oldDecor) {
        console.log("remove decor", this.decorId);
        oldDecor.classList.remove(...this.edgeClasses);
        oldDecor.classList.remove(...this.focusClasses);
        this.decorEl= false;
        this.decorId= false;
        this.edgeIdx= -1;
      }
  }
};
