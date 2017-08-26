
<template>
    <div>
        <p id="exp">
            Block difficulty.
        </p>
        <div style="height:512px; width: 100%;">
            <canvas id="canvas-difficulty" style="height: 512px;"></canvas>
        </div>
        <br>
        <p>
            The BCH chain is currently running at
            <strong>{{ sprintf('%.2f', data.BCH.averages.last.avg_diff / data.BTC.averages.last.avg_diff * 100 ) }}%</strong> difficulty compared to the BTC chain.
        </p>
        <br>
        <hr>
        <expected></expected>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'
import expected from './difficulty-expected.vue'

export default {
    name: 'pow-difficulty',
    components: {
        expected: expected
    },
    mounted() {
        this.redraw()

        if(this.$route.hash === '#exp') {
            location.hash = ''
            this.$nextTick(() => {
                location.hash = '#exp'
            })
        }
    },
    methods: {
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            let diff = (obj) => {
                return obj.blocks > 0 ? this.sprintf('%.0f', obj.avg_diff) : null
            }

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => diff(obj)),
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => diff(obj)),
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
                                labelString: 'Difficulty'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return this.numeral('0,0', value)
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-difficulty', config, labels)
        }
    }
}
</script>

<style>

</style>
