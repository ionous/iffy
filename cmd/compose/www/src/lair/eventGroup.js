class EventGroup {
  constructor(el, that, events) {
    this.el= el;
    this.calls= Object.keys(events).map((type) => {
      const method= events[type];
      const call= function(evt) { return that[method](evt); }
      el.addEventListener(type, call);
      return {
        type, call,
      };
    });
  }
  silence() {
    const el= this.el;
    for (const n of this.calls) {
       el.removeEventListener(n.type, n.call);
    }
  }
};
