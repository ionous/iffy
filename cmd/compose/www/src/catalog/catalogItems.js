
class CatalogItem {
  constructor( name, path ) {
    this.name= name;
    this.path= path; // ex. curr/sub
  }
  get fullpath() {
    const { path, name } = this;
    return path? `${path}/${name}`: name;
  }
}

class CatalogFolder extends CatalogItem {
  constructor( name, path ) {
    super(name, path);
    this.contents= false;
  }
}

class CatalogFile extends CatalogItem {
  constructor( name, path ) {
    super(name, path);
    this.story= false;
  }
}

