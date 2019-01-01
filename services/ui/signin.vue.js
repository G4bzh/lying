import * as URL from "./url.js"

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
        <md-input placeholder="Enter Username" v-model="username"></md-input>
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

      if (this.isSignUp) {
        axios({
          method: "put",
          url: URL.SIGNIN,
          data: {
            id: this.id,
            password: this.password,
            username: this.username
          }
        }).then(response => {

          this.$router.push("/");


        }).catch(error => {

          if (error.response) {
            this.message = error.response.data.msg;
          } else {
            this.message = "Unexpected Error";
          }

        });
      } else {

        axios({
          method: "post",
          url: URL.SIGNIN,
          data: {
            id: this.id,
            password: this.password
          }
        }).then(response => {

          localStorage.setItem('user', JSON.stringify({id:this.id,token:response.data.token,name:response.data.username}));

          if (this.$route.query.redirect) {
            this.$router.push(this.$route.query.redirect);
          } else {
            this.$router.push("/");
          }

        }).catch(error => {

          if (error.response) {
            this.message = error.response.data.msg;
          } else {
            this.message = "Unexpected Error" + error;
          }

          localStorage.removeItem('user');

        });
      }
    }
  }
}
