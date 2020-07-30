class Shortcuts {
  constructor(redux) {
    this.mousetrap = new Mousetrap();
    // where combo= 'ctrl+shift+up'
    this.mousetrap.stopCallback = function(e, el, combo) {
      if (combo.indexOf("mod+") !== -1) {
        return false;
      }
      // this is the original callback:
      // if the element has the class "mousetrap" then no need to stop
      if ((' ' + el.className + ' ').indexOf(' mousetrap ') > -1) {
        return false;
      }
      // stop for input, select, and textarea
      return el.tagName == 'INPUT' ||
        el.tagName == 'SELECT' ||
        el.tagName == 'TEXTAREA' ||
        (el.contentEditable && el.contentEditable == 'true');
    };
    this.mousetrap.bind('mod+s', function(e) {
      const { story } = app.$data;
      if (redux.changed) {
        const serial = story.serialize();
        localStorage.setItem("save", serial);
        redux.changed = 0;
        console.log("saved", serial);
      }
      return false;
    });
    this.mousetrap.bind('mod+z', function(e) {
      console.log("undo");
      redux.undo(true);
      console.log(redux.changed ? "needs save" : "up to date");
      return false;
    });
    this.mousetrap.bind(['mod+y', 'mod+shift+z'], function(e) {
      console.log("redo");
      redux.redo(true);
      console.log(redux.changed ? "needs save" : "up to date");
      return false;
    });
  }
}
