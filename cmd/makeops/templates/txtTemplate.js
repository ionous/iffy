// txtTemplate.js
'use strict';
module.exports =
`type {{Pascal name}} string
{{#if desc}}

{{>spec spec=this}}
{{/if}}
`;

/* input: { name: 'number', uses: 'num', group: [] },
}*/

/* output:
type Number float64
func (*Number) Compose() composer.Spec {
    return composer.Spec{
    Name:  "argument",
    Spec:  "its {name:variable_name} is {from:assignment}",
    Group: "patterns",
  }
}
*/
