import * as URL from "./url.js"
import ZoneEdit from "./zonedit.vue.js"

export default {
  name: "MyDnsZone",
  template: `
  <div>

      <zone-edit rrname="toto" v-bind:rrttl=4></zone-edit>
      Zone {{ $route.params.zone }}
      <ul>
        <li v-for="rr in rrs">
          {{ rr.name }} {{ rr.ttl }} {{ rr.class }} {{ rr.type }} {{ rr.rdata }}
        </li>
      </ul>
  </div>
  `,
  data: function() {
    return {
      rrs : []
    }
  },
  components: {
    ZoneEdit
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
