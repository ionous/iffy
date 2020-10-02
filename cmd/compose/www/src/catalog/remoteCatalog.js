class RemoteCatalog extends Cataloger {
  constructor(nodes) {
    super();
    this.store= new CatalogStore(nodes);
    this.base= "/stories/";
    this._saving= false;
  }
  get saving() {
    return this._saving;
  }
  saveStories() {
    if (!this._saving) {
      this._saving= true;
      const json= JSON.stringify(this.store);
      this._put("", json, (res)=>{
        console.log("SAVED:", res);
        this._saving= false;
      });
    }
  }
  // inject a directory listing into folder
  loadFolder(folder) {
    const { path } = folder;
    this._get(path, (contents)=>{
      if (Array.isArray(contents)) {
        folder.contents= this.readContents(path, contents);
      }
    });
  }
  // inject a story file into folder.
  loadStory(file) {
    const { path } = file;
    let story= this.store.getStory(path);
    if (story) {
      file.story= story;
    } else {
      this._get(path, (storyData)=>{
        if (storyData) {
          story= this.store.storeStory(path, storyData);
          file.story= story;
        }
      });
    }
  }
  run(action, file, options, cb) {
    const { path } = file;
    this._send("POST", `${path}/${action}`, cb, options);
  }
  _get(path, cb) {
    this._send("GET", path, cb);
  }
  _put(path, body, cb) {
    this._send("PUT", path, cb, body);
  }
  _send(method, path, cb, body) {
    const url= this.base+path;
    console.log("xml http request:", method, url);
    var xhr = new XMLHttpRequest();
    xhr.addEventListener("load", ()=>{
      console.log("xml http response:", method, url, xhr.statusText);
      let data= true;
      if (xhr.response) {
        try {
          data= JSON.parse(xhr.response);
        } catch (e) {
          data= false;
        }
      }
      cb(data);
    });
    xhr.addEventListener("abort", ()=>cb(false));
    xhr.addEventListener("error", ()=>cb(false));
    xhr.open(method, url);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(body);
  }

  // ["/curr","/proj1","/proj2","/shared", "currStory.if"]
  readContents(path, got) {
    console.log("gotten:", got);
    return got.map((el)=> {
      const isFolder= el.startsWith("/");
      const name= el.slice(isFolder?1:0);
      return isFolder? new CatalogFolder(name, path): new CatalogFile(name, path);
    });
  }
};
