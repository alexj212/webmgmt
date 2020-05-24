import VueTerminal from './components/VueTerminal.vue'

const terminal = {
  version: /* eslint-disable no-undef */ __VERSION__,
  install,
  VueTerminal
}

if (typeof window !== 'undefined' && window.Vue) {
  window.Vue.use(install)
}

export default terminal

function install (Vue) {
  if (install.installed) {
    return
  }
  Vue.component(VueTerminal.name, VueTerminal)
}
