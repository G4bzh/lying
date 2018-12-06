export default {
  name: "SignIn",
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

        <md-button class="md-raised md-primary" v-on:click="doSignIn">Sign In</md-button>
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
    doSignIn: function() {
      self = this;
      axios({
        method: "post",
        url: "http://auth.lyingto.me:9080/v1/login",
        data: {
          id: self.id,
          password: self.password
        }
      }).then(function (response) {

        localStorage.setItem('user', JSON.stringify({id:self.id,token:response.data.token}));
        self.$parent.isAuth = true;

        if (self.$route.query.redirect) {
          self.$router.push(self.$route.query.redirect);
        } else {
          self.$router.push("/");
        }

      }).catch(function (error) {

        if (error.response) {
          self.message = error.response.data.msg;
        } else {
          self.message = "Unexpected Error" + error;
        }

        self.$parent.isAuth = false;
        localStorage.removeItem('user');

      });
    }
  }
}
