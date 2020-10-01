
const mockCatalog= [{
  "curr": [
    "currStory.if", {
      "alt": [
        "fileA.if",
        "fileB.if"
      ],
      "sub": [
        "fileC.if",
        "fileD.if",
        "fileE.if"
      ],
    },
  ],
  "proj": [
    "fileA.if",
    "fileB.if",
    "fileC.if",
    "fileD.if",
    "fileE.if"
  ],
  "empty": []
}];

class MockCatalog extends Cataloger {
  constructor(nodes) {
    super();
    this.store= new CatalogStore(nodes);
  }
  // maybe takes a string, list, or none
  // if none saves all.
  saveStories() {
    const json= JSON.stringify(this.store, 0,2);
    console.log("SAVED:", json);
  }
  // injects the real contents of the targeted file
  // into the passed CatalogFile object.
  loadStory(file) {
    const path= file.path;
    let story= this.store.getStory(path);
    if (!story) {
      const pick= Types.slats("story_statement");
      const order= path.split("").reduce((a,v)=>(a+v.charCodeAt(0)), 0);
      // const type= pick[ Math.floor(Math.random() * pick.length) ];
      const type= pick[order%pick.length];
      const storyData= Types.createItem(type.name);
      story= this.store.storeStory(path, storyData);
    }
    file.story= story;
  }
  // injects the real contents of the targeted folder
  // into the passed CatalogFolder object.
  loadFolder(folder) {
    let mockFolder= mockCatalog;
    const path= folder.path;
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
