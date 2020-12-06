// num.js
'use strict';

module.exports =
`// {{Pascal name}} requires a user-specified number.
type {{Pascal name}} float64

func (*{{Pascal name}}) Num() (closed bool, choices []float64) {
    return {{#if (IsClosed this)}}true{{else}}false{{/if}}, []float64{
    {{#each (Choices @this)~}}"{{this}}",{{/each}}
  }
}

{{>spec spec=this}}
`;

/* input: { name: 'number', uses: 'num', group: [] },
} */

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
