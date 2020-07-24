//
function shiftMixin(n=null) {
  return {
    created() {
      document.addEventListener("keydown", (e) => {
        const shift= e.key === "Shift";
        this.shift= this.shift || shift;
        console.log("keydown", e.key, shift);
      });
      document.addEventListener("keyup", (e) => {
        const shift= e.key === "Shift";
        this.shift= this.shift && !shift;
        console.log("keyup", e.key, shift);
      });
      window.addEventListener("blur", (e) => {
        this.shift= false;
      });
    },
    data: {
      shift: false,
    }
  };
};
//

