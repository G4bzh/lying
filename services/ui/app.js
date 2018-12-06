Vue.use(VueMaterial.default)
Vue.use(VueRouter);

// Workaround to bypass "to=" error in md componement
Vue.component('router-link', Vue.options.components.RouterLink);
Vue.component('router-view', Vue.options.components.RouterView);

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
