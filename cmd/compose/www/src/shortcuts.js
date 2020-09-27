class Shortcuts {
  constructor(redux) {
    const trap= new Mousetrap();
    this.mousetrap = trap;
    //
    var shifty= false;
    function changeShift(x) {
      if (x!== shifty) {
        app.shift= shifty= x;
      }
    }
    window.addEventListener("blur", (e) => {
      changeShift(false);
    });//
    const handleKey= trap.handleKey;
    trap.handleKey= function(ch,mod,e) {
      // note: https://jsfiddle.net/sqrbnaqw/
      // safari doesnt allow drag and drop to start when the shift key is held down.
      const shift= (mod.indexOf("shift")>=0) || (mod.indexOf("alt")>=0);
      if (e.type === "keydown") {
        changeShift(shifty || shift);
      }
      if (e.type === "keyup") {
        changeShift(shifty && shift);
      }
      handleKey.call(this, ch, mod, e);
    };
    // Mousetrap calls stopCallback to determine whether to stop keyboard events
    //  combo is ex. 'ctrl+shift+up'
    trap.stopCallback = function(e, el, combo) {
      // if (combo.indexOf("mod+") !== -1) {
      //   return false;
      // }
      // note: this is the original mousetrap callback:
      // if the element has the class "mousetrap" then allow the event
      if ((' ' + el.className + ' ').indexOf(' mousetrap ') > -1) {
        return false;
      }
      // stop mousetrap handlers for input, select, and textarea
      return el.tagName == 'INPUT' ||
            el.tagName == 'SELECT' ||
            el.tagName == 'TEXTAREA' ||
            (el.contentEditable && el.contentEditable == 'true');
    };
    trap.bind('mod+s', function(e) {
      const { story } = app;
      // its confusing to not have things save when you explicitly ask
      if (true /*redux.changed*/) {
        const serial = story.serialize();
        localStorage.setItem("save", serial);
        redux.changed = 0;
        console.log("saved", serial);
      }
      return false;
    });
    trap.bind('mod+z', function(e) {
      console.log("undo");
      redux.undo(true);
      console.log(redux.changed ? "needs save" : "up to date");
      return false;
    });
    trap.bind(['mod+y', 'mod+shift+z'], function(e) {
      console.log("redo");
      redux.redo(true);
      console.log(redux.changed ? "needs save" : "up to date");
      return false;
    });
    trap.bind('esc', function(e) {
      const { copier } = app;
      copier.cancel("escape");
      return false;
    });
    trap.bind('mod+c', function(e) {
      const { copier } = app;
      copier.start("shortcut");
      return false;
    });
  }
}
