export default {
  name: "MyDnsZone",
  template: `
  <div>

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
      url: "http://dnscfg.lyingto.me:9053//v1/public/zone/" +  this.$route.params.zone
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
