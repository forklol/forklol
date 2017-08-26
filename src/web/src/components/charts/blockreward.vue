
<template>
    <div>
        <p>
            Average block reward (coinbase+fees) in USD.
        </p>
        <div style="height:512px; width: 100%;">
            <canvas id="canvas-reward" style="height: 512px;"></canvas>
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
    name: 'blocks-reward',
    computed: {
        reward_ratio() {
            let data = this.data
            let leader = this.leader
            let loser = this.loser
            return sprintf('%.2f', (data[leader].reward_avg / 100000000 * data[leader].price) / (data[loser].reward_avg / 100000000 * data[loser].price))
        },
        stats() {
            return {
                header: ['3h', '6h', '12h', '1d', '7d'],
                f: (coin, avg) => {
                    return this.money(this.reward(this.data[coin].averages[avg]))
                },
                coins: [{
                    coin: 'BTC',
                    values: ['3h', '6h', '12h', '1d', '7d']
                }, {
                    coin: 'BCH',
                    values: ['3h', '6h', '12h', '1d', '7d']
                }]
            }
        }
    },
    mounted() {
        this.redraw()
    },
    methods: {
        reward(obj) {
            return obj.blocks > 0 && obj.price > 1 ? this.sprintf('%.2f', obj.reward_avg / 100000000 * obj.price) : null
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                backgroundColor: this.btc_cbg,
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => this.reward(obj)),
                fill: true,
            })

            this.add_dataset({
                label: 'BCH',
                backgroundColor: this.bch_cbg,
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.reward(obj)),
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
                                labelString: 'Blockreward (USD)'
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

            this.draw('canvas-reward', config, labels)
        },
    }
}
</script>

<style>

</style>
