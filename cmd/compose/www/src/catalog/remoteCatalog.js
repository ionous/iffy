class RemoteCatalog extends Cataloger {
  constructor(nodes) {
    super();
    this.store= new CatalogStore(nodes);
  }
  // inject into folder
  getFiles(folder, parts=[]) {
      const path= "/files/" + parts.join("/");
      var xhr = new XMLHttpRequest();
      xhr.addEventListener("load", ()=>{
        const got= JSON.parse(xhr.response);
        if (Array.isArray(got)) {
          folder.contents= this.readContents(got);
        }
      });
      console.log("GET", path);
      xhr.open("GET", path);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(path);
  }
  readContents(path, got) {
    return got.map((el)=> {
      const isFolder= el.startsWith("/");
      const name= el.slice(isFolder?1:0);
      return isFolder? new CatalogFolder(name, path): new CatalogFile(name, path);
    });
  }
};
