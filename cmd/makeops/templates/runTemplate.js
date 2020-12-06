// runTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} requires various parameters.
type {{Pascal name}} struct {
{{#each with.params}}
  {{Pascal @key}} {{{Lede this}}}{{Pascal type}}{{{Tail this}}}
{{/each}}
}
{{#if with.slots}}

{{#each with.slots}}
func (*{{Pascal ../name}}) {{Pascal this}}() {}
{{/each}}
{{/if}}

{{>spec spec=this}}
`;

/*input: {
name: 'pattern_variables_decl',
 desc:
  { label: 'Declare pattern variables',
    short: 'Storage for values used during the execution of a pattern.',
    long: '' },
 uses: 'run',
 group: [],
 with:
  { slots: [ 'story_statement' ],
    tokens:
     [ 'The pattern ',
       '$PATTERN_NAME',
       ' requires ',
       '$VARIABLE_DECL',
       '.' ],
    params:
     { '$PATTERN_NAME':
        { label: 'pattern name',
          type: 'pattern_name',
          filters: [ 'quote' ] },
       '$VARIABLE_DECL':
        { label: 'variable decl',
          type: 'variable_decl',
          repeats: true,
          filters: [ 'comma-and' ] } } } },
*/
