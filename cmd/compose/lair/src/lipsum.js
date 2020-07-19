const lipsum= `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum nec lorem malesuada, condimentum nibh ac, viverra justo. Pellentesque eleifend lectus in quam rhoncus, a sodales orci rutrum. Donec eu nulla elementum, tincidunt nunc id, consequat eros. Cras laoreet facilisis neque id viverra. Vivamus a semper nisl. Nulla ultricies lectus sed rutrum pulvinar. Aliquam in diam efficitur est volutpat sollicitudin nec eu massa. Sed tempus, augue eget vehicula tristique, odio elit suscipit erat, sit amet congue est erat at mauris. Maecenas scelerisque dapibus metus, at pulvinar augue congue eu.`

const allWords = lipsum.split(' ');
let lastWord=0;

// a slice of words from the above lipsum string.
// id, words, text
class Lipsum {
  constructor(cnt) {
    const idx= lastWord;
    const words= cnt? Lipsum.words(idx, cnt): ["<blank>"];
    lastWord= (lastWord+(cnt||0)) % allWords.length;
    this.text= words.join(" ");
  }
  // generate an array of arrays of strings [[a,b,c],[d,e,f]]
  static list(...wordcounts) {
    let wordset= wordcounts.map((c) => {
      const lipsum= new Lipsum(c);
      const parent= new Item();
      parent.content= lipsum.text.split(". ").map((x) => {
        return new Item(parent, `${x.trim().replace(/\.|,$/,'')}.`);
      });
      return parent;
    });
    return wordset;
  }
  // return an array of
  static words(idx, cnt) {
    // when the end index is greater, that's fine.
    const out= allWords.slice(idx, idx+ cnt);
    const rem= cnt - out.length;
    if (rem > 0) {
      const rest= allWords.slice(0, rem);
      Array.prototype.push.apply(out, rest)
    }
    return out;
  }
};
