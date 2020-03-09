function getStory() {
  const serial= localStorage.getItem("save");
  return serial? JSON.parse(serial): newStory();
}
