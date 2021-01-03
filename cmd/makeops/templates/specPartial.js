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
{{#if (IsStr name)}}
{{#unless (IsClosed this)}}
    OpenStrings: true,
{{/unless}}
{{#if (Choices this)}}
    Strings: []string{
      {{#each (Choices @this)~}}"{{this.value}}",{{#unless @last}} {{/unless}}{{/each}}
    },
{{/if}}
{{/if}}
  }
}
`;
