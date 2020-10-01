class Cataloger {
  // saveStories(future: string, list, or none)
  // loadStory(fileItem) -> read file and inject into the passed object.
  // loadFolder(folderItem) -> read the directory and inject into the passed object.
}


// catalog store tracks individual files.
class CatalogStore {
  constructor(nodes) {
    this.nodes= nodes;
    this.cache= {};
  }
  getStory(path) {
    return this.cache[path];
  }
  storeStory(path, storyData) {
    const { nodes, cache } = this;
    const story= nodes.unroll(storyData);
    cache[path]= story;
    return story;
  }
  toJSON() {
    const { cache } = this;
    return Object.keys(cache).map(path=> {
      const story= cache[path];
      return {
        path,
        story,
      };
    });
  }
}
