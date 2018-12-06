Vue.use(VueMaterial.default)
Vue.use(VueRouter);

import SignIn from './signin.vue.js';
import SignUp from './signup.vue.js';


const routes = [{
  path: '/signin',
  component: SignIn
}, {
  path: '/signup',
  component: SignUp
}];

const router = new VueRouter({
    routes // short for `routes: routes`
})

var app = new Vue({
  router,
  el: '#app',
  components : {
    SignIn,
    SignUp
  }
});
