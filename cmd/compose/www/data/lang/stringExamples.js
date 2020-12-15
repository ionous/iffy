// a sampling of string pickings:
function makeLang(make) {
  make.flow("story", "here's a test {string}");

  // a user editable string:
  // - issue: the popup is slightly edge too far to the left
  make.str("string");

  // an inline pick list
  // - issue: clicking a link popups up an edit box. ( it should just accept the text )
  make.str("string", "{usually}, {always}, {seldom}");

  // an inline pick with editing
  // - issue: deleting the editable text ( "" ) hides the text link permanently.
  make.str("string", "{usually}, {always}, {seldom}");

  // an inline edit box. the same as what make.str("string") generates
  make.str("string",  "{string}");

  // this is "good" for when you might need a constant in an option.
  make.str("string", "{usually}");
}
