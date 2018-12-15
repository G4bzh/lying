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

          <md-list-item to="/mylies/mydashboard">
            <md-icon>dashboard</md-icon>
            <span class="md-list-item-text">Dashboard</span>
          </md-list-item>

          <md-list-item to="/mylies/mydns" md-expand>
            <md-icon>dns</md-icon>
            <span class="md-list-item-text">DNS Configuration</span>
            <md-list slot="md-expand">
              <md-list-item class="md-inset" to="/mylies/mydns/config">Genral Configuration</md-list-item>
              <md-list-item class="md-inset" v-for="zone in zones" :key="zone.domain" :to="'/mylies/mydns/zone/' + zone.domain">{{zone.domain}}</md-list-item>
            </md-list>
          </md-list-item>

          <md-list-item to="/mylies/mysettings">
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
    return {
      zones: []
    }
  },
  computed: {
    token: function () {
      const user = JSON.parse(localStorage.getItem('user'));
      if (user) {
        if (user.token) {
          return user.token;
        }
      }
      return "Unknown";
    }
  },
  mounted: function () {
    self = this;

    axios({
      method: "get",
      headers: {'Authorization' : 'Bearer ' + self.token },
      url: "http://dnscfg.lyingto.me:9053/v1/public/zones"
    }).then(function (response) {

        self.zones = response.data;


    }).catch(function (error) {

      if (error.response) {
        return error.response.data.msg;
      } else {
        return "Unexpected Error" + error;
      }

    });
  }
};
