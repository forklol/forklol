<template>
  <div>
    <div style="height:100%; width: 100%;">
      <canvas id="canvas-feemarket" :style="{height: chart_height}"></canvas>
    </div>
  </div>
</template>

<script>
  import moment from 'moment'
  import Chart from 'chart.js'
  import {sprintf} from 'sprintf-js'
  import numeral from 'numeral'

  export default {
    name: 'feemarket-chart',
    mounted () {
      this.redraw()
    },
    data() {
      return {
        chart_height: '256px'
      }
    },
    methods: {
      redraw () {
        let labels = this.data.BTC.summary.fee_market.map((obj) => moment.unix(obj.timestamp).format("MMM, DD"))

        this.add_dataset({
          type: 'bar',
          fill: true,
          label: 'Fees',
          borderColor: this.btc_c,
          backgroundColor: this.btc_c,
          data: this.data.BTC.summary.fee_market.map((obj) => this.sprintf('%.2f', obj.fees)),
          yAxisID: 'sum',
        })

        this.add_dataset({
          type: 'line',
          label: '%',
          borderColor: this.btc_c,
          borderDash: [3, 3],
          data: this.data.BTC.summary.fee_market.map((obj) => this.sprintf('%.2f', obj.percentage)),
          yAxisID: 'pct',
        }, (d) => {
          d.lineTension = 0.2
        })

        this.add_dataset({
          type: 'bar',
          fill: true,
          label: 'Fees',
          borderColor: this.bch_c,
          backgroundColor: this.bch_c,
          data: this.data.BCH.summary.fee_market.map((obj) => this.sprintf('%.2f', obj.fees)),
          yAxisID: 'sum',
        })

        this.add_dataset({
          type: 'line',
          label: '%',
          borderColor: this.bch_c,
          borderDash: [3, 3],
          data: this.data.BCH.summary.fee_market.map((obj) => this.sprintf('%.2f', obj.percentage)),
          yAxisID: 'pct',
        }, (d) => {
          d.lineTension = 0.2
        })


        var config = {
          type: 'bar',
          options: {
            title: {
              display: true,
              text: 'Fee market development (past 30 days)',
            },
            scales: {
              xAxes: [{
                display: true,
              }],
              yAxes: [{
                id: 'sum',
                display: true,
                position: 'right',
                scaleLabel: {
                  display: true,
                  labelString: 'TX fees per day (USD)'
                },
                ticks: {
                  callback: (value, index, values) => {
                    return this.numeral('$0,0', value)

                    let steps = [10000, 100000, 100000, 1000000, 5000000]

                    if (steps.indexOf(value) !== -1) {
                      return this.numeral('$0,0', value)
                    }
                    return null
                  },
                  min: 0,
                  maxTicksLimit: 10
                },
              }, {
                display: true,
                id: 'pct',
                position: 'left',
                scaleLabel: {
                  display: true,
                  labelString: 'Percentage of total reward'
                },
                ticks: {
                  callback: (val) => {
                    return this.sprintf('%.0f%%', val)
                  }
                }
              }]
            }
          }
        };

        this.draw('canvas-feemarket', config, labels)
      },

    }
  }
</script>

<style>

</style>
