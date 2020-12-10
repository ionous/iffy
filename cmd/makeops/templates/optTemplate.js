// optTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} swaps between various options
type {{Pascal name}} struct {
  At  reader.Position \`if:"internal"\`
  Opt interface{}
}

{{>spec spec=this}}

func (*{{Pascal name}}) Choices() map[string]interface{} {
  return map[string]interface{} {
{{#each with.params}}
    "{{Lower @key}}": (*{{Pascal type}})(nil),
{{/each}}
  }
}
`;
