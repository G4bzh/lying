export default {
  name: "MyLies",
  template: `
  <div>
  <md-app>

    <md-app-toolbar class="md-primary">
      <span class="md-title">My Lies</span>
    </md-app-toolbar>

    <md-app-drawer md-permanent="card">

      <md-list>

        <md-list-item>
          <md-icon>dashboard</md-icon>
          <span class="md-list-item-text">Dashboard</span>
        </md-list-item>

        <md-list-item>
          <md-icon>dns</md-icon>
          <span class="md-list-item-text">DNS Configuration</span>
        </md-list-item>

        <md-list-item>
          <md-icon>settings</md-icon>
          <span class="md-list-item-text">Settings</span>
        </md-list-item>

      </md-list>

    </md-app-drawer>

    <md-app-content>
      Hello {{username}}
    </md-app-content>
    </md-app>

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
