
<template>
    <div>
        <p>
            Proof of work equivalent days to rewrite the BCH chain from moment of fork. A 3 day average is used to determine hashrate.
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-attack" style="height: 512px;"></canvas>
        </div>
        <br>
        <p>
            The ratio of total work since the split divided by estimate of hashrates at that time. The 5%, 10% and 25% are percentages of global hashpower. For example, it will take about {{ replace_chain }} for hashpower equivalent to that of BTC to rewrite the BCH chain from moment of fork until now. For a miner with 25% hashpower it will take about {{ replace_chain_rate(0.25) }}.
        </p>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'

export default {
    name: 'rewrite-fork',
    computed: {
        replace_chain() {
            return moment.duration({
                minutes: this.data.BCH.pow.split.work / this.data.BTC.pow.current.work_actual * 10
            }).humanize()
        },
    },
    mounted() {
        this.redraw()
    },
    methods: {
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            var prework = 32729585000856628.00

            var rate = this.data['BTC'].pow.current.work_actual + this.data['BCH'].pow.current.work_actual
            let attack = (obj, r) => {
                if (obj.blocks < 1) return null
                if (!r) return null
                return this.sprintf('%.2f', (obj.work - prework) / r * 10 / 60 / 24)
            }

            this.add_dataset({
                label: '5%',
                borderWidth: 1,
                borderColor: '#42cbf4',
                spanGaps: true,
                data: this.data.BCH.history['all'].map((obj) => {
                    return attack(obj, this.rate(obj.timestamp) * 0.05)
                }),
            })

            this.add_dataset({
                label: '10%',
                borderWidth: 1,
                borderColor: '#ad5eed',
                spanGaps: true,
                data: this.data.BCH.history['all'].map((obj) => {
                    return attack(obj, this.rate(obj.timestamp) * 0.10)
                }),
            })

            this.add_dataset({
                label: '25%',
                borderWidth: 1,
                borderColor: '#ff0000',
                spanGaps: true,
                data: this.data.BCH.history['all'].map((obj) => {
                    return attack(obj, this.rate(obj.timestamp) * 0.25)
                }),
            })

            this.add_dataset({
                label: 'BTC hashrate',
                borderWidth: 1,
                borderColor: this.btc_c,
                spanGaps: true,
                data: this.data.BCH.history['all'].map((obj) => {
                    return attack(obj, this.tsrate('BTC', obj.timestamp))
                }),
            })

            this.add_dataset({
                label: 'BCH hashrate',
                borderWidth: 1,
                borderColor: this.bch_c,
                spanGaps: true,
                data: this.data.BCH.history['all'].map((obj) => {
                    return attack(obj, this.tsrate('BCH', obj.timestamp))
                }),
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
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-attack', config, labels)
        },

        replace_chain_rate(rate) {
            return moment.duration({
                minutes: this.data.BCH.pow.split.work / this.global_rate(rate) * 10
            }).humanize()
        },
        global_rate(pct) {
            return (this.data['BTC'].pow.current.work_actual + this.data['BCH'].pow.current.work_actual) * pct
        },
        rate(ts) {
            let btc = this.data['BTC'].history['all'].find(h => h.timestamp === ts)
            let bch = this.data['BCH'].history['all'].find(h => h.timestamp === ts)

            if (!btc || !bch) {
                return null
            }

            return btc.rate + bch.rate
        },
        tsrate(coin, ts) {
            let o = this.data[coin].history['all'].find(h => h.timestamp === ts)
            return o ? o.rate : null
        },
    }
}
</script>

<style scoped>
@media (max-width: 575px) {
    .canvas-wrapper {
        height: 60%;
    }
}
</style>
