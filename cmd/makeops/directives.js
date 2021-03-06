'use strict';

const TagParser = require('./tags.js');

class TypeSet {
  constructor() {
    this.all= {};
    this.slots= {}; // slot name => [ runs that implement the slot ]
    this.groups= {}; // group name => [ runs that implement the group ]
  }
  get(typeName) {
    return this.all[typeName];
  }
  has(typeName) {
    return !!this.get(typeName);
  }
  // implements but isnt the same type as...
  implements(doesThis, implementThat) {
    const t= this.get(doesThis);
    return t && (t.with && t.with.slots && t.with.slots.indexOf(implementThat)>=0);
  }
  // implements or is the stame type as...
  areCompatible(doesThis, implementThat) {
    const t= this.get(doesThis);
    return t && ((t.name === implementThat) ||
                ((t.with && t.with.slots && t.with.slots.indexOf(implementThat)>=0)));
  }
  // object { string name; string uses;
  //  union { string; object { label, short, long string } } desc;
  //  object with?; }
  newType(type) {
    const name= type.name;
    if (name in this.all) {
      throw new Error(`redefining type ${name}`);
    }
    this.all[ name ]= type;
    return type;
  }
  newItem(typeName, value) {
    if (!(typeName in this.all)) {
      throw new Error(`expected type, got '${typeName}'`);
    }
    return { type:typeName, value };
  }
}

module.exports = class Make {
  constructor() {
    this.types= new TypeSet();
    this.currGroups= [];
  }

  // introduce the passed group name to types
  // created during the passed function.
  group(name, ...descFn) {
    var desc, fn;
    const [a,b]= descFn;
    if (b===undefined) {
      fn= a;
      desc= "";
    } else {
      desc= a;
      fn= b;
    }
    const n= name.toLowerCase();
    this.currGroups.push(n);
    if (!this.types.has(n)) {
      this.newType( n, "group", desc );
    } else if (desc) {
      throw new Error(`group ${name} already declared`);
    }
    fn();
    this.currGroups.pop();
  }

  // typeName string, [ slot or slots string(s) ], msg string, desc multi-part string.
  // msg is a "format string" -- with token, types, etc.
  // desc is a description with the form: "label: short description. long description...".
  flow(name, ...slotMsgDesc) {
    const [b,c,d]= slotMsgDesc;
    var tags, slots, desc;
    // assume msg is the first parameter
    const firstTags= TagParser.parse(b);
    if (Object.keys(firstTags.args).length) {
      tags= firstTags;
      desc= c; // so desc is second
    } else {
      // for backwards compat with directiveTests; check that c had text before assigning tags
      const secondTags= TagParser.parse(c);
      if (!d && c && !Object.keys(secondTags.args).length) {
        tags= firstTags;
        desc= c;
      } else {
        const slot_or_slots= b;
        tags= secondTags;
        slots= Array.isArray(slot_or_slots)? slot_or_slots: [slot_or_slots],
        desc= d;
      }
    }
    return this.newType(name, "flow", desc,  {
        slots: slots,
        tokens: tags.keys,
        params: tags.args,
        spec: tags.msg,
    });
  }

  slot( name, desc= null ) {
    return this.newType(name, "slot", desc);
  }

  // displays types inline ( vs. slot and slat dropdowns )
  swap( name, msg, desc= null ) {
    const tags= TagParser.parse(msg);
    return this.newType(name, "swap", desc, {
        tokens: tags.keys,
        params: tags.args,
        spec: tags.msg,
    });
  }
  str( name, msg=null, desc= null ) {
    return this.makeStr(name, "str", msg, desc);
  }

  // pick or enter a small bit of text.
  makeStr( name, uses, msg, desc) {
    const settings={ asValues: true, nullValue: name };
    let tags= TagParser.parse(msg, settings);
    // msg had no tags: it's either a desc, or it was set/left as null.
    if (!Object.keys(tags.args).length) {
      desc= desc || msg;
      msg = `{${name}}`;
      tags= TagParser.parse(msg, settings);
    }
    return this.newType(name, uses, desc, {
        tokens: tags.keys,
        params: tags.args,
        spec: tags.msg,
     });
  }

  // multiline text
  txt( name, msg=null, desc= null ) {
    return this.makeStr(name, "txt", msg, desc);
  }

  num( name, desc= null ) {
    return this.newType(name, "num", desc);
  }

  newType(name, uses, desc, withspec=null) {
    const group= this.currGroups.length && this.currGroups;
    return this.types.newType(Object.assign(
      {name:name},
      desc&&{desc:Make.makeDesc(name, desc)},
      uses&&{uses},
      group&&{group},
      withspec&&{with:withspec}
    ));
  }

  // label, a human readable name
  // short, a short description
  // long, additional details
  static makeDesc(name, desc) {
    var ret;
    if (typeof desc !== 'string') {
      ret= desc;
    } else {
      let label= "";
      let short= "";
      let long= "";
      const i= desc.indexOf('.');
      if (i >= 0) {
        long= desc.substring(i+1).trimLeft();
        desc= desc.substring(0,i+1);
      }
      const j= desc.indexOf(':');
      if (j < 0) {
        short= desc;
      } else {
        label=  desc.substring(0,j);
        short= desc.substring(j+1).trimLeft();
      }
      ret= (label || long)? {
        label: (label || name),
        short,
        long
      }: short;
    }
    return ret;
  }

  // given a text description, add a new type
  // used by autogenerated/autogenerating types.
  newFromSpec(spec) {
    const d= Object.assign(spec,
      spec.desc&&{desc:Make.makeDesc(spec.name, spec.desc)},
    );
    if (d.spec) {
      const tags= TagParser.parse(d.spec);
      const w= d.with || {};
      w.tokens= tags.keys;
      w.params= tags.args;
      d["with"]= w;
    }
    this.types.newType(d);
  }
}
