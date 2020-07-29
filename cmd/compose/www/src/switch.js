//------------------------------------------------
// expands to a control appropriate for a node based on its type.
Vue.component('mk-switch', {
  props: {
    node: Node,
 },
 render(createElement) {
    const defaultComponent= "span";
    var component;
    const { node } = this;
    if (!node) {
      component= defaultComponent;
    } else if (node.isArray) {
      component= "mk-repeater-ctrl";
    } else if (node.plainText) {
      component= "mk-plain-text";
    } else if (node.item) {
      const { itemType } = node;
      if (!itemType) {
        component= defaultComponent;
      } else {
        // search for a template particular to the item's underlying type.
        component= `mk-${itemType.name}-ctrl`
        // if not, use a generic control based on item's role.
        if (!(component in Vue.options.components)) {
          component= `mk-${itemType.uses}-ctrl`;
        }
      }
    }
    return createElement( component, {
      props: {
        node,
      }
    });
  }
});
