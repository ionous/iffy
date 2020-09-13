var _= function(){
  makeLang(new Make(new Types()));
  const all= {};
  const nodes= new Nodes( all, "tests" );
  const types= allTypes.all;
  for (const typeName in types) {
    const t= types[typeName];
    if (t.uses !== "group") {
      console.log("try", typeName);
      nodes.newFromType(typeName);
    }
  }
  console.log("done");

}();
