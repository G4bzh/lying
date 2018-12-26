import * as URL from "./url.js"

export default {
  name: "MyDashboard",
  template: `
  <div class="md-layout md-gutter md-size-50 md-small-size-100">

    <md-card class="md-layout-item" md-with-hover @click.native="$router.push('/mylies/mysettings')">
      <md-ripple>
        <md-card-header>
          <div class="md-title">Welcome</div>
        </md-card-header>

        <md-card-content>
          {{username}}
        </md-card-content>
      </md-ripple>
    </md-card>

    <md-card class="md-layout-item" md-with-hover @click.native="$router.push('/mylies/mydns')">
      <md-ripple>
        <md-card-header>
          <div class="md-title">Your Zones</div>
        </md-card-header>

        <md-card-content>
          {{ n }}
        </md-card-content>

      </md-ripple>
    </md-card>

  </div>
  `,
  data: function() {
    return {
      n: 0
    }
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
    },
    token: function () {
      const user = JSON.parse(localStorage.getItem('user'));
      if (user) {
        if (user.token) {
          return user.token;
        }
      }
      return "Unknown";
  }},
  mounted: function () {

    axios({
      method: "get",
      headers: {'Authorization' : 'Bearer ' + this.token },
      url: URL.GETZONES
    }).then(response => {

      this.n = response.data.length;


    }).catch(error => {

      if (error.response) {
        return error.response.data.msg;
      } else {
        return "Unexpected Error" + error;
      }

    });
  }
};
