function getStory() {
  const serial= localStorage.getItem("save");
  console.log(serial);
  return serial? JSON.parse(serial): newStory();
}
