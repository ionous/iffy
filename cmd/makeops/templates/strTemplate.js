// strTemplate.js
'use strict';
// there's a few different strs:
// - fixed choices that could be ints
// ---- lets not use them
// - some fixed choices and some dynamic choices
// --- { tokens: [ '$NOUN_NAME' ],
// - a fully dynamic str ( ex. noun_name )
// --- { tokens: [ '$NOUN_NAME' ],
// IsDynamic
//
module.exports =
`type {{Pascal name}} string
{{#if (IsDynamic this)}}

func (*{{Pascal name}}) Dynamic() string {
  return "{{Lower name}}";
}
{{/if}}{{#if (IsEnumerated this)}}
func (*{{Pascal name}}) Choices() []string {
  return []string{
    {{#each (Choices this)}}"{{this}}", {{/each}}
  }
}
{{/if}}

{{>spec spec=this}}
`;

/*
type Something string

func (*Something) Dynamic() string {}
func (*Something) Choices() []string {
  return []string{
     "a", "b", "c",
  }
}
*/


/* input: {
  name: 'an',
  uses: 'str',
  group: [],
  with: {
    tokens: [ '$A', ' or ', '$AN' ],
    params: { '$A': 'a', '$AN': [Object] },
    spec: '{a} or {an}'
  }
 */
