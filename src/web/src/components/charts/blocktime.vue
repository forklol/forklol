
<template>
    <div>
        <p>
            Average number of blocks found every hour.
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-blocktime" style="height: 644px;"></canvas>
        </div>
        <br>
        <small>
            <table class="table table-sm table-bordered">
                <thead class="thead-dark">
                    <tr>
                        <th>Coin</th>
                        <th>3h</th>
                        <th>6h</th>
                        <th>12h</th>
                        <th>1d</th>
                        <th>7d</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>BTC</td>
                        <td>{{ blocks_hr('BTC', '3h', 3 * 60) }}</td>
                        <td>{{ blocks_hr('BTC', '6h', 6 * 60) }}</td>
                        <td>{{ blocks_hr('BTC', '12h', 12 * 60) }}</td>
                        <td>{{ blocks_hr('BTC', '1d', 24*60) }}</td>
                        <td>{{ blocks_hr('BTC', '7d', 7 * 24 * 60) }}</td>
                    </tr>
                    <tr>
                        <td>BCH</td>
                        <td>{{ blocks_hr('BCH', '3h', 3 * 60) }}</td>
                        <td>{{ blocks_hr('BCH', '6h', 6 * 60) }}</td>
                        <td>{{ blocks_hr('BCH', '12h', 12 * 60) }}</td>
                        <td>{{ blocks_hr('BCH', '1d', 24*60) }}</td>
                        <td>{{ blocks_hr('BCH', '7d', 7 * 24 * 60) }}</td>
                    </tr>
                </tbody>
            </table>
        </small>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'
import { sprintf } from 'sprintf-js'

export default {
    name: 'blocks-time',
    data() {
        return {
        }
    },
    mounted() {
        this.redraw()
    },
    methods: {
        blocks_hr(coin, period, min) {
            let obj = this.data[coin].averages[period]
            if(obj.blocks < 1) return '0.00'
            return obj.blocks > 0 ? (sprintf('%.2f', (3600 / (min * 60)) * obj.blocks)) : null
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            let t = this.period_time()
            let blocktime = (obj, coin) => {
                let b = 3600 / t * obj.blocks
                return obj.blocks > 0 ? (sprintf('%.2f', b)) : 0
            }

            let btc_data = this.data.BTC.history['all'].map((obj) => blocktime(obj, 'BTC'))
            let bch_data = this.data.BCH.history['all'].map((obj) => blocktime(obj, 'BCH'))

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                data: btc_data,
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                data: bch_data,
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
                                labelString: 'Blocks per hour'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    if(value != 5 && value != 3 && value < 10 && index % 2) return null
                                    return this.sprintf('%.2f', value)
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-blocktime', config, labels)
        }
    }
}
</script>

<style>

</style>
