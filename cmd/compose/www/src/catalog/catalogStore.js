class Cataloger {}

class CatalogStore {
  constructor(nodes) {
    this.nodes= nodes;
    this.cache= {};
  }
  getFile(path) {
    return this.cache[path];
  }
  loadFile(path, storyData) {
    const { nodes, cache } = this;
    const story= nodes.unroll(storyData);
    cache[path]= story;
    return story;
  }
}
