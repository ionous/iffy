//------------------------------------------------
// expands to a control appropriate for a node based on its type.
Vue.component('mk-switch', {
 props: ["node", "param", "token"], // node is actually a single node, or an array of nodes
 render(createElement) {
    var component;
    const { node, param, token } = this;
    if (token && token[0]!=='$') {
      component= "mk-plain-text";
    } else if (Array.isArray(node)) {
      // note: we could look at "param.repeats" but param is used for both the array slot
      // and the array elements
      component= "mk-repeater-ctrl";
    } else if (node) {
      const { itemType } = node;
      if (itemType) {
        // search for a template particular to the item's underlying type.
        component= `mk-${itemType.name}-ctrl`
        // if not, use a generic control based on item's role.
        if (!(component in Vue.options.components)) {
          component= `mk-${itemType.uses}-ctrl`;
        }
      }
    }
    return component && createElement( component, {
      key: node&& node.id, // can be empty for repeats...
      props: {
        node,
        param,
        token,
      }
    });
  }
});
