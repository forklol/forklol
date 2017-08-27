
<template>
    <div>
        <p>
            Average blocksizes.
        </p>
        <div style="height:512px; width: 100%;">
            <canvas id="canvas-blocksize" style="height: 512px;"></canvas>
        </div>
        <br>
        <stat-table :stats="stats"></stat-table>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'
import { sprintf } from 'sprintf-js'

export default {
    name: 'blocks-size',
    data() {
        return {

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
                    return this.sprintf('%.2f', this.data[coin].averages[avg].size / 1000)
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
                data: this.data.BTC.history['all'].map((obj) => this.sprintf('%.2f', obj.size / 1000)),
                fill: true
            })

            this.add_dataset({
                label: 'BCH',
                backgroundColor: this.bch_cbg,
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => obj.blocks > 0 ? this.sprintf('%.2f', obj.size / 1000) : null),
                fill: true
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
                                labelString: 'Avg. blocksize'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return value + ' KiB'
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-blocksize', config, labels)
        },
    },
    destroy() {
        this.chart.destroy()
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
