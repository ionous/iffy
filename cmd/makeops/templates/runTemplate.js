// runTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} requires various parameters.
type {{Pascal name}} struct {
{{#each with.params}}
  {{#unless (IsAnonymous @key type)}}
  {{Pascal @key}}
  {{~/unless~}}
  {{~#unless (IsSlot type)}} *{{~/unless~}}
  {{Pascal type}}
{{/each}}
}
{{#if with.slots}}

{{#each with.slots}}
var _ {{Pascal this}} = (*{{Pascal ../name}})(nil)
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
