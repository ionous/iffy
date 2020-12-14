// slots.js
'use strict';
module.exports =
`var Slots = []composer.Slot{
{{~#each allTypes}}{{#if (IsSlot name)}}
  {
    Name: "{{name}}",
    Type: (*{{Pascal name}})(nil),
    Desc: "{{DescOf this}}",
  }
{{~#unless @last}},{{/unless~}}
{{~/if}}{{/each}}
}

var Model = []composer.Slat{
{{~#each allTypes}}{{#if (IsSlat name)}}
  (*{{Pascal name}})(nil),
{{~/if}}{{/each}}
}`;

/* input: {
     name: 'story_statement',
     desc: 'Phrase',
     uses: 'slot',
     group: []
} */

/* output:
 var Slots = []composer.Slot{{
  Name: "comparator",
  Type: (*Comparator)(nil),
  Desc: "Comparison Types: Helper used when comparing two numbers, objects, pieces of text, etc.",
}, {
  Name: "assignment",
  Type: (*Assignment)(nil),
  Desc: "Assignments: Helper used when setting variables.",
}}

var Slats = []composer.Slat{
  (*Something)(nil),
}
*/
