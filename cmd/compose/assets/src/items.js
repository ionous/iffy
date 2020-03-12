// items wrap values to make even string and number primitives unique.
// items exist for every required value, even if the value itself is empty.
// currently, ids are serialized. that's not the long term plan.
class Ids {
  constructor() {
    this.counter=0;
    this.base= "id-" + Date.now().toString(16) + "-";
  }
  nextId() {
    return this.base + this.counter++;
  }
}

// inner class for Types
// while a single global Type class simplifies code, it hurts testing.
// this provides a way to have a mixture of both.
class TypeSet {
  constructor() {
    this.all= {};
    this.slots= {}; // slot name => [ runs that implement the slot ]
    this.groups= {}; // group name => [ runs that implement the group ]
    this.ids= new Ids();
  }

  get(typeName) {
    return this.all[typeName];
  }
  has(typeName) {
    return !!this.get(typeName);
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
    return { id:this.ids.nextId(), type:typeName, value };
  }
}

let allTypes; // TypeSet singleton, contained/reset by Types.

// global
class Types {
  constructor() {
    this.allTypes= allTypes= new TypeSet();
  }

  // register a new named type
  static newType(spec) {
    return allTypes.newType(spec);
  }

  // spec for the named type
  static get(typeName) {
    return allTypes.all[typeName];
  }

  // types compatible with the named slot
  static slats(slotTypeName) {
    const slats= [];
    const all= allTypes.all;
    for (const typeName in all) {
      const type= all[typeName];
      if (type.uses === 'run') {
        const slots= type.with.slots;
        if (slots && slots.includes(slotTypeName)) {
          slats.push(type);
        }
      }
    }
    return slats;
  }

  static labelOf(t) {
    const label= !t.desc ? t.name.replace('_', ' ')  /* friendlyish name */ :
          (typeof t.desc === 'string') ? t.desc :
          t.desc.label;
    return label;
  }
  static shortOf(t) {
      const short= t.desc && t.desc.short ? t.desc.short : "";
      return short;
  }
  static longOf(t) {
    const long= t.desc && t.desc.long ? t.desc.long : "";
    return long;
  }
  static groupsOf(t) {
    const g= t.group;
    return g? (Array.isArray(g)? g: [g]) : [];
  }

  // produce an item/tree for the named type;
  // filling out defaults for all required fields.
  static createItem(typeName, ctx=null) {
    console.log("Types.createItem", typeName);
    if (typeof typeName !== 'string') {
      throw new Error("expected type string");
    }
    var ret;
    const type= allTypes.all[ typeName ];
    const { uses } =  type;
    switch (uses) {
      case "run": {
        const data= {};
        const { params } = type.with;
        for ( const token in params ) {
          const param= params[token];
          if (!param.optional) {
            const val= Types.createItem( param.type, {
              token: token,
              param: param
            });
            // if the param repeats then we'll wind up with an array (of items)
            data[token]= param.repeats? (val? [val]: []): val;
          }
        }
        ret= allTypes.newItem(type.name, data);
      }
      break;
      case "slot":
      case "opt": {
        // note: "initially", if any, is: object { string type; object value; }
        const pair= Types._unpack(ctx);
        if (!pair) {
          // ex. Item("story_statement", null)
          ret= allTypes.newItem(type.name, null);
        } else {
          const { type:slatType, value:slatValue } = pair;
          ret= Types.createItem(slatType, slatValue);
        }
      }
      break;
      case "str":
      case "txt": {
        // ex. Item("attribute", "testing")
        const value= Types._unpack(ctx, "");
        ret= allTypes.newItem(type.name, value);
      }
      break;
      case "num": {
        const value= Types._unpack(ctx, 0);
        ret= allTypes.newItem(type.name, value);
      }
      break;
      default:
        throw new Error(`unknown type ${uses}`);
      break;
    }
    return ret;
  }

  static _unpack( ctx, defaults ) {
    var ret;
    if (ctx && ctx.param && ("initially" in ctx.param)) {
      ret= ctx.param.initially;
    } else {
      ret= defaults;
    }
    return ret;
  }
}
