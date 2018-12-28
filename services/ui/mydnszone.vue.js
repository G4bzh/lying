import * as URL from "./url.js"
import RrEdit from "./rredit.vue.js"

export default {
  name: "MyDnsZone",
  template: `
  <div>
    <md-card>


    <md-card-header>

      <div class="md-title">Zone {{ $route.params.zone }}</div>
    </md-card-header>

    <md-card-content>


      <div v-for="rr,index in rrs" :key="rr.name+rr.type+rr.data">
          <rr-edit
            v-bind:rrname=rr.name
            v-bind:rrttl=rr.ttl
            v-bind:rrclass=rr.class
            v-bind:rrtype=rr.type
            v-bind:rrdata=rr.rdata
            v-bind:rrindex=index
            @rr-remove="doRemove">
          </rr-edit>
      </div>

      <md-progress-bar md-mode="query" v-if="loading" ></md-progress-bar>
      <rr-edit @rr-add="doAdd"  v-else></rr-edit>

    </md-card-content>

    </md-card>
  </div>
  `,
  data: function() {
    return {
      rrs : [],
      loading: true
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
      this.loading = false;

    }).catch(function (error) {

      if (error.response) {
        return error.response.data.msg;
      } else {
        return "Unexpected Error" + error;
      }

    });
  },
  methods: {
    doAdd: function() {
      console.log("Adding Object");
      return;
    },
    doRemove: function(ind) {
      this.rrs =  this.rrs.filter(function(e,i) { return i != ind });
      return;
    }
  }
};
