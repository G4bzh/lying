Vue.use(VueMaterial.default)
Vue.use(VueRouter);

// Workaround to bypass "to=" error in md componement
Vue.component('router-link', Vue.options.components.RouterLink);
Vue.component('router-view', Vue.options.components.RouterView);

import SignIn from './signin.vue.js';
import SignUp from './signup.vue.js';
import SignOut from './signout.vue.js';
import Dashboard from './dashboard.vue.js';

const routes = [{
  path: '/signin',
  component: SignIn
}, {
  path: '/signup',
  component: SignUp
}, {
  path: '/signout',
  component: SignOut  
}, {
  path: '/dashboard',
  component: Dashboard,
  meta: {
    requiresAuth: true
  }
}];

const router = new VueRouter({
    routes // short for `routes: routes`
})

// Auth guard
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const user = JSON.parse(localStorage.getItem('user'));

  if (requiresAuth && !user) {
     next({
       path: '/signin',
       query: { redirect: to.fullPath }
     });
  } else {
    next();
  }
});

var app = new Vue({
  router,
  el: '#app',
  components : {
    Dashboard,
    SignIn,
    SignUp
  },
  data : {
    isAuth : false
  }
});
