export default {
  name: "login",
  template: `
  <div>

    <md-content class="md-elevation-3">

      <div class="md-title">Sign In </div>

      <md-field md-clearable>
        <label>Username</label>
        <md-input placeholder="Enter Username" v-model="id"></md-input>
      </md-field>

      <md-field>
        <label>Password</label>
        <md-input v-model="password" type="password" placeholder="Enter Password"></md-input>
      </md-field>

        <md-button class="md-raised md-primary" v-on:click="doLogin">Login</md-button>
        {{ message }}

    </md-content>


  </div>
  `,
  data: function() {
    return {
      message: "",
      id: "",
      password: ""
    }
  },
  methods : {
    doLogin: function() {
      self = this;
      axios({
        method: "post",
        url: "http://auth.lyingto.me:9080/v1/login",
        data: {
          id: self.id,
          password: self.password
        }
      }).then(function (response) {

        self.message = response.data.token;

      }).catch(function (error) {

        if (error.response) {
          self.message = error.response.data.msg;
        } else {
          self.message = "Unexpected Error";
        }

      });
    }
  }
}
