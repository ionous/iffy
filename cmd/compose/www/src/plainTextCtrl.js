Vue.component('mk-plain-text', {
   template:
    `<span class="mk-plain-text"
    >{{node.plainText}} </span>`,
  props: {
    node: Node,
  },
});
