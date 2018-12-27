import * as URL from "./url.js"
import RrEdit from "./rredit.vue.js"

export default {
  name: "MyDnsZone",
  template: `
  <div>

      Zone {{ $route.params.zone }}
      <div v-for="rr,index in rrs">
          <rr-edit
            v-bind:rrname=rr.name
            v-bind:rrttl=rr.ttl
            v-bind:rrclass=rr.class
            v-bind:rrtype=rr.type
            v-bind:rrdata=rr.rdata
            v-bind:rrindex=index>
          </rr-edit>
      </div>
      <rr-edit></rr-edit>

  </div>
  `,
  data: function() {
    return {
      rrs : []
    }
  },
  components: {
    RrEdit
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

    axios({
      method: "get",
      headers: {'Authorization' : 'Bearer ' + this.token },
      url: URL.GETZONE +  this.$route.params.zone
    }).then(response => {

      this.rrs = response.data.rrs;

    }).catch(function (error) {

      if (error.response) {
        return error.response.data.msg;
      } else {
        return "Unexpected Error" + error;
      }

    });
  }
};
