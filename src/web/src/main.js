require("jquery")
require("font-awesome-webpack")
require('bootstrap-loader')

// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import Notifications from 'vue-notification'
import App from './App'
import router from './router'
import store from './store'


import commonMixin from './mixins/common.js'
import chartMixin from './mixins/chart.js'
import statTable from './components/misc/stat-table.vue'

Vue.mixin(commonMixin)
Vue.mixin(chartMixin)
Vue.use(Notifications)
Vue.component('stat-table', statTable)

import Chart from 'chart.js'

Chart.defaults.global.responsive = true
Chart.defaults.global.maintainAspectRatio = false
Chart.defaults.global.animation = false
Chart.defaults.global.elements.line.backgroundColor = "rgba(0,0,0,0)"
Chart.defaults.global.elements.line.fill = false
Chart.defaults.global.elements.line.borderWidth = 1
Chart.defaults.global.elements.point.radius = 0
//Chart.defaults.global.tooltips = { "mode": false, "intersect": false }
Chart.defaults.global.title.display = false
Chart.defaults.global.title.display = false
Chart.defaults.global.defaultFontColor = '#000'
Chart.defaults.global.defaultFontSize = 12

/* eslint-disable no-new */
new Vue({
    el: '#app',
    router,
    store,
    template: '<App/>',
    components: { App }
  })
