export default {
  name: "Dashboard",
  template: `
  <div>
    <md-content class="md-elevation-3">

      <div class="md-title">Dashboard </div>

      Hello {{ username }}

    </md-content>

  </div>
  `,
  data: function() {
    return {}
  },
  computed: {
    username: function () {
      const user = JSON.parse(localStorage.getItem('user'));
      if (user) {
        if (user.name) {
          return user.name;
        }
      }
      return "Unknown";
    }
  }
};
