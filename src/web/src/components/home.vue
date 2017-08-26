<template>
    <div>
        <p>
            Mining profitability based on fees, exchange rate and mining difficulty. 
        </p>
        <div class="container">
            <div class="row justify-content-md-center">
                <div class="col col-xl-7 col-lg-8 col-md-8">

                    <table class="table table-sm">
                        <thead>
                            <tr class="trfirst">
                                <th></th>
                                <th>current*
                                    <span class="mini">(expected**)</span>
                                </th>
                                <th>6h<sup>avg</sup>
                                </th>
                                <th>1d<sup>avg</sup>
                                </th>
                                <th>7d<sup>avg</sup>
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-bind:class="{'bgdarilead': leader == 'BTC'}">
                                <td>
                                    <strong>BTC</strong>
                                </td>
                                <td>
                                    <span v-html="dari_period('BTC', 'BCH', 'last', true)"></span>
                                    <small v-html="next_dari_coin('BTC', 'BCH')"></small>
                                </td>
                                <td v-html="dari_period('BTC', 'BCH', '6h')"></td>
                                <td v-html="dari_period('BTC', 'BCH', '1d')"></td>
                                <td v-html="dari_period('BTC', 'BCH', '7d')"></td>
                            </tr>
                            <tr v-bind:class="{'bgdarilead': leader == 'BCH'}">
                                <td>
                                    <strong>BCH</strong>
                                </td>

                                <td>
                                    <span v-html="dari_period('BCH', 'BTC', 'last', true)"></span>
                                    <small v-html="next_dari_coin('BCH', 'BTC')"></small>
                                </td>
                                <td v-html="dari_period('BCH', 'BTC', '6h')"></td>
                                <td v-html="dari_period('BCH', 'BTC', '1d')"></td>
                                <td v-html="dari_period('BCH', 'BTC', '7d')"></td>
                            </tr>
                        </tbody>
                    </table>
                    <p class="text-center">
                        <small>
                            <template v-if="!empty_any">
                                * Based on the last block and exchange rates {{ money(data.BTC.averages.last.price) }}<sup>BTC</sup> {{ money(data.BCH.averages.last.price) }}<sup>BCH</sup>
                            </template>
                            <template v-else>
                                * Empty BTC or BCH block found, skipping calculations.
                            </template>
                            <br>
                            ** After the {{ first_retarget }} difficulty adjustment 
                            <span v-if="first_retarget">
                                (~{{ time_retarget }}/{{ left }} blocks)
                            </span>
                        </small>
                    </p>
                    <p class="text-center">
                        <small>
                        </small>
                    </p>
                </div>
            </div>
        </div>

        <dari></dari>
        <br>
        <table class="table table-sm">
            <thead>
                <tr>
                    <th>Recent changes</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>
                        Added <router-link to="pow/difficulty#exp">expected vs actual</router-link> difficulty chart.
                    </td>
                </tr>
                <tr>
                    <td>
                        Added <router-link :to="{name: 'pow.work'}">total work</router-link> and <router-link :to="{name: 'blocks.height'}">block height</router-link> charts.
                    </td>
                </tr>
                <tr>
                    <td>
                        Added <router-link :to="{name: 'reward.inflation'}">inflation</router-link> chart.
                    </td>
                </tr>
                <tr>
                    <td>
                        Many have asked so here it is:
                        <router-link :to="{name: 'contact'}">donations</router-link>.
                    </td>
                </tr>
                <tr>
                    <td>
                        Added
                        <router-link :to="{name: 'pow.hashrateabs'}">absolute hashrate</router-link> chart.
                    </td>
                </tr>
                <tr>
                    <td>
                        Added long polling for live data. No more need to refresh page.
                    </td>
                </tr>
                <tr>
                    <td>
                        Added
                        <router-link :to="{name: 'pow.retarget'}">difficulty retarget</router-link> chart and estimates.
                    </td>
                </tr>
            </tbody>
        </table>

    </div>
</template>

<script>
import dari from './charts/dari.vue'
import humanize from 'humanize-duration'

export default {
    name: 'home',
    components: {
        dari
    },
    data() {
        return {
            chart: null,
            last_dari: {'BTC': 1, 'BCH': 1}
        }
    },
    computed: {
        first_retarget(){
            let at = this.estimate('BTC', this.data['BTC'].averages.last, true)
            let bt = this.estimate('BCH', this.data['BCH'].averages.last, true)

            if(!at || !bt) return null

            return at < bt ? 'BTC' : 'BCH'
        },
        time_retarget() {
            let coin = this.first_retarget
            if(coin === null) return null
            let mins = this.estimate(coin, this.data[coin].averages.last, true)
            return humanize(mins * 60 * 1000, { largest: 2 })
        },
        left() {
            let l = this.first_retarget
            if(l === null) return null
            let h = this.data[l].averages.last.height
            return (2016 - h % 2016)
        },
        empty_any() {
            return false //this.data.BTC.averages.last.txs === 1 || this.data.BCH.averages.last.txs === 1
        }
    },
    methods: {
        dari_period(coin, coin2, p, strong = false) {
            let a = this.data[coin].averages[p]
            let b = this.data[coin2].averages[p]
            if(a.blocks < 1 || b.blocks < 1) return null

            let pct = 0
            let space = '______'

            if(a.txs === 1 || b.txs === 1) {
                pct = this.last_dari[coin]
            } else {
                pct = this.data[coin].averages[p].dari / this.data[coin2].averages[p].dari
                if (pct < 1) pct = 1
                this.last_dari[coin] = pct
            }

            let res = sprintf('%.2f', pct)

            if (strong) {
                return `<span class="${pct > 1 ? 'advantage' : 'disadvantage_x'}">${pct > 0 ? res + 'x' : space}</span>`
            }

            return `<span class="${pct > 1 ? 'advantage_sm' : 'disadvantage_sm_x'}">${pct > 0 ? res + 'x' : space}</span>`
        },
        next_dari_coin(coin, coin2) {
            let at = this.estimate(coin, this.data[coin].averages.last, true)
            let bt = this.estimate(coin2, this.data[coin2].averages.last, true)

            let a = this.next_dari(coin, this.data[coin].averages.last)
            let b = this.next_dari(coin2, this.data[coin2].averages.last)

            let tbd = `<span class="mini">(<span class="disadvantage_sm_x">T.B.D.</span>)</span>`

            if (!at || !bt || !a || !b) {
                return tbd
            }

            if(at < bt) {
                b = this.data[coin2].averages.last.dari
            } else {
                a = this.data[coin].averages.last.dari
            }

            let pct = a / b
            if (pct < 1) pct = 1

            let res = this.sprintf('%.2f', pct)
            let space = '______'
            return `<span class="mini">(<span class="${pct > 1 ? 'advantage_sm' : 'disadvantage_sm_x'}">${pct > 0 ? res + 'x' : space}</span>)</span>`
        },
        next_dari(coin, obj) {
            let pct = this.next_pct(coin, true)
            return this.sprintf('%.3f', obj.dari * pct)
        },
        estimate(coin, obj, m = false) {
            let prev = this.data[coin].retargets.slice(-1)[0]
            let last = obj

            if (prev === undefined) prev = this.preforktarget

            let mod = last.height % 2016
            if(mod === 0 || mod === 1) return null


            let left = 2016 - mod
            let blkhr = obj.rate12h / obj.avg_diff * 6

            let hr = left / blkhr

            let days = hr / 24
            let mins = hr * 60

            return m === false ? this.sprintf('%.2f', days) : mins
        },
    },
    mounted() {

    },
}
</script>

<style>
.trfirst th {
    border-top: none;
}

.advantage {
    font-weight: bold;
    color: green;
}

.disadvantage {
    font-weight: bold;
    color: red;
}

.disadvantage_x {
    font-weight: bold;
}

.advantage_sm {
    color: green;
}

.disadvantage_sm {
    color: red;
}

.bgdarilead {
    background-color: rgba(58, 201, 56, 0.1);
}

th>sup {
    font-weight: normal;
}

.mini {
    font-size: 0.75rem;
}

@media (max-width: 575px) {
    .canvas-wrapper {
        height: 60%;
    }
}
</style>
