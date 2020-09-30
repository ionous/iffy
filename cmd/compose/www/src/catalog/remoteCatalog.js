// new files?

const mockCatalog= [{
  "curr": [
    "currStory.if.js", {
      "alt": [
        "fileA.if.js",
        "fileB.if.js"
      ],
      "sub": [
        "fileC.if.js",
        "fileD.if.js",
        "fileE.if.js"
      ],
    },
  ],
  "proj": [
    "fileA.if.js",
    "fileB.if.js",
    "fileC.if.js",
    "fileD.if.js",
    "fileE.if.js"
  ],
  "empty": []
}];

class CatalogItem {
  constructor( name, path ) {
    this.name= name;
    this.path= path; // ex. curr/sub
  }
  get fullpath() {
    const { path, name } = this;
    return path? `${path}/${name}`: name;
  }
};

class CatalogFolder extends CatalogItem {
  constructor( name, path ) {
    super(name, path);
    this.contents= false;
  }
}

class MockCatalog {
  getFiles(model) {
    let mockFolder= mockCatalog;
    const path= model.fullpath;
    if (path.length) {
      // advance through the mockup data
      path.split("/").forEach((part)=>{
        // the last element is where the directories live.
        const last= mockFolder[mockFolder.length-1];
        const dirs= (typeof last !== 'string')? last: {};
        if (!part in dirs) {
          throw new Error(`unknown part ${part} in ${dirs}`);
        }
        mockFolder= dirs[part];
      });
    }
    model.contents= MockCatalog.buildContents(mockFolder, path);
  }
  static chunk(path) {
    const i= path.indexOf("/");
    return (i<0)? path: path.slice(0,i);
  }
  static buildContents(mockFolder, mockPath) {
    const ret= [];
    // walk the items in the folder:
    for (const item of mockFolder) {
      // file vs sub-folder map
      if (typeof item === 'string') {
        ret.push( new CatalogItem(item, mockPath) );
      } else {
        // note: if there's a map, its always the last item of the folder
        Object.keys(item).forEach((mockDir) => {
          ret.push( new CatalogFolder(mockDir, mockPath) );
        });
      }
    }
    return ret;
  }
};

class RemoteCatalog {
  // inject into model
  getFiles(model, parts=[]) {
      const path= "/files/" + parts.join("/");
      var xhr = new XMLHttpRequest();
      xhr.addEventListener("load", ()=>{
        const got= JSON.parse(xhr.response);
        if (Array.isArray(got)) {
          model.contents= this.readContents(got);
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
      return isFolder? new CatalogFolder(name, path): new CatalogItem(name, path);
    });
  }
};
