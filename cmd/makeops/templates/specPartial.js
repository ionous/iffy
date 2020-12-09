// specPartial.js
'use strict';
module.exports =`func (*{{Pascal name}}) Compose() composer.Spec {
  return composer.Spec{
    Name: "{{name}}",
{{#if desc}}
    Desc: \`{{{DescOf this}}}\`,
{{/if}}
{{#if group}}
    Group: "{{GroupOf this}}",
{{/if}}
{{#if with.spec}}
    Spec: "{{{with.spec}}}",
{{/if}}
  }
}
`;

/* input: { name: '', dsec: {}, group: {} },
} */

/* output:
func (*Number) Compose() composer.Spec {
    return composer.Spec{
    Name:  "argument",
    Spec:  "its {name:variable_name} is {from:assignment}",
    Group: "patterns",
  }
}
*/
