// numTemplate.js
'use strict';
module.exports =
`// {{Pascal name}} requires a user-specified number.
type {{Pascal name}} struct {
  At reader.Position \`if:"internal"\`
  Val float64
}

func (*{{Pascal name}}) Num()  []float64 {
    return []float64{
    {{#each (Choices @this)~}}"{{this}}",{{/each}}
  }
}

{{>spec spec=this}}
`;
