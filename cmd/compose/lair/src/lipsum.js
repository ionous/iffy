const lipsum= `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum nec lorem malesuada, condimentum nibh ac, viverra justo. Pellentesque eleifend lectus in quam rhoncus, a sodales orci rutrum. Donec eu nulla elementum, tincidunt nunc id, consequat eros. Cras laoreet facilisis neque id viverra. Vivamus a semper nisl. Nulla ultricies lectus sed rutrum pulvinar. Aliquam in diam efficitur est volutpat sollicitudin nec eu massa. Sed tempus, augue eget vehicula tristique, odio elit suscipit erat, sit amet congue est erat at mauris. Maecenas scelerisque dapibus metus, at pulvinar augue congue eu.`
const allWords = lipsum.split(' ');
let lastWord=0;

// slice words from the above lipsum string.
// id, words, text
class Lipsum {
  // generate an array of items containing items with arrays of sentences.
  static list(...wordcounts) {
    let wordset= wordcounts.map((cnt) => {
      const parent= new Item();
      const lipsum= cnt? Lipsum.newSentences(cnt): [];
      parent.content= lipsum.map((sentence) => {
        return new Item(parent, sentence);
      });
      return parent;
    });
    return wordset;
  }
  // return an array of sentences.
  static newSentences(cnt) {
    const lipsum= Lipsum.newString(cnt);
    return lipsum.split(". ").map((x) => {
      const noTrailingPunct= x.trim().replace(/\.|,$/,'');
      return noTrailingPunct.charAt(0).toUpperCase()+
             noTrailingPunct.slice(1)+ '.';
    });
  }
  // return a string of cnt words, maybe a proper sentence. maybe not.
  static newString(cnt) {
    const idx= lastWord;
    const words= Lipsum.newWords(idx, cnt);
    lastWord= (lastWord+(cnt||0)) % allWords.length;
    return words.join(" ");
  }
  // return an array of words
  static newWords(idx, cnt) {
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
