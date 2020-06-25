import Vue from 'vue'
import VueMarkdown from 'vue-markdown'
import App from './App.vue'
import router from './router'
import store from './store'

Vue.config.productionTip = false

Vue.component('vue-markdown', VueMarkdown);
new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
