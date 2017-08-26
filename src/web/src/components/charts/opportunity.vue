
<template>
    <div>
        <small>
            <div class="alert alert-info">
                Once segwit activates on the BTC chain this chart cannot be constructed in a reliable manner anymore (it already is a bit iffy). The chains wil diverge too much to make a sensible comparison between the two. Once that happens this chart will be removed.
            </div>
        </small>
        <p>
            The sum of potentially missed rewards by BCH miners compared to what they could earn on the BTC chain. Combined, BCH miners have currently missed out on an estimated total of {{ money(total_cap) }}.
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-opportunity" style="height: 512px;"></canvas>
        </div>
        <br>
        <br>
        <p>
            These numbers are based on blocks that were actually found and not only an estimation based on difficulty/hashrate. The hashrate (or
            <router-link :to="{name: 'pow.speed'}">speed</router-link>) is accounted for when calculating the difference. There is also a small compensation applied to fees to correct for loss of fee pressure. These numbers are not perfect, only an aproximation but they come close to reality.
            <br>
            <br> Also, this chart does not take into account the future gains in exchange rate of the coins. It assumes miners will sell the coins as soon as possible (hodling is not factored in ;)).
        </p>
    </div>
</template>
<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'

export default {
    name: 'reward-opportunity',
    data() {
        return {
            total: 0.0,
            zero: false
        }
    },
    computed: {
        time() {
            return this.BTC.history['all'][1].timestamp - this.BTC.history['all'][0].timestamp
        },
        lg() {
            return this.total > 0 ? 'gained' : 'lost'
        },
        total_cap() {
            if (this.zero) return 0
            return this.total < 0 ? '0.00' : this.total
        }
    },
    mounted() {
        this.redraw()
    },
    methods: {
        redraw() {
            this.total = 0.0
            this.zero = false

            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BCH',
                backgroundColor: this.bch_cbg,
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.balance(obj)),
                fill: true,
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
                                labelString: 'Opportunity cost (USD)'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return this.money(value)
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-opportunity', config, labels)
        },

        balance(obj) {
            if (this.zero) return 0
            let btc = this.find_ts('BTC', obj.timestamp)
            let bch = this.find_ts('BCH', obj.timestamp)

            if (!btc || !bch || btc.blocks < 1 || bch.blocks < 1) {
                return this.total
            }

            let thr = btc.rate + bch.rate
            let rpct = bch.rate / thr

            let btc_speed = btc.rate / btc.avg_diff
            let btc_speed_total = (btc.rate + bch.rate) / btc.avg_diff
            let speed_diff = btc_speed_total - btc_speed

            let new_reward = (btc.reward_avg - ((btc.reward_avg - 1250000000) * (3 * speed_diff)))

            let btc_total = new_reward / 100000000 * (bch.blocks * (bch.avg_diff / btc.avg_diff)) * btc.price
            let bch_total = bch.reward_avg / 100000000 * bch.blocks * bch.price

            this.total += (btc_total - bch_total)

            if (this.total < 0) {
                this.zero = true
                this.total = 0
            }

            return this.sprintf('%.2f', this.total)
        }
    }
}
</script>

<style>

</style>
