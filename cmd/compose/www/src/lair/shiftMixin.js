// listen to shift key presses
function shiftMixin() {
  return {
    created() {
      document.addEventListener("keydown", (e) => {
        const shift= e.key === "Shift";
        this.$root.shift= this.$root.shift || shift;
        console.log("keydown", e.key, this.$root.shift);
      });
      document.addEventListener("keyup", (e) => {
        const shift= e.key === "Shift";
        this.$root.shift= this.$root.shift && !shift;
        console.log("keyup", e.key, this.$root.shift);
      });
      window.addEventListener("blur", (e) => {
        this.$root.shift= false;
      });
    },
    data: {
      shift: false,
    }
  };
};
//

