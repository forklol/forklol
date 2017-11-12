<template>
  <div>
    <p>
      Relative speed of blocks found compared to the target of 6 blocks per hour (although the 6 blocks per hour is only a proxy for the actual target of 2016 blocks every two weeks).
    </p>
    <div style="height:100%; width: 100%;">
      <canvas id="canvas-speed" style="height: 512px;"></canvas>
    </div>
    <br>
    <stat-table :stats="stats"></stat-table>
  </div>
</template>

<script>
  import moment from 'moment'
  import Chart from 'chart.js'
  import Vue from 'vue'

  export default {
    name: 'pow-speed',
    mounted () {
      this.redraw()
    },
    computed: {
      stats () {
        return {
          header: ['3h', '12h', '1d', '3d', '7d'],
          f: (coin, avg) => {
            let spd = this.speed(this.data[coin].averages.last, 'rate' + avg)
            return spd + "%"
          },
          coins: [{
            coin: 'BTC',
            values: ['3h', '12h', '1d', '3d', '7d']
          }, {
            coin: 'BCH',
            values: ['3h', '12h', '1d', '3d', '7d']
          }]
        }

      }
    },
    methods: {
      speed (obj, rate = null) {
        if (obj.blocks < 1) return '0.00'

        if (rate !== null) {
          let spd = (obj[rate] / obj.avg_diff) - 1
          spd *= 100

          let formatted = this.sprintf('%.2f', spd)
          return spd > 0 ? '+' + formatted : formatted
        }

        return null
      },
      redraw () {
        let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

        this.add_dataset({
          label: 'BTC',
          borderColor: this.btc_c,
          data: this.data.BTC.history['all'].map((obj) => this.speed(obj, 'rate3h'))
        })

        this.add_dataset({
          label: 'BCH',
          borderColor: this.bch_c,
          data: this.data.BCH.history['all'].map((obj) => this.speed(obj, 'rate3h'))
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
                position: 'right',
                scaleLabel: {
                  display: true,
                  labelString: 'Relative speed'
                },
                ticks: {
                  callback(val) {
                    return val + "%"
                  },
                  min: -100,
                }
              }]
            }
          }
        };

        this.draw('canvas-speed', config, labels)
      }
    }
  }
</script>

<style>

</style>
