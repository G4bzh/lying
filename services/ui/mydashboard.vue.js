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
          {{ domains }}
        </md-card-content>

        <md-card-expand>

          <md-card-actions md-alignment="space-between">
            <md-card-expand-trigger>
              <md-button class="md-icon-button">
                <md-icon>keyboard_arrow_down</md-icon>
              </md-button>
            </md-card-expand-trigger>
          </md-card-actions>

          <md-card-expand-content>
            <md-card-content>
              {{ token }}
            </md-card-content>
          </md-card-expand-content>
        </md-card-expand>

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
    return {
      n_domains : 2
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
    },
    domains: function() {
      self = this;

      axios({
        method: "get",
        headers: {'Authorization' : 'Bearer ' + self.token },
        url: "http://dnscfg.lyingto.me:9053/v1/public/zones"
      }).then(function (response) {

        return response.data.domain[0];


      }).catch(function (error) {

        if (error.response) {
          return error.response.data.msg;
        } else {
          return "Unexpected Error" + error;
        }

      });
    }
  }
};
