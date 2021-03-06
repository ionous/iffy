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

func (*{{Pascal name}}) Choices() (choices map[string]string) {
  return map[string]string{
    {{#each (Choices @this)~}}"{{this.token}}": "{{this.value}}",{{#unless @last}} {{/unless}}{{/each}}
  }
}

{{>spec spec=this}}
`;
