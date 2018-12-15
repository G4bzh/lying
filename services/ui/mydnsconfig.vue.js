export default {
  name: "MyDnsConfig",
  template: `
  <div>
    DNS Config
    Forwarders :
    <ul>
      <li v-for="forwarder in forwarders">
        {{ forwarder }}
      </li>
    </ul>
  </div>
  `,
  data: function() {
    return {
      forwarders: []
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
      url: "http://dnscfg.lyingto.me:9053//v1/public/forwarders"
    }).then(response => {

        this.forwarders = response.data;


    }).catch(error => {

      if (error.response) {
        return error.response.data.msg;
      } else {
        return "Unexpected Error" + error;
      }

    });
  }
};
