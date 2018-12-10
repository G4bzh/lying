export default {
  name: "SignIn",
  template: `
  <div>

    <md-content class="md-elevation-3">

      <div class="md-title">{{ title }}</div>

      <md-field md-clearable>
        <label>Mail Address</label>
        <md-input placeholder="Enter Email Address" v-model="id"></md-input>
      </md-field>

      <md-field md-clearable v-show=isSignUp>
        <label>Username</label>
        <md-input placeholder="Enter an Username" v-model="username"></md-input>
      </md-field>

      <md-field>
        <label>Password</label>
        <md-input v-model="password" type="password" placeholder="Enter Password"></md-input>
      </md-field>

      <md-checkbox v-model="isSignUp">I want to Sign Up</md-checkbox>
      <md-button class="md-raised md-primary" v-on:click="doSign">{{ title }}</md-button>
        {{ message }}

    </md-content>


  </div>
  `,
  data: function() {
    return {
      message: "",
      id: "",
      password: "",
      username: "",
      isSignUp : false
    }
  },
  computed: {
    title: function() {
      return (this.isSignUp ? "Sign Up" : "Sign In");
    }
  },
  methods : {
    doSign: function() {
      self = this;
      if (self.isSignUp) {
        axios({
          method: "put",
          url: "http://auth.lyingto.me:9080/v1/login",
          data: {
            id: self.id,
            password: self.password,
            username: self.username
          }
        }).then(function (response) {

          self.$router.push("/");


        }).catch(function (error) {

          if (error.response) {
            self.message = error.response.data.msg;
          } else {
            self.message = "Unexpected Error";
          }

        });
      } else {

        axios({
          method: "post",
          url: "http://auth.lyingto.me:9080/v1/login",
          data: {
            id: self.id,
            password: self.password
          }
        }).then(function (response) {

          localStorage.setItem('user', JSON.stringify({id:self.id,token:response.data.token,name:response.data.username}));

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

          localStorage.removeItem('user');

        });
      }
    }
  }
}
