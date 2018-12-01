export default {
  name: "login",
  template: `
  <div>
      <label><b>Username</b></label>
      <input type="text" placeholder="Enter Username" v-model="id" >

      <label><b>Password</b></label>
      <input type="password" placeholder="Enter Password" v-model="passwd" >

      <button v-on:click="doLogin">Login</button>

      {{ message }}

  </div>
  `,
  data: function() {
    return {
      message: "",
      id: "",
      passwd: ""
    }
  },
  methods : {
    doLogin: function() {
      this.message = this.id + "/" + this.passwd
    }
  }
}
