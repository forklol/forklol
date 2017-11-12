<template>
  <div>
    <p>
      BTC
    </p>
    <div style="height:512px; width: 100%;">
      <canvas id="canvas-daa" style="height: 512px;"></canvas>
    </div>
    <br>
    <p>
      In the period from {{ date_from }} to {{ date_to }} a total number of {{ blocks}} blocks were mined.
      This results in an average of <u>{{ blocktime }} minutes</u>
      to find a block.
    </p>
    <p>
      <small>
        Average number of blocks per hour over the entire period (measured over {{ interval / 60 }} minute intervals):
        <table class="table table-sm table-bordered times-table">
          <thead class="thead-dark chead">
          <tr>
            <th class="noc">Blocks/hr:</th>
            <th>15+</th>
            <th>>10-15</th>
            <th>>8-10</th>
            <th>>4-8</th>
            <th>>3-4</th>
            <th>>2-3</th>
            <th>>1-2</th>
            <th><1</th>
          </tr>
          </thead>
          <tbody>
          <tr>
            <td class="noc">Blocks:</td>
            <td>{{ times['15'] }}</td>
            <td>{{ times['10'] }}</td>
            <td>{{ times['7'] }}</td>
            <td>{{ times['5'] }}</td>
            <td>{{ times['3'] }}</td>
            <td>{{ times['2'] }}</td>
            <td>{{ times['1'] }}</td>
            <td>{{ times['0'] }}</td>
          </tr>
          <tr>
            <td class="noc">Percentage:</td>
            <td>{{ times_pct('15') }}</td>
            <td>{{ times_pct('10') }}</td>
            <td>{{ times_pct('7') }}</td>
            <td>{{ times_pct('5') }}</td>
            <td>{{ times_pct('3') }}</td>
            <td>{{ times_pct('2') }}</td>
            <td>{{ times_pct('1') }}</td>
            <td>{{ times_pct('0') }}</td>
          </tr>
          <tr>
            <td></td>
            <td colspan="3" class="times-high">
              {{ block_times_high }}
            </td>
            <td class="times-norm">
              {{ block_times_norm }}
            </td>
            <td colspan="4" class="times-low">
              {{ block_times_low }}
            </td>
          </tr>
          </tbody>
        </table>
      </small>


      <small>
        Number of blocks mined during periods of relative block reward per amount of work done compared to BCH.
        <table class="table table-sm table-bordered times-table">
          <thead class="thead-dark chead">
          <tr>
            <th class="noc">DARI ratio:</th>
            <th>< 0.7x</th>
            <th>0.7x - 0.9x</th>
            <th colspan="2">0.9x - 1.1x</th>
            <th>1.1x - 1.30x</th>
            <th>> 1.30x</th>
          </tr>
          </thead>
          <tbody>
          <tr>
            <td class="noc">Blocks:</td>
            <td>{{ blocks_dari.vlow }}</td>
            <td>{{ blocks_dari.low }}</td>
            <td colspan="2">{{ blocks_dari.norm }}</td>
            <td>{{ blocks_dari.high }}</td>
            <td>{{ blocks_dari.vhigh }}</td>
          </tr>
          <tr>
            <td class="noc">Percentage:</td>
            <td>{{ dari_pct('vlow') }}</td>
            <td>{{ dari_pct('low') }}</td>
            <td colspan="2">{{ dari_pct('norm') }}</td>
            <td>{{ dari_pct('high') }}</td>
            <td>{{ dari_pct('vhigh') }}</td>
          </tr>
          <tr>
            <td></td>
            <td colspan="3" class="times-high">{{ sprintf('%.2f%%', dari_low / blocks * 100) }}</td>
            <td colspan="3" class="times-norm">{{ sprintf('%.2f%%', dari_high / blocks * 100) }}</td>
          </tr>
          </tbody>
        </table>
      </small>
    </p>
  </div>
</template>

<script>
  import moment from 'moment'
  import Chart from 'chart.js'
  import Vue from 'vue'
  import expected from './difficulty-expected.vue'

  export default {
    name: 'daa',
    data () {
      return {
        ts: null,
        blocks: 0,
        dari_low: 0,
        dari_high: 0,
        blocks_min: 0,
        blocks_max: 0,
        interval: 0,
        times: {
          '15': 0,
          '10': 0,
          '7': 0,
          '5': 0,
          '3': 0,
          '2': 0,
          '1': 0,
          '0': 0
        },
        blocks_dari: {
          'vlow': 0,
          'low': 0,
          'norm': 0,
          'high': 0,
          'vhigh': 0
        },
      }
    },
    mounted () {
      this.redraw()
    },
    computed: {
      block_times_high () {
        let t = this.times
        return this.sprintf('%.2f%%', (t['15'] + t['10'] + t['7']) / this.blocks * 100)
      },
      block_times_norm () {
        let t = this.times
        return this.sprintf('%.2f%%', (t['5']) / this.blocks * 100)
      },
      block_times_low () {
        let t = this.times
        return this.sprintf('%.2f%%', (t['3'] + t['2'] + t['1'] + t['0']) / this.blocks * 100)
      },
      correcting () {
        return this.blocksec > 600 ? 'overcorrecting' : 'undercorrecting'
      },
      correctpct () {
        return this.sprintf('%.2f%%', (this.blocksec > 600 ? this.blocksec / 600 - 1 : (600 / this.blocksec) - 1) * 100)
      },
      blocktime () {
        return this.sprintf('%.2f', this.blocksec / 60)
      },
      blocksec () {
        return (moment().unix() - this.ts) / this.blocks
      },
      date_from () {
        return moment.unix(this.ts).format('MMM. Do')
      },
      date_to () {
        return moment().format('MMM. Do')
      }
    },
    methods: {
      times_pct (k) {
        let val = this.times[k]
        return val > 0 ? this.sprintf('%.1f%%', val / this.blocks * 100) : '0%'
      },
      dari_pct (k) {
        let val = this.blocks_dari[k]
        return val > 0 ? this.sprintf('%.1f%%', val / this.blocks * 100) : '0%'
      },
      redraw () {

        this.blocks_dari = {
          'vlow': 0,
          'low': 0,
          'norm': 0,
          'high': 0,
          'vhigh': 0
        }

        this.times = {
          '15': 0,
          '10': 0,
          '7': 0,
          '5': 0,
          '3': 0,
          '2': 0,
          '1': 0,
          '0': 0
        }

        let ts = 1510602000
        let d = this.data.BTC.history['all']

        if (d[0].timestamp > ts) {
          ts = d[0].timestamp
        }

        let interval = d[1].timestamp - d[0].timestamp
        this.interval = interval

        this.ts = ts
        this.blocks = this.data.BTC.history['all'].reduce((blocks, obj) => {
          if (obj.timestamp > ts) {
            blocks += obj.blocks
          }
          return blocks
        }, 0)

        this.dari_low = this.data.BTC.history['all'].reduce((blocks, obj) => {
          if (obj.timestamp > ts) {

            let btc = this.ts_coin("BCH", obj.timestamp)
            if (btc.dari > obj.dari) {
              blocks += obj.blocks
            }


            let ratio = btc.dari / obj.dari
            if (ratio > 1.30) {
              this.blocks_dari.vlow += obj.blocks
            } else if (ratio > 1.1) {
              this.blocks_dari.low += obj.blocks
            } else if (ratio < 0.7) {
              this.blocks_dari.vhigh += obj.blocks
            } else if (ratio < 0.9) {
              this.blocks_dari.high += obj.blocks
            } else {
              this.blocks_dari.norm += obj.blocks
            }
          }
          return blocks
        }, 0)

        this.dari_high = this.data.BTC.history['all'].reduce((blocks, obj) => {
          if (obj.timestamp > ts) {

            let btc = this.ts_coin("BCH", obj.timestamp)
            if (btc.dari < obj.dari) {
              blocks += obj.blocks
            }
          }
          return blocks
        }, 0)

        this.blocks_min = this.data.BTC.history['all'].reduce((blocks, obj) => {
          if (obj.timestamp > ts) {
            if (obj.blocks < blocks) {
              blocks = obj.blocks
            }
          }
          return blocks
        }, 9999999)

        this.blocks_max = this.data.BTC.history['all'].reduce((blocks, obj) => {
          if (obj.timestamp > ts) {
            if (obj.blocks > blocks) {
              blocks = obj.blocks
            }
          }
          return blocks
        }, 0)

        this.data.BTC.history['all'].forEach((obj) => {
          if (obj.timestamp > ts) {
            if (obj.blocks === 0) {
              this.times['x'] += 1
            }

            let min = interval / obj.blocks / 60

            if (min < 4) {
              this.times['15'] += obj.blocks
            } else if (min < 6) {
              this.times['10'] += obj.blocks
            } else if (min < 7.5) {
              this.times['7'] += obj.blocks
            } else if (min < 15) {
              this.times['5'] += obj.blocks
            } else if (min < 20) {
              this.times['3'] += obj.blocks
            } else if (min < 30) {
              this.times['2'] += obj.blocks
            } else if (min < 60) {
              this.times['1'] += obj.blocks
            } else {
              this.times['0'] += obj.blocks
            }
          }
        })

        //let ts = d[Math.floor(d.length / 2)].timestamp

        let labels = this.data.BTC.history['all'].reduce((dates, obj) => {
          if (obj.timestamp > ts) {
            dates.push(moment.unix(obj.timestamp).format("MMM, DD HH:mm"))
          }
          return dates
        }, [])

        let act = (vals, obj) => {
          if (obj.timestamp > ts) {
            vals.push(obj.blocks > 0 ? this.sprintf('%.0f', obj.avg_diff) : null)
          }
          return vals
        }

        let perf = (vals, obj) => {
          if (obj.timestamp > ts) {
            //vals.push(obj.blocks > 0 ? this.sprintf('%.0f', obj.ratep) : null)

            let o = this.tsrate('BCH', obj.timestamp)
            if (!o) o = 0
            vals.push(o ? sprintf('%.2f', obj.ratep / (obj.ratep + o) * 100) : null)
          }
          return vals
        }

        let dari = (vals, obj) => {
          if (obj.timestamp > ts) {
            vals.push(obj.blocks > 0 ? this.sprintf('%.2f', obj.dari / this.ts_coin('BCH', obj.timestamp).dari) : null)
          }
          return vals
        }

        let blks = (vals, obj) => {
          if (obj.timestamp > ts) {
            //vals.push(obj.blocks > 0 ? this.sprintf('%.2f', (obj.rate3h / obj.avg_diff - 1)) : null)
            vals.push(this.sprintf('%.2f', (3600 / interval) * obj.blocks))
          }
          return vals
        }

        let blks_hi = (vals, obj) => {
          if (obj.timestamp > ts) {
            let rat =  obj.dari / this.ts_coin('BCH', obj.timestamp).dari
            vals.push(rat > 1 ? this.sprintf('%.2f', obj.blocks) : null)
          }
          return vals
        }

        let blks_lo = (vals, obj) => {
          if (obj.timestamp > ts) {
            let rat =  obj.dari / this.ts_coin('BCH', obj.timestamp).dari
            vals.push(rat <= 1 ? this.sprintf('%.2f', obj.blocks) : null)
          }
          return vals
        }



        this.add_dataset({
          type: 'line',
          label: 'Difficulty',
          borderColor: this.bch_c,
          data: this.data.BTC.history['all'].reduce(act, []),
          spanGaps: true,
          yAxisID: 'diff',
          borderDash: [3, 3],
        }, (d) => {
          d.lineTension = 0.5
        })


        /*        this.add_dataset({
                  type: 'line',
                  label: 'Perfect',
                  borderColor: '#ff0000',
                  data: this.data.BTC.history['all'].reduce(perf, []),
                  spanGaps: true,
                  yAxisID: 'diff'
                }, (d) => {
                  d.lineTension = 0.5
                })*/

        this.add_dataset({
          type: 'line',
          label: 'DARI ratio vs BCH',
          borderColor: '#ff0000',
          fill: false,
          borderWidth: 1,
          data: this.data.BTC.history['all'].reduce(dari, []),
          yAxisID: 'dari'
        }, (d) => {
          d.lineTension = 0.2
        })

        this.add_dataset({
          type: 'bar',
          label: 'Blocks/hr',
          backgroundColor: '#eeeeee',
          borderColor: '#eeeeee',
          fill: true,
          borderWidth: 0,
          data: this.data.BTC.history['all'].reduce(blks, []),
          yAxisID: 'ratio'
        })

        /*        this.add_dataset({
                  type: 'bar',
                  label: 'Blocks Lo',
                  backgroundColor: '#0000ff',
                  borderColor: '#eeeeee',
                  fill: true,
                  borderWidth: 0,
                  data: this.data.BTC.history['all'].reduce(blks_lo, []),
                  yAxisID: 'ratio'
                })

                this.add_dataset({
                  type: 'bar',
                  label: 'Blocks Hi',
                  backgroundColor: '#00ff00',
                  borderColor: '#eeeeee',
                  fill: true,
                  borderWidth: 0,
                  data: this.data.BTC.history['all'].reduce(blks_hi, []),
                  yAxisID: 'ratio'
                })*/

        var config = {
          type: 'bar',
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
                id: 'diff',
                display: false,
                position: 'right',
                scaleLabel: {
                  display: true,
                  labelString: 'Difficulty'
                },
                ticks: {
                  callback: (value, index, values) => {
                    return this.numeral('0,0', value)
                  },
                  min: 0,
                  maxTicksLimit: 10
                },
              }, {
                display: true,
                id: 'ratio',
                position: 'left',
                scaleLabel: {
                  display: true,
                  labelString: 'Blocks per hour'
                },
                gridLines: {
                  color: [null, null, null, null, null, null, null, null, '#666', null, null, null]
                },
                ticks: {
                  min: 0,
                  max: 30,
                  stepSize: 3,
                  callback: (val) => {
                    return val
                  }
                }
              },
                {
                  display: true,
                  id: 'dari',
                  position: 'right',
                  scaleLabel: {
                    display: true,
                    labelString: 'DARI ratio'
                  },
                  gridLines: {
                    color: [null, null, null, null, null, '#666', null, null, null]
                  },
                  ticks: {
                    min: 0,
                    max: 2
                  }
                }
              ]
            }
          }
        };

        this.draw('canvas-daa', config, labels)
      }
    }
  }
</script>

<style>
  .times-table {
    text-align: center;
  }

  thead.chead th {
    text-align: center;
  }

  .noc {
    text-align: left !important;
  }

  .times-high {
    background-color: rgba(255, 0, 0, 0.1);
  }

  .times-norm {
    background-color: rgba(0, 255, 59, 0.11);
  }

  .times-low {
    background-color: rgba(28, 114, 255, 0.11);
  }
</style>

