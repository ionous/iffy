// strTemplate.js
'use strict';
// there's a few different strs:
// - fixed choices that could be ints
// ---- lets not use them
// - some fixed choices and some dynamic choices
// --- { tokens: [ '$NOUN_NAME' ],
// - a fully dynamic str ( ex. noun_name )
// --- { tokens: [ '$NOUN_NAME' ],
// IsClosed
//
module.exports =
`// {{Pascal name}} requires a user-specified string.
type {{Pascal name}} string

func (*{{Pascal name}}) Str() (closed bool, choices []string) {
  return {{#if (IsClosed this)}}true{{else}}false{{/if}}, []string{
    {{#each (Choices @this)~}}"{{this}}",{{/each}}
  }
}

{{>spec spec=this}}
`;

/*
type Something string

func (*Something) Choices() (closed bool, choices[] string {}) {
  return true, []string{
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
