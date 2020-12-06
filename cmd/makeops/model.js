// allTypes.js
'use strict';
const fs = require('fs'); // filesystem for loading iffy language file
const vm = require('vm'); // virtual machine for parsing iffy language file
const Make = require('./directives.js'); // composer directives

// load the language file; brings 'localLang()' into global scope.
const filename= `../compose/www/data/lang/iffy.js`
vm.runInThisContext(fs.readFileSync(filename));
const m= new Make();
localLang(m);
//
const sorted = {};
Object.keys(m.types.all).sort().forEach((key) => {
  sorted[key] = m.types.all[key];
});

module.exports=sorted;
