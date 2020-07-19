let lastItem=0;
class Item {
  constructor(parent, content) {
    this.id= `id-${lastItem++}`;
    this.parent= parent;
    this.content= content;
  }
}
