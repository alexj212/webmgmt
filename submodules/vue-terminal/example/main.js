import Vue from 'vue'
import terminal from '../src'
import App from './App.vue'

Vue.use(terminal)

/* eslint-disable no-new */
new Vue({
  render(createElement) {
    return createElement(App)
  }
}).$mount('#app')
