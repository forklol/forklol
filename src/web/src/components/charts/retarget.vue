<template>
  <div>
    <p>
      Estimated days for difficulty retarget.
    </p>
    <div style="height:100%; width: 100%;">
      <canvas id="canvas-retarget" style="height: 512px;"></canvas>
    </div>
    <br>
    <small>
      <p>
        Estimated next difficulty. The 12 hour hashrate average is used to make a prediction about upcoming blocks in the period.
      </p>

      <table class="table table-sm table-bordered">
        <thead class="thead-dark">
        <tr>
          <th>Coin</th>
          <th>Date (UTC 24h)</th>
          <th>Remaining</th>
          <th>Change</th>
        </tr>
        </thead>
        <tbody>
        <tr>
          <td>BTC</td>
          <td>{{ date('BTC') }}</td>
          <td>
            {{ days('BTC') }}
            <small>({{ left('BTC') }} blk)</small>
          </td>
          <td>
            <i v-if="next_pct('BTC', true) > 1" class="fa fa-chevron-up up"></i>
            <i v-if="next_pct('BTC', true) < 1" class="fa fa-chevron-down down"></i>
            {{ next('BTC') }}
          </td>
        </tr>
        <tr>
          <td>BCH *</td>
          <td>-</td>
          <td>-</td>
          <td>-</td>
        </tr>
        </tbody>
      </table>
      * On November 13th 2017 BCH switches to a faster changing difficulty algorithm which makes it impractical to make predictions about its behaviour.
    </small>
  </div>
</template>

<script>
  import moment from 'moment'
  import Chart from 'chart.js'
  import Vue from 'vue'
  import humanize from 'humanize-duration'

  export default {
    name: 'pow-retarget',
    mounted () {
      this.redraw()
    },
    methods: {
      redraw () {
        let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

        this.add_dataset({
          label: 'BTC',
          borderColor: this.btc_c,
          data: this.data.BTC.history['all'].map((obj) => this.estimate('BTC', obj)),
          spanGaps: true,
        })

        this.add_dataset({
          label: 'BCH',
          borderColor: this.bch_c,
          data: this.data.BCH.history['all'].map((obj) => this.estimate('BCH', obj)),
          spanGaps: true,
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
                  labelString: 'Days'
                },
                ticks: {
                  callback: (value, index, values) => {
                    return this.numeral('0,0', value)
                  }
                }
              }]
            }
          }
        };

        this.draw('canvas-retarget', config, labels)
      },
    }
  }
</script>

<style scoped>
  .up {
    color: green;
  }

  .down {
    color: red;
  }
</style>
