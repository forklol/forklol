<template>
  <div>
    <p>
      Normalized miner income for work done (DARI)
    </p>
    <div style="height:100%; width: 100%;">
      <canvas id="canvas-dari" style="height: 512px;"></canvas>
    </div>
    <br>
    <p>
      The DARI (Difficulty Adjusted Reward Index) is a way to compare the rewards of two chains that share the same reward system. The following calculation is used to determine the DARI:
      <br>
    </p>
    <div class="text-center">
      <code>
        (block coinbase+fees in satoshis) / (block difficulty) * (exchange rate in USD)
      </code>
    </div>
    <br>
    <p>
      With this result we can compare the two chains. For example, using the DARI of the last block we can estimate that a {{ leader
      }} miner could potentially collect ({{ dari(leader) }} / {{ dari(loser)}} =)
      <strong>{{ leader_ratio }}</strong> times the reward (in USD) that a {{ loser }} miner could collect.
    </p>

    <p>
      Since the inception of BCH we have seen that some miners switch based on which coin is most rewarding for the work they perform.
      The data below is an attempt to capture some of the effects of this so called "chain hopping".
      <br>
      <br>
      Whether the chain hopping actually results in more income for a miner depends on many factors. Most notably the pool
      <a href="https://en.bitcoin.it/wiki/Comparison_of_mining_pools" target="_blank">reward type</a>, the exchange rate
      (since coins need to mature 100 blocks before they can be spent), the cash out strategy of the miner and of course plain old luck :)
    </p>

    <ul class="nav nav-pills justify-content-center small">
      <li class="nav-item">
        <router-link class="nav-link" :to="{name: 'reward.dari.btc'}">BTC</router-link>
      </li>
      <li class="nav-item">
        <router-link class="nav-link" :to="{name: 'reward.dari.bch'}">BCH</router-link>
      </li>
    </ul>
    <br>

    <transition name="fade">
      <router-view :key="$route.fullPath"></router-view>
    </transition>
  </div>
</template>

<script>
  import moment from 'moment'
  import Chart from 'chart.js'
  import {sprintf} from 'sprintf-js'
  import numeral from 'numeral'
  import algo from './algo.vue'

  export default {
    name: 'dari',
    components: {
      algo
    },
    mounted () {
      this.redraw()
    },
    methods: {
      redraw () {
        let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

        this.add_dataset({
          label: 'BTC',
          borderColor: this.btc_c,
          data: this.data.BTC.history['all'].map((obj) => this.obj_dari(obj)),
          yAxisID: 'dari',
          spanGaps: true
        })

        this.add_dataset({
          label: 'Difficulty',
          borderColor: this.btc_c,
          borderDash: [3, 3],
          data: this.data.BTC.history['all'].map((obj) => this.obj_diff(obj)),
          yAxisID: 'diff',
          spanGaps: true
        })

        this.add_dataset({
          label: 'BCH',
          borderColor: this.bch_c,
          data: this.data.BCH.history['all'].map((obj) => this.obj_dari(obj)),
          yAxisID: 'dari',
          spanGaps: true
        })

        this.add_dataset({
          label: 'Difficulty',
          borderColor: this.bch_c,
          borderDash: [3, 3],
          data: this.data.BCH.history['all'].map((obj) => this.obj_diff(obj)),
          yAxisID: 'diff',
          spanGaps: true
        })

        var config = {
          type: 'line',
          options: {
            scales: {
              xAxes: [{
                display: true,
                scaleLabel: {
                  display: true,
                  labelString: 'Date/Time'
                }
              }],
              yAxes: [{
                display: true,
                id: 'dari',
                position: 'right',
                scaleLabel: {
                  display: true,
                  labelString: 'DARI'
                }
              }, {
                display: false,
                id: 'diff',
                scaleLabel: {
                  display: true,
                  labelString: 'Relative difficulty'
                }
              }]
            }
          }
        };

        this.draw('canvas-dari', config, labels)
      },
      obj_dari (obj) {
        return obj.blocks > 0 && obj.price > 1 ? this.sprintf('%.2f', obj.dari) : null
      },
      obj_diff (obj) {
        return obj.blocks > 0 && obj.price > 1 ? this.sprintf('%.2f', obj.avg_diff / 10000000000) : null
      },
    }
  }
</script>

<style>

</style>
