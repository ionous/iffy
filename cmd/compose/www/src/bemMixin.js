function bemMixin(n=null) {
  return {
    methods: {
      bemBlock(mod=false) {
        const blockName= n|| this.$options._componentTag;
        const ar= [blockName];
        if (mod) {
          ar.push(blockName+"--"+mod);
        }
        return ar;
      },
      bemElem(el, mod=false) {
        const blockName= n|| this.$options._componentTag;
        const elName= blockName+ "__" + el;
        const ar= [elName];
        if (mod) {
          ar.push(elName+"--"+mod);
        }
        return ar;
      },
    }
  }
};