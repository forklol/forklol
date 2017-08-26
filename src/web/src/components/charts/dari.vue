<template>
    <div>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-dari" style="height: 512px;"></canvas>
        </div>
        <br>
        <p>
            The DARI (Difficulty Adjusted Reward Index) is a way to compare the rewards of two chains that share the same reward system. The following calculation is used to determine the DARI:
            <br>
        </p>
        <div class="text-center">
            <code>
                (block coinbase+fees in satoshis) / (block difficulty) * (exchange rate<sup>[1]</sup> in USD)
            </code>
        </div>
        <br>
        <p>
            With this result we can compare the two chains. For example, using the DARI of the last block we can estimate that a {{ leader }} miner could potentially collect ({{ dari(leader) }} / {{ dari(loser)}} =)
            <strong>{{ leader_ratio }}</strong> times the reward (in USD) that a {{ loser }} miner could collect.
        </p>
        <p>
            <small>[1] The exchange rate at time the block was found.</small>
        </p>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import { sprintf } from 'sprintf-js'
import numeral from 'numeral'

export default {
    name: 'dari',
    mounted() {
        this.redraw()
    },
    methods: {
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => this.obj_dari(obj)),
                yAxisID: 'dari'
            })

            this.add_dataset({
                label: 'Difficulty',
                borderColor: this.btc_c,
                borderDash: [3, 3],
                data: this.data.BTC.history['all'].map((obj) => this.obj_diff(obj)),
                yAxisID: 'diff'
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.obj_dari(obj)),
                yAxisID: 'dari'
            })

            this.add_dataset({
                label: 'Difficulty',
                borderColor: this.bch_c,
                borderDash: [3, 3],
                data: this.data.BCH.history['all'].map((obj) => this.obj_diff(obj)),
                yAxisID: 'diff'
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
                            id: 'dari',
                            position: 'right',
                            scaleLabel: {
                                display: true,
                                labelString: 'DARI'
                            }
                        }, {
                            display: false,
                            id: 'diff',
                            scaleLabel: {
                                display: true,
                                labelString: 'Relative difficulty'
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-dari', config, labels)
        },
        obj_dari(obj) {
            return obj.blocks > 0 && obj.price > 1 ? this.sprintf('%.2f', obj.dari) : null
        },
        obj_diff(obj) {
            return obj.blocks > 0 && obj.price > 1 ? this.sprintf('%.2f', obj.avg_diff / 10000000000) : null
        },
    }
}
</script>

<style>

</style>
