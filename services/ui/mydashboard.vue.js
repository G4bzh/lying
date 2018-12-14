export default {
  name: "MyDashboard",
  template: `
  <div>
    <md-card md-with-hover>

        <md-card-header>
          <div class="md-title">Welcome</div>
        </md-card-header>

        <md-card-content>
          {{username}}
        </md-card-content>

    </md-card>

    <md-card md-with-hover @click.native="$router.push('/mylies/mydns')">
      <md-ripple >
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
    self = this;

    axios({
      method: "get",
      headers: {'Authorization' : 'Bearer ' + self.token },
      url: "http://dnscfg.lyingto.me:9053/v1/public/zones"
    }).then(function (response) {

      self.n = response.data.length;


    }).catch(function (error) {

      if (error.response) {
        return error.response.data.msg;
      } else {
        return "Unexpected Error" + error;
      }

    });
  }
};
