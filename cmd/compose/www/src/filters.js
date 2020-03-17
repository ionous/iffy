const Filters= {
   capitalize(value) {
    let ret= '';
    if (value) {
      const str = value.toString().split("_").join(" ");
      ret= str.charAt(0).toUpperCase() + str.slice(1);
    }
    return ret;
  },
  titlecase(value) {
    let ret= '';
    if (value) {
      const parts= Filters.capitalize(value).split(" ");
      ret= parts.map((str)=>  str.charAt(0).toUpperCase() + str.slice(1)).
                 join(" ");
    }
    return ret;
  },
  none(value) {
    return value;
  }
};
for (const k in Filters) {
  Vue.filter(k, Filters[k]);
}
