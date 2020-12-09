// strTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} requires a user-specified string.
type {{Pascal name}} struct {
  At reader.Position \`if:"internal"\`
  Str string
}

func (op *{{Pascal name}}) String() string {
  return op.Str
}

func (*{{Pascal name}}) Choices() (closed bool, choices []string) {
  return {{#if (IsClosed this)}}true{{else}}false{{/if}}, []string{
    {{#each (Choices @this)~}}"{{this}}",{{/each}}
  }
}

{{>spec spec=this}}
`;
