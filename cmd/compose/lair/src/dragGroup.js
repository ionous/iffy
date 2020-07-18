
// implementation specific to em-table.
class DragGroup {
  constructor(list, dropper) {
    this.list= list;
    this.dropper= dropper;
  };
  dragOver(over,dt) {
    this.dropper.setTarget(this.list, over);
    dt.dropEffect= "copy";
  }
  drag() {
    this.dropper.updateTarget(this.list);
  }
  dragLeave(leave, dt) {
    this.dropper.leaving= this.list;
  }
  dragEnd() {
    this.dropper.reset(true);
  }
  dragStart(start, dt) {
    this.dropper.setSource(this.list, start);
    const tgt= this._getDragImage(start, dt);
    Dropper.setDragData(dt, tgt, this._serializeItem(start));
  }
  drop(drop, dt) {
    const dropGroup= this.list;
    const {idx:dropIdx}= drop;
    const {idx:dragIdx, list:dragGroup} = this.dropper.source;
    // add and remove can ( sometimes ) cause dragend not to fire.
    // fix? while moving items is quick and easy
    // technically, we should create new items here by serialization --
    // and wait to remove items in drag end.
    //
    let width= 1;
    if (dragGroup.inline) {
      width= Number.MAX_VALUE;
    }
    if (dropGroup === dragGroup) {
      dragGroup.move(dragIdx, dropIdx, width);
    } else {
      let rub= dragGroup.items.splice(dragIdx, width);
      const at= Math.min(Math.max(0,dropIdx+1), dropGroup.items.length);

      // if we are moving multiple items from an inline group to a block group
      const merge= (rub.length > 1 && dragGroup.inline && !dropGroup.inline);
      if (!merge) {
        dropGroup.items.splice(at,0,...rub);
      } else {
        const text= rub.map((x)=>x.text).join(" ");
        console.log("merging", text);
        const obj= { id: rub[0].id, text  };
        dropGroup.items.splice(at,0, obj);
      }
    }

    // clear b/c we dont always get dragEnd.
    this.dropper.reset(true);
  }
  // fix: this should be more ....
  _serializeItem(start) {
    const item= this.list.items[start.idx];
    return {
      'text/plain': item.text,
    };
  }
  _getDragImage(start, dt) {
    let tgt= start.el;
    // create a temporary set of elements for an image
    // the blur drag source style is left to the .highlight
    if (this.list.inline) {
      tgt = document.createElement("span");
      let sib= start.el;
      while (1) {
        const add = sib.cloneNode(true);
        tgt.appendChild(add);
        sib= sib.nextSibling;
        if (!sib || TargetFinder.getData(sib, "dragIdx") === undefined) {
          break;
        }
      }
    }
    return tgt;
  }
};
