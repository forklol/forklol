<template>
  <div id="app">
    <div v-if="loaded === false" class="loading">
      Loading data..
      <br>
      <br>
      <img src="/static/logo.png" alt="">
    </div>
    <div v-if="loaded">
      <!-- As a link -->
      <nav class="navbar navbar-light bg-faded lol-nav">
        <a class="navbar-brand" href="/">
          <img src="/static/logo.png" alt="fork.lol" width="30" height="30"> fork.lol
        </a>
      </nav>

      <div class="container-fluid">

        <div class="row justify-content-md-center">
          <div class="col-xl-10 col-lg-12 col-md-12 col-sm-12 col-xs-12">

            <div class="card lol-window">
              <div class="card-block">
                <ul class="nav nav-pills">
                  <li class="nav-item">
                    <router-link class="nav-link" :to="{name: 'home'}" exact>
                      <i class="fa fa-home"></i>
                    </router-link>
                  </li>
                  <li class="nav-item">
                    <router-link class="nav-link" :to="{name: 'reward'}">Reward</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link class="nav-link" :to="{name: 'blocks'}">Blocks</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link class="nav-link" :to="{name: 'tx'}">Tx</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link class="nav-link" :to="{name: 'pow'}">POW</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link class="nav-link" :to="{name: 'security'}">Security</router-link>
                  </li>
                </ul>
                <hr>
                <div class="lol-window-inner">
                  <transition name="fade">
                    <router-view></router-view>
                  </transition>
                </div>
              </div>

              <div class="card-footer text-muted text-center">
                created by
                <router-link :to="{name: 'contact'}">trippysalmon</router-link>
                / {{ version }} / blockdata kindly provided by
                <a target="_blank" href="http://blockchair.com">blockchair.com</a> / exchange rate by
                <a href="http://www.bitcoinaverage.com">bitcoinaverage.com</a>
              </div>
            </div>
          </div>
        </div>
      </div>
      <br>

    </div>
    <notifications/>
  </div>
</template>

<script>
  import Chart from 'chart.js'
  import axios from 'axios'
  import moment from 'moment'
  import numeral from 'numeral'
  import {sprintf} from 'sprintf-js'


  export default {
    name: 'app',
    data () {
      return {
        visible: true,
        gone: false,
      }
    },
    methods: {
      refresh () {
        this.$store.dispatch('get_last')
      },
      notify (t) {
        this.$notify({
          title: t,
          duration: 5000,
        })
      }
    },
    watch: {
      loaded () {
        this.$store.dispatch('poll_data')
        //this.$store.dispatch('get_last')
      },
      btc_h (n, o) {
        if (o == 0) return

        let t = 'New BTC block found: #' + n

        let nw = n - o

        if (this.gone) {
          if (nw > 1) {
            t = `${n - o} new BTC blocks found since you were last here`
          } else {
            t = `${n - o} new BTC block found since you were last here`
          }
        }

        this.notify(t)

        setTimeout(() => {
          this.gone = false
        }, 100)
      },
      bch_h (n, o) {
        if (o == 0) return

        let t = 'New BCH block found: #' + n

        let nw = n - o

        if (this.gone) {
          if (nw > 1) {
            t = `${n - o} new BCH blocks found since you were last here`
          } else {
            t = `${n - o} new BCH block found since you were last here`
          }
        }

        this.notify(t)

        setTimeout(() => {
          this.gone = false
        }, 100)
      },
      visible (n, o) {
        if (n === true) {
          this.refresh()
          this.$store.commit('autorefresh', true)
          this.$store.commit('visible', true)
          //this.$store.dispatch('get_last')
          this.$store.dispatch('poll_data')
          setTimeout(() => {
            this.gone = false
          }, 5000)
        } else {
          //this.$store.commit('autorefresh', false)
          this.$store.commit('visible', false)
          this.$store.dispatch('poll_stop')
          this.gone = true
        }
      }
    },
    computed: {
      version () {
        return process.env.VERSION
      }
    },
    mounted () {
      this.$store.dispatch('get_last', false)

      if (this.$route.name === null) {
        this.$router.replace({path: '', redirect: '/'})
      }

      var stateKey, eventKey, keys = {
        hidden: "visibilitychange",
        webkitHidden: "webkitvisibilitychange",
        mozHidden: "mozvisibilitychange",
        msHidden: "msvisibilitychange"
      };

      for (stateKey in keys) {
        if (stateKey in document) {
          eventKey = keys[stateKey];
          break;
        }
      }

      if (process.env.NODE_ENV === 'development') {
        document.title = document.title + ' DEV'
      }

      document.addEventListener(eventKey, (v) => {
        this.visible = !document[stateKey]
      });
    }
  }
</script>

<style>
  html {
    overflow-y: scroll;
  }

  .loading {
    text-align: center;
    margin-top: 48px;
  }

  #app, .nav-tabs .nav-link.active {
    color: #000000 !important;
  }

  .nav-tabs {
    font-size: 0.8rem;
  }

  @media screen and (max-width: 460px) {
    .nav-tabs, .nav-pills {
      overflow-y: scroll;
      white-space: nowrap
    }
  }

  .notifications {
    width: auto !important;
  }

  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    margin-top: 0px;
  }

  .lol-nav {
    margin-bottom: 16px;
  }

  .lol-window {
    margin-left: 1%;
    margin-right: 1%;
  }

  .lol-window-inner {
    padding: 0px;
    padding-top: 0px;
  }

  .lol-window-chart {
    padding: 16px;
  }

  .card-footer {
    font-size: 12px;
  }

  .card-block {
    padding: 0.75rem !important;
  }

  .fade-leave-active {
    transition: opacity 0s
  }

  .fade-enter-active {
    transition: opacity .2s
  }

  .fade-enter,
  .fade-leave-to
    /* .fade-leave-active below version 2.1.8 */

  {
    opacity: 0
  }
</style>
