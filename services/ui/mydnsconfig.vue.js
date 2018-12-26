import * as URL from "./url.js"

export default {
  name: "MyDnsConfig",
  template: `
  <div class="md-layout md-gutter">

    <md-card class="md-layout-item">

      <md-card-header>
        <div class="md-title">Forwarders</div>
      </md-card-header>

      <md-card-content>
        <md-chips v-model="forwarders" md-placeholder="Add forwarder..." ></md-chips>
      </md-card-content>

    </md-card>

    <md-card class="md-layout-item">

      <md-card-header>
        <div class="md-title">Public IPs</div>
      </md-card-header>

      <md-card-content>
        <md-chip v-for="i in ips" :key="i" md-deletable  @md-delete=rmIP>{{ i }}</md-chip>
        <md-button class="md-primary" @click="addIP">Add my current ip : {{ ip }}</md-button>
      </md-card-content>

    </md-card>

  </div>
  `,
  data: function() {
    return {
      forwarders: [],
      ips: [],
      ip: ""
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
      url: URL.GETFORWARDERS
    }).then(response => {

        this.forwarders = response.data;


    }).catch(error => {

      if (error.response) {
        return error.response.data.msg;
      } else {
        return "Unexpected Error" + error;
      }

    });

    axios({
      method: "get",
      url: URL.GETPUBLICIP
    }).then(response => {

        this.ip = response.data.ip;


    }).catch(error => {

      if (error.response) {
        this.ip =  error.response.data.msg;
      } else {
        this.ip = "Unexpected Error" + error;
      }

    });

  },
  methods: {
    rmIP: function (event) {
        this.ips =  this.ips.filter(function(e) { return event.currentTarget.parentElement.textContent.trim() != e.trim() })
    },
    addIP: function () {

      if (! this.ips.includes(this.ip)) {
        this.ips.push(this.ip);
      }

    }
  }
};