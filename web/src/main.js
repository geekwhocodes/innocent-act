import Vue from 'vue';
import Buefy from 'buefy';
import Unicon from 'vue-unicons';
import { uni500Px } from 'vue-unicons/src/icons';

import App from './App.vue';
import router from './router';
import store from './store';
import variables from './assets/variables.scss';

Vue.use(Buefy, {});
Vue.config.productionTip = false;

Unicon.add(uni500Px);
Vue.use(Unicon, {
  fill: variables.primaryColor,
});

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
