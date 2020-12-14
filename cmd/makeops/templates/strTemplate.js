// strTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} requires a user-specified string.
type {{Pascal name}} struct {
  At  reader.Position \`if:"internal"\`
  Str string
}

func (op *{{Pascal name}}) String() string {
  return op.Str
}

func (*{{Pascal name}}) Choices() (closed bool, choices map[string]string) {
  return {{#if (IsClosed this)}}true{{else}}false{{/if}}, map[string]string{
    {{#each (Choices @this)~}}"{{this}}": "{{Lower this}}",{{#unless @last}} {{/unless}}{{/each}}
  }
}

{{>spec spec=this}}
`;
