Vue.component('mk-tools', {
  template:
  `<div class="mk-buttons-form">
      <button disabled>Play</button>
      <button :disabled="!allow.testing" @click="onTest">Test</button>
      <span v-if="msg">{{msg}}</span>
    </div>`,
  data() {
    return {
      msg: "",
      allow: {
        testing: true,
      },
    }
  },
  methods: {
    onTest() {
      // https://xhr.spec.whatwg.org/#events
      // load: any success
      // progress: etc.
      // timeout: only if the timeout is set
      // abort:  ex the client called XMLHttpRequest.abort().
      // loadend: any completion
      const xhr = new XMLHttpRequest();
      //
      xhr.addEventListener("loadend", () => {
        this.allow.testing= true;
      });
      xhr.addEventListener("load", (evt) => {
        this.msg ="";
      });
      xhr.addEventListener("error", (evt) => {
        this.msg= "An unknown error occurred.";
        console.log(xhr.statusText);
      });
      this.msg= "Connecting...";
      this.allow.testing= false;
      const { story } = app.$data;
      const serial = JSON.stringify(story.item, 0, 2);
      // console.log("testing", serial);
      xhr.open("PUT", "/story/test");
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(serial);
    },

  },
  mixins: [bemMixin()],
});


