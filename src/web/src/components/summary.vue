<template>
  <div>

    <div class="container">
      <div class="row justify-content-md-center">
        <div class="col col-xl-7 col-lg-8 col-md-10 col-sm-12">

          <table class="table table-sm table-summary">
            <tbody>
            <tr>
              <th class="noboto"></th>
              <td class="noboto">BTC/USD</td>
              <td class="noboto">BCH/USD</td>
              <td class="noboto">BCH/BTC</td>
            </tr>
            <tr>
              <td>Now</td>
              <td v-for="coin in coins">
                <!--          <span v-html="price_change(coin, true)"></span> -->
                {{ money(sum(coin).price.now, false) }}
                (<span v-html="price_change(coin)"></span>)
              </td>
              <td>
                {{ sprintf('%.3f', sum('BCH').price.now / sum('BTC').price.now) }}
                (<span v-html="price_ratio(false)"></span>)
              </td>
            </tr>
            <tr>
              <td>
                -1d
              </td>
              <td v-for="coin in coins">
                {{ money(sum(coin).price.day_1, false) }}
              </td>
              <td>
                {{ sprintf('%.3f', sum('BCH').price.day_1 / sum('BTC').price.day_1) }}
              </td>
            </tr>

            <tr>
              <td>
                -7d
              </td>
              <td v-for="coin in coins">
                {{ money(sum(coin).price.day_7, false) }}
              </td>
              <td>
                {{ sprintf('%.3f', sum('BCH').price.day_7 / sum('BTC').price.day_7) }}
              </td>
            </tr>

            <tr>
              <td>
                -30d
              </td>
              <td v-for="coin in coins">
                {{ money(sum(coin).price.day_30, false) }}
              </td>
              <td>
                {{ sprintf('%.3f', sum('BCH').price.day_30 / sum('BTC').price.day_30) }}
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <rdari></rdari>


    <div class="container">
      <div class="row justify-content-md-center">
        <div class="col col-xl-7 col-lg-8 col-md-10">

          <table class="table table-sm table-summary">
            <tbody>

            <tr>
              <td class="noboto"></td>
              <td class="noboto" align="center"><strong>BTC</strong></td>
              <td class="noboto" align="center"><strong>BCH</strong></td>
            </tr>

            <tr>
              <th colspan="3" class="summary-header">
                <span class="toggle-stat" @click="toggle('dari')">
                  <i v-if="showstats.dari" class="fa fa-chevron-down"></i>
                  <i v-if="!showstats.dari" class="fa fa-chevron-up"></i>
                  DARI ratio
                  <small>(period averages)</small>
                </span>
                <router-link class="chartlink" :to="{name: 'reward.dari'}">chart</router-link>
              </th>
            </tr>

            <template v-if="showstats.dari">
              <tr>
                <td>Now*</td>
                <td v-for="coin in coins">
                  <strong>
                    <span v-html="dari_rel(coin, 'blk_3')"></span>
                  </strong>
                </td>
              </tr>

              <tr>
                <td>6 hour</td>
                <td v-for="coin in coins">
                  <span v-html="dari_rel(coin, 'hour_6')"></span>
                </td>
              </tr>

              <tr>
                <td>1 day</td>
                <td v-for="coin in coins">
                  <span v-html="dari_rel(coin, 'day_1')"></span>
                </td>
              </tr>

              <tr>
                <td>7 day</td>
                <td v-for="coin in coins">
                  <span v-html="dari_rel(coin, 'day_7')"></span>
                </td>
              </tr>

              <tr>
                <td colspan="3">
                  <small>
                    * The average of the last 3 blocks is used to even out tx fees.
                  </small>
                </td>
              </tr>
            </template>

            <tr>
              <th colspan="3" class="summary-header">
                <span class="toggle-stat" @click="toggle('hashrate')">
                  <i v-if="showstats.hashrate" class="fa fa-chevron-down"></i>
                  <i v-if="!showstats.hashrate" class="fa fa-chevron-up"></i>
                  Relative hashrate
                  <small>(period averages)</small>
                </span>
                <router-link class="chartlink" :to="{name: 'pow.hashrate'}">chart</router-link>
              </th>
            </tr>
            <template v-if="showstats.hashrate">

              <tr>
                <td>6 hours</td>
                <td v-for="coin in coins">
                  {{ hashrate_rel(coin, 'hour_6') }}
                </td>
              </tr>
              <tr>
                <td>1 day</td>
                <td v-for="coin in coins">
                  {{ hashrate_rel(coin, 'day_1') }}
                </td>
              </tr>
              <tr>
                <td>7 day</td>
                <td v-for="coin in coins">
                  {{ hashrate_rel(coin, 'day_7') }}
                </td>
              </tr>
              <tr>
                <td>30 day</td>
                <td v-for="coin in coins">
                  {{ hashrate_rel(coin, 'day_30') }}
                </td>
              </tr>

            </template>


            <tr>
              <th colspan="3" class="summary-header">

                <span class="toggle-stat" @click="toggle('blockrate')">
                  <i v-if="showstats.blockrate" class="fa fa-chevron-down"></i>
                  <i v-if="!showstats.blockrate" class="fa fa-chevron-up"></i>
                  Blocks found
                  <small>(period sum/rate)</small>
                </span>
                <router-link class="chartlink" :to="{name: 'blocks.time'}">chart</router-link>
              </th>
            </tr>
            <template v-if="showstats.blockrate">
              <tr>
                <td>Last hour</td>
                <td v-for="coin in coins">
                  {{ sum(coin).blocks.hour_1 }}
                </td>
              </tr>
              <tr>
                <td>6 hours</td>
                <td v-for="coin in coins">
                  {{ sprintf('%.0f', sum(coin).blocks.hour_6) }} ({{ sprintf('%.2f', sum(coin).blocks.per_hour.hour_6)
                  }}/hr)
                </td>
              </tr>
              <tr>
                <td>24 hours</td>
                <td v-for="coin in coins">
                  {{ sprintf('%.0f', sum(coin).blocks.day_1) }} ({{ sprintf('%.2f', sum(coin).blocks.per_hour.day_1)
                  }}/hr)
                </td>
              </tr>
              <tr>
                <td>2 weeks</td>
                <td v-for="coin in coins">
                  {{ sprintf('%.0f', sum(coin).blocks.day_14) }} ({{ sprintf('%.2f', sum(coin).blocks.per_hour.day_14)
                  }}/hr)
                </td>
              </tr>
            </template>


            <tr>
              <th colspan="3" class="summary-header">
                <span class="toggle-stat" @click="toggle('txs')">
                  <i v-if="showstats.txs" class="fa fa-chevron-down"></i>
                  <i v-if="!showstats.txs" class="fa fa-chevron-up"></i>
                  Confirmed transactions
                  <small>(24h averages)</small>
                </span>
                <router-link class="chartlink" :to="{name: 'tx.txs'}">chart</router-link>
              </th>
            </tr>
            <template v-if="showstats.txs">
              <tr>
                <td>Per block</td>
                <td v-for="coin in coins">{{ sprintf('%.0f', sum(coin).txs.block) }}</td>
              </tr>
              <tr>
                <td>Per second</td>
                <td v-for="coin in coins">{{ sprintf('%.2f', sum(coin).txs.second) }}</td>
              </tr>
              <tr>
                <td>Per hour</td>
                <td v-for="coin in coins">{{ sprintf('%.0f', sum(coin).txs.hour) }}</td>
              </tr>
              <tr>
                <td>Per day</td>
                <td v-for="coin in coins">{{ sprintf('%.0f', sum(coin).txs.day) }}</td>
              </tr>
            </template>

            <tr>
              <th colspan="3" class="summary-header">
                <span class="toggle-stat" @click="toggle('txfee')">
                  <i v-if="showstats.txfee" class="fa fa-chevron-down"></i>
                  <i v-if="!showstats.txfee" class="fa fa-chevron-up"></i>
                  Transaction fees
                  <small>(24h averages)</small>
                </span>
                <router-link class="chartlink" :to="{name: 'tx.fee'}">chart</router-link>
              </th>
            </tr>
            <template v-if="showstats.txfee">
              <tr>
                <td>Sat/B</td>
                <td v-for="coin in coins">{{ sprintf('%.2f', sum(coin).tx_fees.sat_b) }}</td>
              </tr>
              <tr>
                <td>USD/kB</td>
                <td v-for="coin in coins">{{ money(sum(coin).tx_fees.usd_kb) }}</td>
              </tr>
            </template>
            <!--        <tr>
                      <td>Normalized*</td>
                      <td v-for="coin in coins">{{ fee_normalized(coin) }}</td>
                    </tr>
                    <tr>
                      <td colspan="3">
                        <small>
                          * USD/kB if both coins were trading at $5,000. For perspective only, actual fees will depend on demand.
                        </small>
                      </td>
                    </tr>-->

            <tr>
              <th colspan="3" class="summary-header">
            <span class="toggle-stat" @click="toggle('reward')">
              <i v-if="showstats.reward" class="fa fa-chevron-down"></i>
              <i v-if="!showstats.reward" class="fa fa-chevron-up"></i>
              Average block reward
            </span>
                <small>
                  (<span v-bind:class="{'rpa': rewardperiod=='block'}"
                         @click="rewardperiod = 'block'">last block</span>,
                  <span v-bind:class="{'rpa': rewardperiod=='day'}" @click="rewardperiod = 'day'">past day</span>,
                  <span v-bind:class="{'rpa': rewardperiod=='week'}" @click="rewardperiod = 'week'">past week</span>)
                </small>

                <router-link class="chartlink" :to="{name: 'reward.blocks'}">chart</router-link>
              </th>
            </tr>
            <template v-if="showstats.reward">
              <tr>
                <td>Coinbase</td>
                <td v-for="coin in coins">
                  {{ money(rp(coin).coinbase / rpb(coin), false) }}
                </td>
              </tr>
              <tr>
                <td>Fees</td>
                <td v-for="coin in coins">
                  {{ money(rp(coin).fees / rpb(coin), false) }}
                  <small>({{ sprintf('%.2f', rp(coin).fees_coin / rpb(coin)) }} {{ coin }})</small>
                </td>
              </tr>
              <tr>
                <td>Total</td>
                <td v-for="coin in coins">
                  {{ money(rp(coin).total / rpb(coin), false) }}
                </td>
              </tr>
              <tr>
                <td>Fee %</td>
                <td v-for="coin in coins">
                  {{ sprintf('%.2f%%', (rp(coin).fee_pct)) }}
                </td>
              </tr>
              <tr>
                <td>DARI</td>
                <td v-for="coin in coins">
                  <span v-show="rewardperiod === 'block'">{{ sprintf('%.2f', (sum(coin).dari.last)) }}</span>
                  <span v-show="rewardperiod === 'day'">{{ sprintf('%.2f', (sum(coin).dari.day_1)) }}</span>
                  <span v-show="rewardperiod === 'week'">{{ sprintf('%.2f', (sum(coin).dari.day_7)) }}</span>
                </td>
              </tr>
            </template>


            <tr>
              <th colspan="3" class="summary-header">
                <span class="toggle-stat" @click="toggle('blocksize')">
                  <i v-if="showstats.blocksize" class="fa fa-chevron-down"></i>
                  <i v-if="!showstats.blocksize" class="fa fa-chevron-up"></i>
                  Average blocksize in kB
                  <small>(excluding empty blocks)</small>
                </span>
                <router-link class="chartlink" :to="{name: 'blocks.size'}">chart</router-link>
              </th>
            </tr>
            <template v-if="showstats.blocksize">
              <tr>
                <td>Today</td>
                <td v-for="coin in coins">{{ sprintf('%.2f', sum(coin).blocksize.day_1) }}</td>
              </tr>

              <tr>
                <td>-7 day</td>
                <td v-for="coin in coins">{{ sprintf('%.2f', sum(coin).blocksize.day_7) }}</td>
              </tr>

              <tr>
                <td>-30 day</td>
                <td v-for="coin in coins">{{ sprintf('%.2f', sum(coin).blocksize.day_30) }}</td>
              </tr>
            </template>


            <tr>
              <th colspan="3" class="summary-header">
                <span class="toggle-stat" @click="toggle('retarget')">
                  <i v-if="showstats.retarget" class="fa fa-chevron-down"></i>
                  <i v-if="!showstats.retarget" class="fa fa-chevron-up"></i>
                  Difficulty retarget
                </span>
                <router-link class="chartlink" :to="{name: 'pow.retarget'}">chart</router-link>
              </th>
            </tr>
            <template v-if="showstats.retarget">
              <tr>
                <td>Date (24h UTC)</td>
                <td>{{ date('BTC') }}</td>
                <td>-</td>
              </tr>
              <tr>
                <td>Remaining</td>
                <td>
                  {{ days('BTC') }}
                  <small>({{ left('BTC') }} blk)</small>
                </td>
                <td>-</td>
              </tr>
              <tr>
                <td>Change</td>
                <td>
                  <i v-if="next_pct('BTC', true) > 1" class="fa fa-chevron-up green"></i>
                  <i v-if="next_pct('BTC', true) < 1" class="fa fa-chevron-down red"></i>
                  {{ next('BTC') }}
                </td>
                <td>-</td>
              </tr>
            </template>

            <tr>
              <th colspan="3" class="summary-header">
                <span class="toggle-stat" @click="toggle('halvening')">
                  <i v-if="showstats.halvening" class="fa fa-chevron-down"></i>
                  <i v-if="!showstats.halvening" class="fa fa-chevron-up"></i>
                  Halvening
                  <small>(based on blocks past 30 days)</small>
                </span>
              </th>
            </tr>
            <template v-if="showstats.halvening">
              <tr>
                <td>Blocks</td>
                <td v-for="coin in coins">{{ sum(coin).halvening.blocks }}</td>
              </tr>
              <tr>
                <td>Days</td>
                <td v-for="coin in coins">{{ sprintf('%.0f', sum(coin).halvening.seconds / (24 * 3600)) }}</td>
              </tr>
              <tr>
                <td>New reward*</td>
                <td v-for="coin in coins">{{ halvening_reward(coin) }}</td>
              </tr>
              <tr>
                <td>New fee %**</td>
                <td v-for="coin in coins">{{ halvening_fees(coin) }}</td>
              </tr>
              <tr>
                <td colspan="3">
                  <small>
                    * Based on 30d average coinbase+fees with current exchange rate ** Income from fees as a percentage of the total block reward
                  </small>
                </td>
              </tr>
            </template>

            </tbody>
          </table>
        </div>
      </div>
    </div>

    <small>
      <template v-if="showstats.halvening">
        <feemarket></feemarket>
      </template>
    </small>


  </div>
</template>

<script>
  import dari from './charts/dari.vue'
  import humanize from 'humanize-duration'
  import feemarket from './charts/feemarket.vue'
  import rdari from './charts/rdari.vue'
  import axios from 'axios'

  export default {
    name: 'summary',
    components: {
      feemarket,
      rdari
    },
    data () {
      return {
        chart: null,
        rewardperiod: 'day',
        showstats: {
          'dari': true,
          'hashrate': true,
          'blockrate': true,
          'txs': true,
          'txfee': true,
          'reward': true,
          'blocksize': true,
          'retarget': true,
          'halvening': true,
        },
        eri: null
      }
    },
    mounted () {
      if (localStorage.getItem('showstats.enabled') === '1') {
        let show = localStorage.getItem('showstats');
        this.showstats = JSON.parse(show)
      }
      let rp = localStorage.getItem('rewardperiod')
      if (rp) this.rewardperiod = rp
      this.startExchangerates()
    },
    destroyed () {
      clearInterval(this.eri)
    },
    computed: {
      coins () {
        return ['BTC', 'BCH']
      },
      visible () {
        return this.$store.getters['visible']
      }
    },
    watch: {
      rewardperiod (n, o) {
        localStorage.setItem('rewardperiod', n);
      },
      visible (n, o) {
        if (n === true) {
          this.startExchangerates()
        } else {
          clearInterval(this.eri)
        }
      }
    },
    methods: {
      startExchangerates () {
        this.eri = setInterval(this.exchangerates, 10000)
      },
      exchangerates () {
        axios.get(process.env.API + 'exchangerate').then((resp) => {
          this.data.BTC.summary.price.now = resp.data.btc
          this.data.BCH.summary.price.now = resp.data.bch
        }).catch((err) => {
          console.log(err)
        })
      },
      sum (coin) {
        return this.data[coin].summary
      },
      toggle (stat, e) {
        this.showstats[stat] = !this.showstats[stat]
        localStorage.setItem('showstats', JSON.stringify(this.showstats));
        localStorage.setItem('showstats.enabled', '1');
      },
      hashrate_rel (coin, d) {
        let s = this.sum(coin).hashrate[d]
        let o = this.sum(this.coino(coin)).hashrate[d]
        let t = s + o

        return this.sprintf('%.2f%%', s / t * 100)
      },
      dari_rel (coin, d) {
        let a = this.data[coin].summary.dari[d]
        let b = this.data[this.coino(coin)].summary.dari[d]

        if (a > b) {
          return `<span class="hidari">${this.sprintf('%.2fx', a / b)}</span>`
        } else {
          return `<span class="lodari">1.00x</span>`
        }
      },
      fee_normalized (coin) {
        let p = 5000 / this.sum(coin).price.now
        return this.money(this.sum(coin).tx_fees.usd_kb * p)
      },
      halvening_reward (coin) {

        let n = this.sum(coin).halvening.reward_new
        let o = this.sum(coin).halvening.reward

        let change = ((n / o) - 1) * 100
        let formatted = this.sprintf('%.2f%%', change)
        let reward = this.money(n * this.sum(coin).price.now, false)

        return `${reward} (${formatted})`
      },
      halvening_fees (coin) {
        let n = (this.sum(coin).halvening.reward_new - 6.25) / this.sum(coin).halvening.reward_new
        let o = (this.sum(coin).halvening.reward - 12.5) / this.sum(coin).halvening.reward
        let pct = this.sprintf('%.2f%%', n * 100)

        return `${pct}`
      },
      price_change (coin, icon = false) {
        let n = this.sum(coin).price.now
        let o = this.sum(coin).price.day_1

        let change = ((n / o) - 1) * 100

        let formatted = this.sprintf('%.2f%%', change)
        formatted = change > 0 ? '+' + formatted : formatted


        if (change > 0) {
          return icon ? `<i class="fa fa-chevron-up green"></i>` : `<span class="green">${formatted}</span>`
        } else {
          return icon ? `<i class="fa fa-chevron-down red"></i>` : `<span class="red">${formatted}</span>`
        }
      },
      price_ratio (icon = false) {
        let n = this.sum('BCH').price.now / this.sum('BTC').price.now
        let o = this.sum('BCH').price.day_1 / this.sum('BTC').price.day_1

        let change = ((n / o) - 1) * 100

        let formatted = this.sprintf('%.2f%%', change)
        formatted = change > 0 ? '+' + formatted : formatted

        if (change > 0) {
          return icon ? `<i class="fa fa-chevron-up green"></i>` : `<span class="green">${formatted}</span>`
        } else {
          return icon ? `<i class="fa fa-chevron-down red"></i>` : `<span class="red">${formatted}</span>`
        }
      },
      rp (coin) { // active reward period
        return this.sum(coin).reward[this.rewardperiod]
      },
      rpb (coin) { // active reward period blocks
        return this.rp(coin).blocks
      }
    },
  }
</script>

<style>
  .red {
    color: red;
  }

  .green {
    color: green;
  }

  table.nomargin {
    margin-bottom: 0;
  }

  table.table-summary td {
    text-align: center;
  }

  table.table-summary td:first-child {
    text-align: left;
  }

  .summary-header {
    /*    text-align: center;*/
  }

  th.noboto, td.noboto {
    border-top: 0;
  }

  span.rpa {
    text-decoration: underline;
  }

  .chartlink {
    float: right;
  }

  .toggle-stat {
    cursor: pointer;
  }

  .hidari {
    color: green;
  }

  .lodari {
  }
</style>
