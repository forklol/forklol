<template>
  <div>
    <div style="height:100%; width: 100%;">
      <canvas id="canvas-rdari" :style="{height: chart_height}"></canvas>
    </div>
  </div>
</template>

<script>
  import moment from 'moment'
  import Chart from 'chart.js'
  import {sprintf} from 'sprintf-js'
  import numeral from 'numeral'

  export default {
    name: 'rdari-chart',
    mounted () {
      this.redraw()
    },
    data () {
      return {
        chart_height: '256px'
      }
    },
    methods: {
      rdari (coin1, coin2, n) {
        let a = this.data[coin1].summary.dari_relative[n]
        let b = this.data[coin2].summary.dari_relative[n]

        if (!a || !b) return null

        if (a > b) {
          return this.sprintf('%.3f', a / b)
        } else {
          return '1.00'
        }
      },
      redraw () {
        let labels = this.data.BTC.summary.dari_relative.map((obj, n) => n)

        this.add_dataset({
          type: 'line',
          label: 'BTC',
          borderColor: this.btc_c,
          spanGaps: true,
          steppedLine: true,
          borderWidth: 1,
          data: this.data.BTC.summary.dari_relative.map((obj, n) => this.rdari('BTC', 'BCH', n)),
        }, (d) => {
          d.lineTension = 0
        })

        this.add_dataset({
          type: 'line',
          label: 'BCH',
          borderColor: this.bch_c,
          spanGaps: true,
          steppedLine: true,
          borderWidth: 1,
          data: this.data.BCH.summary.dari_relative.map((obj, n) => this.rdari('BCH', 'BTC', n)),
        }, (d) => {
          d.lineTension = 0
        })

        let config = {
          type: 'line',
          options: {
            legend: {
              display: true
            },
            title: {
              display: true,
              text: 'Relative miner reward for work done (DARI ratio, past week)',
            },
            scales: {
              xAxes: [{
                display: false,
              }],
              yAxes: [{
                display: true,
                position: 'right',
                scaleLabel: {
                  display: false,
                  labelString: 'Relative DARI'
                },
                ticks: {
                  callback: (v) => {
                    return this.sprintf('%.2fx', v)
                  }
                }
              }]
            }
          }
        };

        this.draw('canvas-rdari', config, labels)
      },

    }
  }
</script>

<style>

</style>
