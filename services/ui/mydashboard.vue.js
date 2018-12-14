export default {
  name: "MyDashboard",
  template: `
  <div>
    <md-card md-with-hover>

        <md-card-header>
          <div class="md-title">Hello</div>
        </md-card-header>

        <md-card-content>
          {{username}}
        </md-card-content>

    </md-card>

    <md-card md-with-hover>

        <md-card-header>
          <div class="md-title">Your Zones</div>
        </md-card-header>

        <md-card-content>
          10
        </md-card-content>

    </md-card>

    <md-card md-with-hover>

        <md-card-header>
          <div class="md-title">Requests</div>
        </md-card-header>

        <md-card-content>
          1234567
        </md-card-content>

    </md-card>

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
