
<template>
    <div>
        <p>
            Hashrate divided by difficulty. A ratio of > 1.0 means (on average) faster blocks, < 1.0 slower. (log scale, 3h averages)</p>
            <div style="height:100%; width: 100%;">
                <canvas id="canvas-speed" style="height: 512px;"></canvas>
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
    name: 'pow-speed',
    mounted() {
        this.redraw()
    },
    computed: {
        stats() {
            return {
                header: ['3h', '12h', '1d', '3d', '7d'],
                f: (coin, avg) => {
                    return this.speed(this.data[coin].averages.last, 'rate'+avg)
                },
                coins: [{
                    coin: 'BTC',
                    values: ['3h', '12h', '1d', '3d', '7d']
                }, {
                    coin: 'BCH',
                    values: ['3h', '12h', '1d', '3d', '7d']
                }]
            }

        }
    },
    methods: {
        speed(obj, rate = null) {
            if(obj.blocks < 1) return '0.00'

            if(rate !== null) {
                return this.sprintf('%.2f', obj[rate] / obj.avg_diff)
            }

            return this.sprintf('%.2f', obj.rate / obj.avg_diff)
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => this.speed(obj, 'rate3h'))
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.speed(obj, 'rate3h'))
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
                                labelString: 'Hashrate / difficulty'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    if(value < 1 && index % 2) return null
                                    return this.sprintf('%.2f', value)
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-speed', config, labels)
        }
    }
}
</script>

<style>

</style>
