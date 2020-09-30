
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

class MockCatalog extends Cataloger {
  constructor(nodes) {
    super();
    this.store= new CatalogStore(nodes);
  }
  // injects the real contents of the targeted file
  // into the passed CatalogFile object.
  loadFile(file) {
    const path= file.fullpath;
    let story= this.store.getFile(path);
    if (!story) {
      const pick= Types.slats("story_statement");
      const order= path.split("").reduce((a,v)=>(a+v.charCodeAt(0)), 0);
      // const type= pick[ Math.floor(Math.random() * pick.length) ];
      const type= pick[order%pick.length];
      const storyData= Types.createItem(type.name);
      story= this.store.loadFile(path, storyData);
    }
    file.story= story;
  }
  // injects the real contents of the targeted folder
  // into the passed CatalogFolder object.
  getFiles(folder) {
    let mockFolder= mockCatalog;
    const path= folder.fullpath;
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
    folder.contents= MockCatalog.buildContents(mockFolder, path);
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
        ret.push( new CatalogFile(item, mockPath) );
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
