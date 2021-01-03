// node makeops.js > ../../ephemera/story/iffy_model.go
'use strict';

const Handlebars = require('handlebars'); // for templates
const allTypes= require('./model.js'); // iffy language file
// console.log(JSON.stringify(allTypes, 0,2 )); return;

// change to tokenized like name
const tokenize= function(name) {
  return '$'+ name.toUpperCase();
};

// change to lower case name
const lower= function(name) {
  if (name && name[0]=== '$') {
    name= name.slice(1);
  }
  return name.toLowerCase();
};

// change to pascal-cased ( golang public )
const pascal= function(name) {
  const els= lower(name).split('_').map(el=> el.charAt(0).toUpperCase() + el.slice(1));
  return els.join('');
};

const strChoices= function(token, strType) {
  const out=[];
  const { with : {params= {}, tokens=[]}= {} } = strType;
  return tokens.filter(t=>t[0]=='$' && t!==token).map((t,i)=> {
    const d= params[t]; // for str this should always be a string, but... directives.js
    return { label: d.label || d , index:i, token:t };
  });
};

const isClosed= function(token, strType) {
  const out=[];
  const { with : { tokens=[]}= {} } = strType;
  return tokens.indexOf(token) < 0;
};

Handlebars.registerHelper('Pascal', pascal);
Handlebars.registerHelper('Lower', lower);

// does the passed string start with a $
Handlebars.registerHelper('IsToken', function(str) {
  return (str && str[0]=== '$');
});

// characters preceding a type declaration
  // "label": "trait",
  // "type": "trait",
  // "optional": true,
  // "repeats": true,
  // "filters": [
  //   "comma-and"
  // ]
Handlebars.registerHelper('Lede', function(param) {
  let out = "";
  const name = param.type;
  const type = allTypes[name];
  if (param.optional) {
    out+= "*";
  }
  if (param.repeats) {
    out+= "[]";
  }
  out+= (name.indexOf("_eval") >= 0) ? "rt." :"";
  // out+= (name.indexOf("_eval") >= 0) ? "rt." :
  //       (type.uses !== 'slot')? "*": "";
  return out;
});

Handlebars.registerHelper('Tail', function(param) {
  return "";//param.optional? ' `if:"optional"`': "";
});

// is the passed name a slot
Handlebars.registerHelper('IsSlot', function(name) {
  const { uses }= allTypes[name];
  return uses === 'slot';
});

Handlebars.registerHelper('IsSlat', function(name) {
  const { uses }= allTypes[name];
  return uses !== 'slot' && uses !== 'group';
});

Handlebars.registerHelper('IsStr', function(name) {
  const { uses }= allTypes[name];
  return uses === 'str';
});

// for uses='str'
Handlebars.registerHelper('IsClosed', function(strType) {
  const token= tokenize(strType.name);
  return isClosed(token, strType);
});

// for uses='str'
Handlebars.registerHelper('Choices', function(strType) {
  const token= tokenize(strType.name);
  return strChoices(token, strType);
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

// load each js file as a handlebars template
const partials= ['spec'];
const sources= ['header', 'num', 'swap', 'flow', 'str', 'slot', 'footer'];
partials.forEach(k=> Handlebars.registerPartial(k, require(`./templates/${k}Partial.js`)));
const templates= Object.fromEntries(sources.map(k=> [k,
  Handlebars.compile(require(`./templates/${k}Template.js`))])
);
templates['txt']= templates['str']; // fix: txt really shouldnt even exist i think
console.log(templates.header({package:'story'}));

// switch to partials?
for (const typeName in allTypes) {
  const type= allTypes[typeName];
  const mytemp= templates[type.uses];
  if (mytemp) {
    console.log(mytemp(type));
  }
}

console.log(templates.footer({package:'story', allTypes}));
