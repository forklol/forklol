
<template>
    <div>
        <p>
            Sum of coindays destroyed in period (log scale).
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-cdd" style="height: 512px;"></canvas>
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
    name: 'block-cdd',
    mounted() {
        this.redraw()
    },
    computed: {
        stats() {
            return {
                header: ['3h', '1d', '7d'],
                f: (coin, avg) => {
                    return this.numeral('0,0', this.cdd(this.data[coin].averages[avg]))
                },
                coins: [{
                    coin: 'BTC',
                    values: ['3h', '1d', '7d']
                }, {
                    coin: 'BCH',
                    values: ['3h', '1d', '7d']
                }]
            }

        }
    },
    methods: {
        cdd(obj) {
            return obj.blocks > 0 ? this.sprintf('%.2f', obj.cdd * obj.blocks) : null
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => this.cdd(obj)),
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.cdd(obj)),
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
                            type: 'logarithmic',
                            scaleLabel: {
                                display: true,
                                labelString: 'Coindays destroyed'
                            },
                            ticks: {
                                callback: (value, index, ticks) => {
                                    if (!(index % parseInt(ticks.length / 10))) {
                                        return this.numeral('0,0', value)
                                    }
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-cdd', config, labels)
        }
    }
}
</script>

<style>

</style>
