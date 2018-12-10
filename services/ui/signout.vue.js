export default {
  name: "SignOut",
  template: `
  <div>
    <md-content class="md-elevation-3">

      <md-button class="md-raised md-primary" v-on:click="doSignOut">Confirm Sign Out</md-button>

    </md-content>

  </div>
  `,
  methods : {
    doSignOut: function() {
      localStorage.removeItem('user');
      this.$router.push("/");
    }
  }
};
