Vue.use(VueMaterial.default)
Vue.use(VueRouter);

// Workaround to bypass "to=" error in md componement
Vue.component('router-link', Vue.options.components.RouterLink);
Vue.component('router-view', Vue.options.components.RouterView);

import SignIn from './signin.vue.js';
import SignOut from './signout.vue.js';
import MyLies from './mylies.vue.js';
import MyDashboard from './mydashboard.vue.js'
import MyDns from './mydns.vue.js'
import MyDnsConfig from './mydnsconfig.vue.js'
import MyDnsZone from './mydnszone.vue.js'
import MySettings from './mysettings.vue.js'

const routes = [{
  path: '/signin',
  component: SignIn
}, {
  path: '/signout',
  component: SignOut
}, {
  path: '/mylies',
  component: MyLies,
  meta: {
    requiresAuth: true
  },
  children: [
    {
      path: '',
      component: MyDashboard
    },
    {
      path: 'mydashboard',
      component: MyDashboard
    },
    {
      path: 'mydns/',
      component: MyDns,
      children: [
        {
          path : '',
          component: MyDnsConfig
        },
        {
          path: 'config',
          component: MyDnsConfig
        },
        {
          path: 'zone/:zone',
          component: MyDnsZone
        }
      ]
    },
    {
      path: 'mysettings',
      component: MySettings
    }
  ]
}];

const router = new VueRouter({
    routes // short for `routes: routes`
})

// Auth guard
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const user = JSON.parse(localStorage.getItem('user'));
  var token = null;

  if (user) {
    if (user.token) {
      token=user.token;
    }
  }

  if (requiresAuth && !token) {
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
  data : {
  },
  methods : {
    isAuth : function() {
      var user = JSON.parse(localStorage.getItem('user'));
      if (user) {
        if (user.token) {
          return true;
        }
      }
      return false;
    }
  }
});
