export default {
  name: "SignUp",
  template: `
  <div>
    <div class="centered-container">
    <md-content class="md-elevation-3">

      <div class="md-title">Sign Up </div>

      <md-field md-clearable>
        <label>Mail Address</label>
        <md-input placeholder="Enter your mail address" v-model="id"></md-input>
      </md-field>

      <md-field>
        <label>Password</label>
        <md-input v-model="password01" type="password" placeholder="Enter a Password"></md-input>
      </md-field>

      <md-field>
        <label>Confirm Password</label>
        <md-input v-model="password02" type="password" placeholder="Confirm Password"></md-input>
      </md-field>

        <md-button class="md-raised md-primary" v-on:click="doSignUp">Sign Up</md-button>
        {{ message }}

    </md-content>
    </div>

  </div>
  `,
  data: function() {
    return {
      message: "",
      id: "",
      password01: "",
      password02: ""
    }
  },
  methods : {
    doSignUp: function() {
      self = this;
      if (self.password01 != self.password02)
      {
          self.message = "Passwords Mismatch";
          return;
      }
      axios({
        method: "put",
        url: "http://auth.lyingto.me:9080/v1/login",
        data: {
          id: self.id,
          password: self.password01
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
