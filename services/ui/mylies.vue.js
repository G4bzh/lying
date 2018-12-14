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

          <md-list-item to="mydashboard">
            <md-icon>dashboard</md-icon>
            <span class="md-list-item-text">Dashboard</span>
          </md-list-item>

          <md-list-item to="mydns">
            <md-icon>dns</md-icon>
            <span class="md-list-item-text">DNS Configuration</span>
          </md-list-item>

          <md-list-item to="mysettings">
            <md-icon>settings</md-icon>
            <span class="md-list-item-text">Settings</span>
          </md-list-item>

        </md-list>

      </md-app-drawer>

      <md-app-content>

        <router-view></router-view>

      </md-app-content>

    </md-app>

  </div>
  `,
  data: function() {
    return {}
  }
};
