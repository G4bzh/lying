export default {
  name: "login",
  template: `
  <div>
      <label><b>Username</b></label>
      <input type="text" placeholder="Enter Username" v-model="id" >

      <label><b>Password</b></label>
      <input type="password" placeholder="Enter Password" v-model="password" >

      <button v-on:click="doLogin">Login</button>

      {{ message }}

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
