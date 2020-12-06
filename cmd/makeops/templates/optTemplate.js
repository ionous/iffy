// optTemplate.js
'use strict';

// for now we'll implement slots as structs with pointers
module.exports =
`// {{Pascal name}} swaps between various options
type {{Pascal name}} struct {
  {{Pascal name}} interface{}
}

{{>spec spec=this}}

func (*{{Pascal name}}) Swap() map[string]interface{} {
  return map[string]interface{} {
{{#each with.params}}
    "{{Lower @key}}": (*{{Pascal type}})(nil),
{{/each}}
  }
}
`;

/*{
  name: 'noun_phrase',
  uses: 'opt',
  group: [],
  with: {
    tokens: [
      '$KIND_OF_NOUN',
      ', ',
      '$NOUN_TRAITS',
      ', or ',
      '$NOUN_RELATION' ],
    params: {
      '$KIND_OF_NOUN': {
          label: 'kind of noun',
          type: 'kind_of_noun'
      },
      '$NOUN_TRAITS': {
          label: 'noun traits',
          type: 'noun_traits'
      },
      '$NOUN_RELATION': {
        label: 'noun relation',
        type: 'noun_relation' }
    }
  }
}
*/
