

class DragHelper {
  constructor() {
    this.count = {
      start:0,
      focus:0, // needed for focus management
      enter:0, // helpful for testing...
      leave:0,
      clear:0,
      end:0
    };
    this.startEl= false;
    this.focusEl= false;
    this.containerEl= false;
    // edge detection
    this.edgeEl= false;
    this.edgeIdx = -1;
    this.edgeDetect= false;
    this.height= 20; // getComputedStyle requires unit conversion ( em, etc-> px )
    const self= this;
    this.update= function(evt) {
      self._update(evt);
    }
    // class data for the focused element
    this.focusClasses= ["em-drag-highlight", "em-drag-mark"];
    this.imageClasses= ["em-drag-image", "em-drag-mark"];
    this.edgeClasses= ["em-item--bot", "em-item--top"];
  }
  // el, dataTransfer
  start(el, dt) {
    ++this.count.start;
    el.classList.add(...this.imageClasses);
    dt.setDragImage(el,10,10); // fix? maybe should be click relative?
    setTimeout(()=>{
      el.classList.remove(...this.imageClasses);
    });
    this.startEl= el;
  }
  end() {
    ++this.count.end;
    this.clearFocus();
    this.setEdge(false);
    this.setBounds(false);
    this.startEl= false;
  }
  // happens after ".start"
  setBounds(el) {
    this.containerEl= el;
    const call= el ? "addEventListener" : "removeEventListener";
    document[call]("drag", this.update, { passive: true});
  }
  _update(evt) {
    // if we're hovering over the start element:
    // dont show a injection border.
    // check for being below
    const container= this.containerEl;
    if (container) {
      let edge= -1;
      const focus= this.focusEl;
      if (focus) {
        const bounds= container.getBoundingClientRect();
        if (evt.clientX > 0/*bounds.x*/) {
          if (focus !== this.startEl) {
            edge= 1;
          }
          const b= evt.clientY - bounds.bottom;
          if (b >=0) {
            edge= (b < this.height)? 0: -1; // top edge of the first element.
          }
        }
      }
      this.setEdge(focus, edge);
    }
  }
  setEdge(newEl, newEdge) {
    const lastEl= this.edgeEl
    if (lastEl !== newEl || (this.edgeIdx !== newEdge)) {
      if (lastEl) {
        lastEl.classList.remove(...this.edgeClasses);
      }
      if (!newEl || newEdge < 0) {
        this.clearEdge();
      } else {
        const name= this.edgeClasses[newEdge];
        newEl.classList.add(name);
        this.edgeEl= newEl;
        this.edgeIdx= newEdge;
        console.log("edge", this.edgeIdx);
        if (this.focusEl) {
          this.focusEl.classList.add(...this.focusClasses);
        }
      }
    }
  }
  clearEdge() {
    if (this.edgeEl !== false && this.edgeIdx !== -1) {
      this.edgeEl= false;
      this.edgeIdx= -1;
      // out of bounds, clear the highlight too
      if (this.focusEl) {
        this.focusEl.classList.remove(...this.focusClasses);
      }
      console.log("edge cleared");
    }
  }
  // forcibly clearFocus the drag highlights
  clearFocus() {
    ++this.count.clear;
    if (this.focusEl) {
      this.focusEl.classList.remove(...this.focusClasses);
      this.focusEl= false;
      this.count.focus=0;
    }
  }
  enter(el) {
    ++this.count.enter;
    if (el !== this.focusEl) {
      console.log("enter");
      this.clearFocus();
      const atStart= (el === this.startEl);
      this.setEdge(el, atStart? -1: 1 );
      el.classList.add(...this.focusClasses);
      this.focusEl= el;
    }
  }
  leave(el) {
    ++this.count.leave;
    if (el === this.focusEl) {
      --this.count.focus;
    }
  }
};
