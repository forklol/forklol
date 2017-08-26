
<template>
    <div>
        <p>
            Average amount of fees per block in USD.
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-fees" style="height: 512px;"></canvas>
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
    name: 'reward-fees',
    computed: {
        feediff() {
            let btc = this.lastfee('BTC')
            let bch = this.lastfee('BCH')

            return btc / bch
        }
    },
    mounted() {
        this.redraw()
    },
    computed: {
        stats() {
            return {
                header: ['3h', '6h', '12h', '1d', '7d'],
                f: (coin, avg) => {
                    return  this.money(this.fees(this.data[coin].averages[avg]))
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
    methods: {
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                backgroundColor: this.btc_cbg,
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => this.fees(obj)),
                fill: true,
            })

            this.add_dataset({
                label: 'BCH',
                backgroundColor: this.bch_cbg,
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.fees(obj)),
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
                                labelString: 'Fees (USD)'
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

            this.draw('canvas-fees', config, labels)
        },
    }
}
</script>

<style>

</style>
