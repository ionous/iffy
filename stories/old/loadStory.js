function getStory() {
  const serial= localStorage.getItem("save");
  console.log(serial);
  const ret= serial? JSON.parse(serial): newStory();
  const got= JSON.stringify(compact(ret), 0, 2);
  // pull some pairs of closing parens together, and opening brackets onto the same line
  const out= got.replace(/([}\]])\n\s+([}\]])/g, "$1$2").replace(/",\n\s+([{\[])/g, '", $1');
  console.log(out);
  return ret;
}
