Vue.use(VueMaterial.default)

import SignIn from './signin.vue.js';
import SignUp from './signup.vue.js';


var app = new Vue({
  el: '#app',
  components : {
    SignIn,
    SignUp
  }
});
