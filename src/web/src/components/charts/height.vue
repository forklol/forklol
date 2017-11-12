
<template>
    <div>
        <p>
            Block height. {{ height_diff_txt }}
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-height" style="height: 512px;"></canvas>
        </div>
        <br>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'
import { sprintf } from 'sprintf-js'

export default {
    name: 'blocks-height',
    data() {
        return {
            prevh: {
                BTC: null,
                BCH: null
            }
        }
    },
    mounted() {
        this.redraw()
    },
    computed: {
        height_diff_txt() {
            let bh = this.data.BTC.averages.last.height
            let ch = this.data.BCH.averages.last.height

            if(bh > ch) {
                return 'BTC is currently ahead by ' + (bh-ch) + ' blocks.'
            } else {
                return 'BCH is currently ahead by ' + (ch-bh) + ' blocks.'
            }
        },
    },
    methods: {
        height(obj, coin, other = null) {
            let h = obj.height > 0 ? obj.height : null

            if(h) {
                this.prevh[coin] = h
                return h
            }

            if(this.prevh[coin] > 481800) {
                return this.prevh[coin]
            }

            return null
        },
        redraw() {

            this.prevh.BTC = null
            this.prevh.BCH = null

            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                backgroundColor: this.btc_cbg,
                borderColor: this.btc_c,
                spanGaps: true,
                data: this.data.BTC.history['all'].map((obj) => this.height(obj, 'BTC')),
            })

            this.add_dataset({
                label: 'BCH',
                backgroundColor: this.bch_cbg,
                borderColor: this.bch_c,
                spanGaps: true,
                data: this.data.BCH.history['all'].map((obj) => this.height(obj, 'BCH', 'BTC')),
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
                                labelString: 'Block chain height'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return value
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-height', config, labels)
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
