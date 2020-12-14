// runTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} requires various parameters.
type {{Pascal name}} struct {
  At reader.Position \`if:"internal"\`
{{#each with.params}}
  {{Pascal @key}} {{{Lede this}}}{{Pascal type}}{{{Tail this}}}
{{/each}}
}

{{>spec spec=this}}
`;


// check that the struct implements its slots
// this doesnt work for stubs....
/*
{{#if with.slots}}

{{#each with.slots}}
var _ {{Pascal this}}= (*{{Pascal ../name}})(nil)
{{/each}}
{{/if}}
*/
