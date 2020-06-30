
// global
class Sds {
  // returns: 'empty', 'slat, 'slot', 'leaf','prim', 'array'
  static isPrim(obj) {
    if (obj === undefined) {
      throw new Error("undefined");
    }
    return ['string', 'number', 'boolean'].indexOf(typeof obj) >=0;
  }
  static typeOf(obj) {
    var ret;
    if (Sds.isPrim(obj)) {
      ret= 'prim';
    } else if (obj===null) {
      ret= 'empty';
    } else if (Array.isArray(obj)) {
      ret= 'array';
    } else {
      const val= obj.value;
      if (Sds.isPrim(val) || Array.isArray(val) || val===null) {
        ret= 'leaf'; // a 'value' that's an array is a primitive array.
      } else if (!val.value) {
        ret= 'slat';
      } else {
        ret= 'slot';
      }
    }
    return ret;
  }
};

// self describing storage
function compact(obj) {
  const arrayType= function(ar) {
    const el= ar[0];
    return el.type || "";
  }
  const makeTag= function(field, obj) {
    // annotation for an array comes from the elements
    const antOf= obj.length>0 ? arrayType(obj): obj.type;

    // if field and type are the same name, elide the annotation
    const antStr= antOf  ? `${antOf}` : "";
    const fieldStr= ((!field) || (field===antOf)) ? "": field;
    const idStr= obj.id? `@${obj.id}`:"";
    return (antStr || idStr) ? `${fieldStr}::${antStr}${idStr}`: fieldStr;
  };
  const annotate= function(field, obj, val) {
    const tag= makeTag(field, obj);
    return [tag, val];
  };
  // obj is the id, type, value node containing the slat.
  const makeSlat= function(obj) {
    const out= {};
    const slat= obj.value;
    for (const field in slat) {
      if (field.startsWith("$")) {
        const n= field.replace("$", "").toLowerCase();
        const el= slat[field];
        const at= makeTag(n, el);
        out[at]= makeValue(el);
      }
    }
    return out;
  };
  const makeSlots= function(ar) {
    const out= [];
    ar.forEach((el) => {
      // each element is a slot, ex. "execute"
      // we want to get at the slat plugged into its value.
      const slat= el.value;
      const val= makeSlat(slat);
      const ant= annotate(el.id, slat, val);
      out.push( ...ant );
    });
    return out;
  };
  const makeLeaves= function(ar) {
    const out= [];
    ar.forEach((el) => {
      if (!el.id) {
        out.push( el.value );
      } else {
        const ant= annotate("", el, el.value);
        out.push( ...ant );
      }
    });
    return out;
  };
  const makeEls= function(ar) {
    const out= [];
    ar.forEach((el) => {
        const val= makeValue(el);
        if (!el.id) {
          out.push( el.value );
        } else {
          const ant= annotate("", el, val);
          out.push( ...ant );
        }
    });
    return out;
  };
  const valueMaker= {
    empty:(v) => null,
    prim: (v) => v,
    slat(obj) {
      return makeSlat(obj);
    },
    slot(obj) {
      const slat= obj.value;
      const val= makeSlat(slat);
      return annotate("", slat, val);
    },
    leaf(obj) {
      return obj.value;
    },
    array(ar){
      var ret;
      if (!ar.length) {
        ret= [];
      } else {
        const elType= Sds.typeOf(ar[0]);
        switch (elType) {
          case 'prim':
            ret= ar;
            break;
          case 'slot':
            ret= makeSlots(ar);
            break;
          case 'leaf':
            ret= makeLeaves(ar);
            break;
          default:
            ret= makeEls(ar);
            break;
        };
      }
      return ret;
    }
  };
  const makeValue= function(val) {
    const valType= Sds.typeOf(val);
    return valueMaker[valType](val);
  };
  //
  const val= makeSlat(obj);
  return annotate("", obj, val);
};
