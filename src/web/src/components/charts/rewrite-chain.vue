
<template>
    <div>
        <p>
            Proof of work equivalent days to rewrite the complete BTC or BCH chain with 100% of available hashrate. Basically the same idea as the "fork rewrite" chart only for the entire blockchain.
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-attack-full" style="height: 512px;"></canvas>
        </div>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'

export default {
    name: 'rewrite-chain',
    computed: {
        replace_chain() {
            return moment.duration({
                minutes: this.data.BCH.pow.split.work / this.data.BTC.pow.current.work_actual * 10
            }).humanize()
        }
    },
    mounted() {
        this.redraw()
    },
    methods: {
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
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            var rate = this.data['BTC'].pow.current.work_actual + this.data['BCH'].pow.current.work_actual

            let attack = (obj) => {
                if (obj.blocks < 1) {
                    return null
                }

                let rate = this.rate(obj.timestamp)
                if (rate === null) {
                    return null
                }

                return this.sprintf('%.2f', (obj.work) / rate * 10 / 60 / 24)
            }

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => attack(obj)),
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => attack(obj)),
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

            this.draw('canvas-attack-full', config, labels)
        }
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
