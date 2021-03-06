// --------------------------------------------------------------------
// categorize a character from a format string
// https://github.com/ionous/makisu/wiki/FormatStrings
class Char {
  constructor(char) {
    this.ch= char;
  }
  opens() {
    return this.ch === "{";
  }
  closes() {
    return this.ch === "}";
  }
  ends() {
    return this.ch === "";
  }
  labels() {
    return this.ch === "%";
  }
  filter() {
    return this.ch === "|";
  }
  filterParam() {
    return this.ch === "=";
  }
  equals(other) {
    return this.ch === other;
  }
  occurs() {
    const occurs= {
      ':': {
        default: true,
      },
      '*': {
        optional: true,
        repeats: true
      },
      '+':{
        repeats: true
      },
      '?':{
        optional: true
      },
    };
    return occurs[this.ch];
  }
  optional() {
    return this.ch === '?' || this.ch === '*';
  }
  isSpecial() {
    return this.opens() || this.closes() || this.filter() ||
            this.labels() || this.occurs();
  }
};

// --------------------------------------------------------------------
// accumulate a chunk of plain text, or tag text from a format string
class Accum {
  constructor() {
    this.chars= [];
  }
  get length() {
    return this.chars.length;
  }
  push(char) {
    if (char.ends()) {
      throw new Error("unexpected end of string");
    }
    // accumulate all chars, including whitespace to optionally trim on flush
    this.chars.push(char.ch);
  }
  flush() {
    const text= this.chars.join("");
    this.chars= [];
    return text;
  }
};

// --------------------------------------------------------------------
// a tag includes arg, type, and label;
// as well as repeat/required settings.
class TagBlock {
  constructor(options) {
    this.block= {};
    this.occurs= null;
    this.options= options || {};
  }
  setOccurs(char) {
    let okay= false;
    const occurs= char.occurs();
    if (occurs) {
      if (!occurs.default) {
        if (this.options.asValues) {
          throw new Error("expected only one occurrence");
        }
        this.occurs= occurs;
      }
      okay= true;
    }
    return okay;
  }
  setArg(str) {
    this._set("arg", TagBlock.strip(str));
  }
  addFilter(str) {
    const filters= this.block['filters'] || [];
    this.block['filters']= filters.concat(TagBlock.strip(str));
  }
  addFilterVal(str, val) {
    const vals= this.block['filterVals'] || {};
    vals[str]= val;
    this.block['filterVals']= vals;
  }
  setLabel(str) {
    // out of this parse optional prefix and suffix
    const match= /(.*?)\[(.+?)\](.*)/.exec(str);
    if (!match) {
      this._set("label", TagBlock.strip(str, true));
    } else {
      const [_, prefix,label,suffix]= match;
      if (prefix) {
        this._set("prefix", TagBlock.strip(prefix, true, true));
      }
      this._set("label", TagBlock.strip(label, true));
      if (suffix) {
        this._set("suffix", TagBlock.strip(suffix, true, true));
      }
    }
  }
  setValue(key, str) {
    this._set(key, TagBlock.strip(str));
  }
  // set "reduce()" for block parameters:
  // target, label, prefix, filters, etc.
  _set(k, str) {
    if (k in this.block) {
      throw new Error(`duplicate key ${k}`);
    }
    this.block[k]= str;
  }
  static strip(str, keepInnerSpace=false, keepOuterSpace=false) {
    const v= keepOuterSpace? str: str.trim();
    if (!keepInnerSpace && v.match(/\s/)) {
      throw new Error(`format ${v} contained unexpected spaces`);
    }
    return v;
  }
  reduce() {
    // unpack to locals;
    // throws if arg is not set.
    const {
      block:{ arg, target=null, label=null, prefix=null, suffix=null, filters=[], filterVals=null },
      occurs,
      options: { asValues=false, nullValue=false }
    } = this;
    // clear old data.
    this.block= {};
    this.occurs= null;
    // target is going to exist, no label
    // otherwise target is
    const targetKey= asValues? "value": "type";
    const targetValue= TagBlock.getTargetValue(target, arg, nullValue);
    const args= (asValues && !label && !target && targetValue!==null)?
        // the asValue path allows dict to collapse to value
        targetValue:
        // rep only sets {repeats,optional} if one or both are true.
        Object.assign({label: label || arg.replace(/[-_]/g, ' ')},
                      {[targetKey]: targetValue},
                      occurs, // ex. repeat, optional
                      filters.length && {filters},
                      filterVals && {filterVals},
                      prefix && {prefix},
                      suffix && {suffix});
    return {
      arg: "$"+ arg.toUpperCase(),
      args: args
    };
  }
  // try to convert to a number otherwise leave as a string
  static getTargetValue(target, arg, nullValue) {
    let ret= null;
    const str= target || arg;
    if (str !== nullValue) {
      const n= Number(str);
      ret= Number.isNaN(n)? str:n;
    }
    return ret;
  }
};

// --------------------------------------------------------------------
// generator with with specs
class TagOutput {
  constructor() {
    this.keys= [];
    this.args= {};
  }
  writeText(str) {
    if (str) {
      this.keys.push(str);
    }
  }
  // finish the in progress TagBlock (msg)
  writeMsg(msg) {
    const block= msg.reduce();
    // de-dupe names
    let arg= block.arg;
    const og= arg;
    for (let i=1; arg in this.args; ++i) {
      arg= og + i.toString();
    }
    this.keys.push(arg);
    this.args[arg]= block.args;
  }
};

// --------------------------------------------------------------------
class TagParser {
  static parse(msg, options) {
    const p= new TagParser(options);
    if (msg) {
      for (let i=0; i< msg.length; ++i) {
        p.onChar(msg[i], msg[i+1]);
      }
    }
    return p.end("");
  }
  constructor(options) {
    this.accum= new Accum();
    this.msg= new TagBlock(options);
    this.escaping= false;
    this.state= this.readingText;
    this.out= new TagOutput();
  }
  end() {
    this.onChar("");
    return this.out;
  }

  onChar(c, next) {
    const char= new Char(c);
    this.state(new Char(c), next);
  }

  readingDone() {
    throw new Error("done");
  }

  // reading normal text, look for opening bracket
  readingText(char) {
      if (char.opens()) {
        this.out.writeText(this.accum.flush());
        this.state= this.readingFirst;
      } else if (char.ends()) {
        this.out.writeText(this.accum.flush());
        this.state= this.readingDone;
      } else {
        this.accum.push(char);
      }
   }
  // reading text directly after an open bracket
  readingFirst(char, next) {
      if (!this.escaping && char.isSpecial() && char.equals(next)) {
        this.escaping= true;
      } else if (this.escaping) {
        this.accum.push(char);
        this.escaping= false;
      } else if (char.labels()) {
        // turns out we were reading a label.
        // ex. {label%...}
        this.msg.setLabel(this.accum.flush());
        this.state= this.readingArg;
      } else if (this.msg.setOccurs(char)) {
        if (!this.accum.length) {
          // started with the occurrence operator?
          // ex. {#arg}
          this.state= this.readingTrailingArg;
        } else {
          // started with an arg.
          // {arg#type}
          this.msg.setArg(this.accum.flush());
          this.state= this.readingType;
        }
      } else {
        this.helpReadingTail(char);
      }
  }
  // after reading a label, we are reading an arg.
  readingArg(char) {
    if (this.msg.setOccurs(char)) {
      this.msg.setArg(this.accum.flush());
      this.state= this.readingType;
    } else {
      this.helpReadingTail(char);
    }
  }

  // explicit types end the block.
  readingType(char) {
    this.helpReadingTail(char, "target");
  }

  // we can have an arg closing the block if the occurrence operator appeared first.
  readingTrailingArg(char)  {
    this.helpReadingTail(char);
  }

  // read till the end of the tag, or till we encounter a filter
  // ex. ...}, or ...|
  helpReadingTail(char, key="arg") {
    if (char.closes()) {
      this.msg.setValue(key, this.accum.flush());
      this.out.writeMsg(this.msg);
      this.state= this.readingText;
    } else if (char.filter()) {
      this.msg.setValue(key, this.accum.flush());
      this.state= this.readingFilter;
    } else if (char.isSpecial()) {
      throw new Error("unexpected character");
    } else if (char.ends()) {
      throw new Error("unexpected end");
    } else {
      this.accum.push(char);
    }
  }

  // read till the end of the tag, accumulating filter text
  // ex. ...|filter|filter}
  readingFilter(char) {
    if (char.closes()) {
      this.msg.addFilter(this.accum.flush());
      this.out.writeMsg(this.msg);
      this.state= this.readingText;
    } else if (char.filter()) {
      this.msg.addFilter(this.accum.flush());
    } else if (char.filterParam()) {
      const filter= this.accum.flush();
      this.msg.addFilter(filter);
      this.state= (newChar) => {
        this.readingFilterVal(newChar, filter);
      };
    } else if (char.isSpecial()) {
      throw new Error("unexpected character");
    } else if (char.ends()) {
      throw new Error("unexpected end");
    } else {
      this.accum.push(char);
    }
  }

  // read till the end of the filter parameter
  // ex. ...|filter:5|filter}
  readingFilterVal(char, filter) {
    if (char.closes()) {
      this.msg.addFilterVal(filter, this.accum.flush());
      this.out.writeMsg(this.msg);
      this.state= this.readingText;
    } else if (char.filter()) {
      this.msg.addFilterVal(filter, this.accum.flush());
      this.state= readingFilter;
    } else if (char.isSpecial()) {
      throw new Error("unexpected character");
    } else if (char.ends()) {
      throw new Error("unexpected end");
    } else {
      this.accum.push(char);
    }
  }

};
