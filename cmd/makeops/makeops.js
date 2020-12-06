// node makeops.js > temp.go
'use strict';

const Handlebars = require('handlebars'); // for templates
const fs = require('fs'); // filesystem for loading iffy language file
const vm = require('vm'); // virtual machine for parsing iffy language file
const Make = require('./directives.js'); // composer directives

// load the language file; brings 'localLang()' into global scope.
const filename= `../compose/www/data/lang/iffy.js`
vm.runInThisContext(fs.readFileSync(filename));
const m= new Make();
localLang(m);
var processing= "";

// change lower case to pascal-cased ( golang public )
const pascal= function(name) {
  if (name && name[0]=== '$') {
    name= name.slice(1).toLowerCase();
  }
  const els= name.split('_').map(el=> el.charAt(0).toUpperCase() + el.slice(1));
  return els.join('');
};

const lower= function(name) {
  return name.slice(1).toLowerCase();
}

// given a strType with the specified pascal'd name, return its list of choice
const strChoices= function(name, strType) {
  const out=[];
  if (!strType) {
    throw new Error(`${name} has no strType processing ${processing}`);
  }
  for (const k in strType.with.params) {
    const p=lower(k);
    if (p === name) {
      out.unshift(p); // move the dynamic key to the front
    } else {
      out.push(p);
    }
  }
  if (!out.length) {
    throw new Error("xxx");
  }
  return out;
};

Handlebars.registerHelper('Pascal', pascal);
Handlebars.registerHelper('Lower', lower);

// try to determine if the field name and field type have the same name
// ( and therefore the golang spec can be anonymous )
Handlebars.registerHelper('IsAnonymous', function(a,b) {
  return pascal(a) == pascal(b);
});

// does the passed string start with a $
Handlebars.registerHelper('IsToken', function(str) {
  return (str && str[0]=== '$');
});

// is the passed name a slot/interface
Handlebars.registerHelper('IsSlot', function(name) {
  return (name.indexOf("_eval") >= 0) || (m.types.all[name].uses=== 'slot');
});

// for uses='str'
Handlebars.registerHelper('IsDynamic', function(strType) {
  const name= lower(strType.name);
  const cs= strChoices(name, strType);
  return cs[0] === name;
});

// for uses='str'
Handlebars.registerHelper('Choices', function(strType) {
  const name= lower(strType.name);
  const cs= strChoices(name, strType);
  return cs[0]===name? cs.slice(1): cs; // remove the dynamic key
});

// for uses='str'
Handlebars.registerHelper('IsEnumerated', function(strType) {
  const name= lower(strType.name);
  const cs= strChoices(name, strType);
  return cs.length >1 || (cs[0] !== name);
});

// flatten desc
Handlebars.registerHelper('DescOf', function (x) {
  let ret='';
  if (x.desc) {
    const desc= x.desc;
    if (typeof desc == 'string') {
      ret= desc;
    } else if (desc) {
      ret= pascal(desc.label || x.name);
      const rest= ((desc.short || '') + ' '+ (desc.long || '')).trim();
      if (rest) {
        ret+= ': ' + rest;
      }
    }
  }
  return ret;
})

// flatten groups
Handlebars.registerHelper('GroupOf', function (desc) {
  return desc.group.join(', ');
})

// console.log(m.currGroups);
console.log("package lang\n")

// load each js file as a handlebars template
const partials= ['spec'];
const sources= ['slots', 'num', 'txt', 'opt', 'run', 'str'];
partials.forEach(k=> Handlebars.registerPartial(k, require(`./templates/${k}Partial.js`)));
const templates= Object.fromEntries(sources.map(k=> [k,
Handlebars.compile(require(`./templates/${k}Template.js`))]));

const sorted= Object.keys(m.types.all).sort();
// console.log(sorted);
// console.log(m.types.all);
// return;

const slots= [];
const groups= [];
for (const typeName of sorted) {
  processing= typeName;
  const type= m.types.all[typeName];
  const mytemp= templates[type.uses];
  if (mytemp) {
    const out= mytemp(type);
    console.log(out);
  } else {
    switch (type.uses) {
      case 'slot':
        slots.push(type);
      case 'group':
        groups.push(type);
      break;
      default:
        console.log('error: Unknown uses', type.uses);
      break;
    }
  }
}

console.log(templates.slots({slots}));
